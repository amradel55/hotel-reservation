package api

import (
	"fmt"

	"github.com/amradel55/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// to-do need to think in a better way to remove roomStore dependance from hotelHandler since the hotel is higher than room.

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hotelStore db.HotelStore, roomStore db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating float64
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qParams HotelQueryParams
	if err := c.QueryParser(&qParams); err != nil {
		return err
	}
	fmt.Println(qParams)

	hotels, err := h.hotelStore.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

// func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
// 	var params types.Cre
// }
