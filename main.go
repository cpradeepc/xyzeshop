package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"xyzeshop/constval"
	"xyzeshop/dbs"
	"xyzeshop/helper"
	"xyzeshop/payloads"
	"xyzeshop/router"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

//var srv *http.Server

var (
	client  *mongo.Client
	contx   context.Context
	contxFn context.CancelFunc
)

func init() {
	fs, err := os.Stat(".env")
	log.Println("env file :", fs)
	if err == nil {
		log.Println("loading... config file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("error  was invoked during loading config file :", err)
		}

		log.Println("succeed load config...")
	}
	client, contx, contxFn = dbs.ConDb()

	//creating system admin
	hashPassword := helper.GenPassHash("1234")
	log.Println("hassPsw :", hashPassword)
	user := payloads.RegdUser{
		Name:     "Admin",
		Email:    "xxxx@xxxx.in",
		Password: hashPassword,
		UserType: constval.AdminUser,
	}

	// log.Println("user data:", user)
	userVarify, err := dbs.Mgr.GetSingleRecordByEmail(user.Email, constval.UserCollection)
	if err != nil {
		_, err := dbs.Mgr.InsertOne(user, constval.UserCollection)
		if err != nil {
			log.Fatal("error was invoked during insert record :", err)
		}
	}
	log.Println("user verify :", userVarify)

	log.Println("init fn done")
}
func main() {
	handler := router.GuestRoutes()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: handler,
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shutdown := helper.GracefulShutdown(srv)

	select {
	case err := <-srvErr:
		shutdown(err)
	case sig := <-quit:
		//dbs.CloseCon(client, contx, contxFn)
		shutdown(sig)
	}
	log.Println("Server exiting")

}
