package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
)

type PaymentWorker struct {
	Body              []byte
	VouTeDarOContexto context.Context
	RetryCount        int
}

var (
	SegureOChann  = make(chan PaymentWorker, 3000)
	SegureOChann2 = make(chan PaymentWorker, 3000)
)

func InitializeWorker(client *sql.DB) {
	defaultURL := os.Getenv("PAYMENT_PROCESSOR_URL_DEFAULT")
	fallbackURL := os.Getenv("PAYMENT_PROCESSOR_URL_FALLBACK")
	const numWorkers = 20

	for i := 1; i <= numWorkers; i++ {
		go workerLoop(client, defaultURL, fallbackURL)
	}

	for i := 1; i <= numWorkers; i++ {
		go retryworkLoop(defaultURL, fallbackURL, client)
	}
}

func workerFunc(client *sql.DB, defaultURL, fallbackURL string, payment PaymentWorker) bool {
	body, ok := ProcessPayment(payment.Body, payment.VouTeDarOContexto, defaultURL)

	if ok {
		database.CreatePaymentHistoryPostgres(client, body, "default")
		return true
	}

	if payment.RetryCount <= 15 {
		return false
	}

	body, ok = ProcessPayment(payment.Body, payment.VouTeDarOContexto, fallbackURL)
	if ok {
		database.CreatePaymentHistoryPostgres(client, body, "fallback")
		return true
	}

	return false
}

func workerLoop(client *sql.DB, defaultURL, fallbackURL string) {
	for payment := range SegureOChann {
		if !workerFunc(client, defaultURL, fallbackURL, payment) {
			SegureOChann2 <- PaymentWorker{
				Body:              payment.Body,
				VouTeDarOContexto: context.TODO(),
			}
		}
	}
}

func retryworkLoop(defaultURL, fallbackURL string, client *sql.DB) {
	for payment := range SegureOChann2 {
		if payment.RetryCount >= 20 {
			continue
		}

		func() {
			cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			payment.VouTeDarOContexto = cxt

			if !workerFunc(client, defaultURL, fallbackURL, payment) {
				payment.RetryCount++
				SegureOChann2 <- payment
			}
		}()
	}
}

func ProcessPayment(paymentBytes []byte, ctx context.Context, theBestURLEver string) (map[string]any, bool) {
	var payment map[string]any
	if err := json.Unmarshal(paymentBytes, &payment); err != nil {
		return nil, false
	}

	// payment["correlationId"], _ = newUUID()

	now := time.Now().UTC()
	isoString := "2006-01-02T15:04:05.000Z"
	payment["requestedAt"] = now.Format(isoString)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Erro ao serializar o pagamento:", err)
		return nil, false
	}

	return payment, sendToPaymentService(paymentBytes, theBestURLEver, ctx)
}

func sendToPaymentService(paymentBytes []byte, reqURL string, ctx context.Context) bool {
	_, status := utils.Request("POST", paymentBytes, reqURL+"/payments", ctx)
	// fmt.Printf("status: %d\n", status)
	return status == 200 || status == 201
}
