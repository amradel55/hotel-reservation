package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/amradel55/hotel-reservation/db"
	"github.com/amradel55/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	ToDate     time.Time `json:"toDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.ToDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	if p.ToDate.Before(p.FromDate) {
		return fmt.Errorf("cannot book a room to date before from date")
	}
	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Status: "error",
			Msg:    "Internal Server Error",
		})
	}
	if ok := h.isRoomAvailableForBooking(c, roomOID, params); !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Status: "error",
			Msg:    fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		ToDate:     params.ToDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(c *fiber.Ctx, roomID primitive.ObjectID, params BookRoomParams) bool {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"toDate": bson.M{
			"$lte": params.ToDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), where)
	return !(len(bookings) > 0 || err != nil)
}
