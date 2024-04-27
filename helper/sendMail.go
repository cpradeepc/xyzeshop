package helper

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"math/big"
	_ "math/rand"

	"os"
	"strconv"
	"xyzeshop/payloads"

	gomail "gopkg.in/mail.v2"
)

func SendMailOtp(req payloads.Verification) (payloads.Verification, error) {
	mailPassKey := os.Getenv("MAIL_PASSKEY")
	if mailPassKey == "" {
		return req, errors.New("error in mail passkey")
	}

	msg := gomail.NewMessage()

	//for sender
	msg.SetHeader("From", os.Getenv("MAIL_USER"))

	//for receiver
	msg.SetHeader("To", req.Email)

	//for subject
	msg.SetHeader("Subject", "Verify OTP")
	otp, err := GenerateOTPCode(4)
	if err != nil {
		return req, errors.New("error in generate otp fn")
	}
	req.OTP, err = strconv.ParseInt(otp, 10, 64)
	if err != nil {
		return req, errors.New("error in parse otp fn")
	}

	//set mail body in html
	//htmlContent := "<p>This is OTP for verification <strong>" + otp + "</strong></p>"
	msg.SetBody("text/plain", "This is OTP for verification : "+otp)
	mailHost := os.Getenv("SMTP_HOST")
	mailPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	mailUser := os.Getenv("MAIL_USER")

	dialer := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassKey)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = dialer.DialAndSend(msg)
	if err != nil {
		return req, err
	}
	return req, nil
}

func GenerateOTPCode(length int) (string, error) {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice), nil
}
