package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "bar"
	bytedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	stringPassword := hex.EncodeToString(bytedPassword)
	fmt.Println(stringPassword)
	comparePassword, err := hex.DecodeString(stringPassword)
	err = bcrypt.CompareHashAndPassword(comparePassword, []byte(password))
	if err != nil {
		panic(err)
	}
}
