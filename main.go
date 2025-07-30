package main

import (
	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
	"github.com/jhamiltonjunior/rinha-de-backend/app/server"
	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	_ "github.com/lib/pq"
)

func main() {
	db := database.InitializePostgresDB()
	defer db.Close()

	worker.InitializeWorker(db)

	// go pingQuantityOfSegureOChann()

	server.ListenAndServe("3000")
}

// func pingQuantityOfSegureOChann() {
// 	for {
// 		fmt.Printf("Quantidade de pagamentos pendentes: %d\n", len(worker.SegureOChann))
// 		fmt.Printf("Quantidade de pagamentos em retry: %d\n", len(worker.SegureOChann2))
// 		fmt.Println("========================================")
// 		time.Sleep(5 * time.Second)
// 	}
// }
