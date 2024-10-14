package db

import "context"

const (
	DBNAME     = "hotel_reservation"
	TestDBNAME = "hotel_reservation_test"
	DBURI      = "mongodb://localhost:27017"
)

type Dropper interface {
	Drop(context.Context) error
}
