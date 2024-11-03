package db

import "context"

const (
	DBNAME     = "hotel_reservation"
	TestDBNAME = "hotel_reservation_test"
	DBURI      = "mongodb://localhost:27017"
)

// does that give the access to hole stores even the domain doesn't need it?
type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Dropper interface {
	Drop(context.Context) error
}
