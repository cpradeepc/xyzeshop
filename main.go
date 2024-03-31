package main

import (
	"log"
	"os"
	"xyzeshop/constval"
	"xyzeshop/dbs"
	"xyzeshop/helper"
	"xyzeshop/payloads"
	"xyzeshop/router"

	"github.com/joho/godotenv"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("loading... config file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("error  was invoked during loading config file :", err)
		}
		log.Println("succeed load config...")
	}

	//creating system admin
	hashPassword := helper.GenPassHash("1234")
	user := payloads.RegdUser{
		Name:     "Admin",
		Email:    "xxxx@xxxx.com",
		Password: hashPassword,
		UserType: constval.AdminUser,
	}

	userVarify := dbs.Mgr.GetSingleRecordByEmail(user.Email, constval.UserCollection)
	if userVarify.Email == "" {
		_, err := dbs.Mgr.InsertOne(user, constval.UserCollection)
		if err != nil {
			log.Fatal("error was invoked during insert record :", err)
		}
	}
}
func main() {
	//calling direct routes by function
	router.GuestRoutes()
}
