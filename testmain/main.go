package main

import (
	"fmt"
	"xyzeshop/helper"
	"xyzeshop/payloads"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("err env: ", err)
	}

	v := payloads.Verification{}
	v.Email = "pkumdeep@gmail.com"

	data, err := helper.SendMailOtp(v)
	fmt.Println("data: ", data, " err: ", err)

}
