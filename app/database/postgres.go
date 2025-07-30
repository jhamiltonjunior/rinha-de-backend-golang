package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var PostgresClient *sql.DB

func InitializePostgresDB() *sql.DB {
	db, err := sql.Open("postgres", "user=admin password=admin123 dbname=rinha_de_backend sslmode=disable host=postgres port=5432")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	PostgresClient = db
	return db
}

func CreatePaymentHistoryPostgres(db *sql.DB, paymentData map[string]interface{}, typeService string) {
	correlationId := paymentData["correlationId"]
	amount := paymentData["amount"]
	requestedAt := paymentData["requestedAt"]

	stmt, err := db.Prepare("INSERT INTO payment_history (correlation_id, amount, requested_at, type) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(correlationId, amount, requestedAt, typeService)

	if err != nil {
		log.Printf("Error inserting payment history: %v\n", err)
	}
	
}

func UpdatePaymentHistoryPostgres(db *sql.DB, correlationId string, typeService string) {
	stmt, err := db.Prepare("UPDATE payment_history SET type = $1 WHERE correlation_id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(typeService, correlationId)
	if err != nil {
		log.Printf("Error updating payment history: %v\n", err)
	}
}

func GetPaymentsHistoryPostgres(db *sql.DB, from, to string) ([]map[string]interface{}, error) {
	query := `
		SELECT
			type, COUNT(*) AS "totalRequests", SUM(amount) AS "totalAmount"
		FROM
			payment_history
		WHERE requested_at BETWEEN $1 AND $2
		GROUP BY type
	`

	rows, err := db.Query(query, from, to)
	if err != nil {
		log.Printf("Error retrieving payment history: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var paymentData []map[string]interface{}
	for rows.Next() {
		var typeServicee string
		var totalRequests int
		var totalAmount float64
		if err := rows.Scan(&typeServicee, &totalRequests, &totalAmount); err != nil {
			log.Printf("Error scanning payment history row: %v\n", err)
			continue
		}
		paymentData = append(paymentData, map[string]interface{}{
			"type":          typeServicee,
			"totalRequests": totalRequests,
			"totalAmount":   totalAmount,
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating payment history rows: %v\n", err)
		return nil, err
	}

	return paymentData, nil
}

func PurgePaymentsHistoryPostgres(db *sql.DB) {
	query := "DELETE FROM payment_history"
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error purging payment history: %v\n", err)
	}
}
