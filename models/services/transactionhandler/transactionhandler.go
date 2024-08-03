package transactionhandler

import "database/sql"

type TransactionHandler struct {
    DB *sql.DB
}

type ITransactionHandler interface {
    Commit() error
    Rollback() error
}
