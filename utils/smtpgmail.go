package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func InitEmail(email string) (string, error) {
	senderEmail := os.Getenv("CONFIG_AUTH_EMAIL")
	to := []string{email}
	cc := []string{senderEmail}
	Generateotp, _ := GenerateOTP(6)
	subject := "Test mail"
	hashCode := EncodeBase64([]byte(Generateotp))
	fmt.Println(hashCode)
	fmt.Println(Generateotp)
	message := fmt.Sprintf("Kode Verifikasi : %s \n Link Verifikasi : https://api-desabangkit.herokuapp.com/users/verification-account?code=%s", Generateotp, hashCode)

	err := sendMail(to, cc, subject, message)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
	return Generateotp, nil
}

func sendMail(to []string, cc []string, subject, message string) error {
	CONFIG_SMTP_HOST := os.Getenv("CONFIG_SMTP_HOST")
	CONFIG_SMTP_PORT := 587
	CONFIG_AUTH_EMAIL := os.Getenv("CONFIG_AUTH_EMAIL")
	CONFIG_SENDER_NAME := os.Getenv("CONFIG_SENDER_NAME")
	CONFIG_AUTH_PASSWORD := os.Getenv("CONFIG_AUTH_PASSWORD")

	body := "From: " + CONFIG_SENDER_NAME + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_AUTH_EMAIL, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
