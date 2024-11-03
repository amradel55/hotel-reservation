package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/amradel55/hotel-reservation/db"
	"github.com/amradel55/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

// to-do apply different type of case like auth failure with wrong password or wrong email.

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "amr.adel@hotel.com",
		FirstName: "Amr",
		LastName:  "Adel",
		Password:  "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}
func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "amr.adel@hotel.com",
		Password: "test",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("content-type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d ", res.StatusCode)
	}
	var authRes AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		t.Error(err)
	}
	if authRes.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}
	// set the encrypted Password to empty string, because we will not return this in any json response.
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authRes.User) {
		t.Fatalf("expected the user to be the inserted user")
	}

}
