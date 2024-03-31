package helper

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"
	"xyzeshop/constval"
	"xyzeshop/payloads"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMailSendGrid(req payloads.Verification) (payloads.Verification, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return req, errors.New("SENDGRID_API_KEY environment")
	}
	//create a sendGrid client
	client := sendgrid.NewSendClient(apiKey)

	//setup the email message
	from := mail.NewEmail("Sender Name", constval.Sender)
	to := mail.NewEmail("Recipient Name", req.Email)
	subject := "OTP verification mail"
	otp := RandomNum()
	req.OTP = int64(otp)
	htmlContent := "<p>This is a test otp for verification <strong>" + strconv.Itoa(otp) + "</strong></p>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	//send the email message
	_, err := client.Send(message)
	if err != nil {
		return req, err
	}
	return req, err
}

// Create a new random number generator with a custom seed (e.g., current time)
// source := rand.NewSource(time.Now().UnixNano())
// rng := rand.New(source)
// Generate a random number of minutes between 1 and 15
// randomMinutes := rng.Intn(15) + 1
func RandomNum() int {
	randSource := rand.NewSource(time.Now().UnixNano())
	rndRang := rand.New(randSource)
	randNum := rndRang.Intn(1000) + 1000
	return randNum
}
