package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type message struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

type response struct {
	Status       string
	Message      string
	EmailMessage string
	EmailSender  string
	EmailAddress string
	TimeTaken    time.Duration
}

func sendMail(msg message) (string, error) {

	if msg.Message == "" {
		return "Error whilst sending email, message was not provided", errors.New("External Client Error")
	} else if msg.Sender == "" {
		return "Error whilst sending email, sender was not provided", errors.New("External Client Error")
	} else if msg.Email == "" {
		return "Error whilst sending email, email was not provided", errors.New("External Client Error")
	}

	PORT, err := strconv.Atoi(os.Getenv("GPORT"))

	HOST := os.Getenv("HOST")
	EMAIL := os.Getenv("EMAIL")
	PWD := os.Getenv("PWD")
	RECIPIENT := os.Getenv("RECIPIENT")

	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", EMAIL)
	m.SetHeader("To", RECIPIENT)
	m.SetHeader("Subject", ("Message from " + msg.Sender + "."))
	m.SetBody("text/html", "<p>Message: "+msg.Message+"</p>"+"<h4>From the address: "+msg.Email+"</h4>")

	d := gomail.NewDialer(HOST, PORT, EMAIL, PWD)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return "Successfully sent email", nil
}
func mailHandle(c *gin.Context) {
	start := time.Now()

	var msg message

	if err := c.BindJSON(&msg); err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to receive request parameters. Please try again.",
		})
		c.AbortWithError(400, err)
	}

	sendMessage, err := sendMail(msg)

	if err != nil {

		res := response{
			Status:       "Failure",
			Message:      sendMessage,
			EmailMessage: msg.Message,
			EmailSender:  msg.Sender,
			EmailAddress: msg.Email,
		}

		jData, _ := json.Marshal(res)

		c.JSON(400, jData)

		c.AbortWithError(400, err)
	}

	elapsed := time.Since(start)

	res := response{
		Status:       "Successful",
		Message:      sendMessage,
		EmailMessage: msg.Message,
		EmailSender:  msg.Sender,
		EmailAddress: msg.Email,
		TimeTaken:    time.Duration(elapsed.Seconds()),
	}

	jRes, err := json.Marshal(res)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	c.JSON(200, jRes)

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-sfs-embed")
		c.Header("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/mail", mailHandle)
	r.Run()
}
