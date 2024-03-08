package api

import (
	"errors"
	"ledgerservice/db"
	"ledgerservice/tools"
	"log"
	"net/http"
)

func PostTransaction(w http.ResponseWriter, r *http.Request, db *db.TransactionDb) {
	var txDto TransactionDto
	var err = tools.DecodeJSONBody(w, r, &txDto)

	if err != nil {
		var mr *tools.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			handleUnknownError(err, w)
		}
	}

	var errList = validate(txDto)
	if len(errList) > 0 {
		handleBodyErrors(errList, w)
		return
	}

	var tx = mapTransactionToDomain(txDto)
	err = db.Save(tx)

	if err != nil {
		handleUnknownError(err, w)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func handleUnknownError(err error, w http.ResponseWriter) {
	log.Print(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func handleBodyErrors(errs []error, w http.ResponseWriter) {
	var errMsg = ""
	for _, e := range errs {
		errMsg += e.Error() + "\n"
	}
	http.Error(w, errMsg, http.StatusBadRequest)
}
