package utils

import "crypto/rand"

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

const str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) *string {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil
	}

	strLength := len(str)
	for i := 0; i < length; i++ {
		buffer[i] = str[int(buffer[i])%strLength]
	}
	ranstr := string(buffer)

	return &ranstr
}

const CapitalNumber = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomCapitalNumber(length int) *string {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil
	}

	strLength := len(str)
	for i := 0; i < length; i++ {
		buffer[i] = str[int(buffer[i])%strLength]
	}
	ranstr := string(buffer)

	return &ranstr
}
