package server

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/jhamiltonjunior/rinha-de-backend/app/handler"
	"github.com/valyala/fasthttp"
)
func ListenAndServe(appPort string) {
	r := router.New()

	r.POST("/payments", handler.Payments)
	r.GET("/payments", handler.PaymentsGEt)
	r.GET("/payments-summary", handler.PaymentsSummary)
	r.POST("/purge-payments", handler.PaymentsPurge)

	fmt.Println("Servidor rodando em http://localhost:" + appPort)
	fasthttp.ListenAndServe(":" + appPort, r.Handler)
}