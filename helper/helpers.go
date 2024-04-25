package helper

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"xyzeshop/payloads"

	"golang.org/x/crypto/bcrypt"
)

// convert string into int
func ConvertStrInt(str string) (int, error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		log.Println("error invoked during convert string to int")
		return 0, err
	}
	return integer, nil
}

// check user validation
func CheckUserValidation(user payloads.GuestUser) error {
	if user.Email == "" {
		return errors.New("email can not be empty")
	}
	if user.Name == "" {
		return errors.New("name can not be empty")
	}
	if user.Phone == "" {
		return errors.New("phone can not be empty")
	}
	if user.Password == "" {
		return errors.New("password can not be empty")
	}

	return nil
}

// generate password into hash
func GenPassHash(pasStr string) string {
	bts, err := bcrypt.GenerateFromPassword([]byte(pasStr), 14) //cost (4 to 31)
	if err != nil {
		return ""
	}
	return string(bts)
}

// check product validation
func CheckProductValidation(product payloads.ProductUser) error {
	if product.Description == "" {
		return errors.New("description can not be empty")
	}
	if product.Name == "" {
		return errors.New("name can not be empty")
	}
	if product.ImageUrl == "" {
		return errors.New("image can not be empty")
	}
	if product.Price == 0 {
		return errors.New("price can not be empty")
	}
	return nil
}

func GracefulShutdown(srv *http.Server) func(reason interface{}) {
	return func(reason interface{}) {
		log.Println("Server Shutdown:", reason)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error Gracefully Shutting Down API:", err)
		}
	}
}
