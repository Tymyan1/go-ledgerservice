package main

import (
	"ledgerservice/api"
	"ledgerservice/db"
	"log"
	"net/http"
)

func main() {

	var dbCon, err = db.CreateConnection()
	if err != nil {
		log.Print("Connecting to db failed")
		panic(err)
	}

	var transactionDb = &db.TransactionDb{DB: dbCon}

	http.HandleFunc("POST /transactions", injectDbConn(transactionDb, api.PostTransaction))

	log.Print("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func injectDbConn(db *db.TransactionDb, handler func(http.ResponseWriter, *http.Request, *db.TransactionDb)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, db)
	}
}
