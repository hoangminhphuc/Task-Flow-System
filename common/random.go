package common

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSequence(n int) string {
    b := make([]rune, n)

    //Generate new seed for random function
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    for i := range b {
        b[i] = letters[r1.Intn(999999) % len(letters)]
    }

    return string(b)
}

func GenSalt(length int) string {
    if length < 0 {
        length = 50
    }
    return randSequence(length)
}



type BcryptHash struct{}

func NewBcryptHash() *BcryptHash {
    return &BcryptHash{}
}

func (*BcryptHash) Hash(plainPassword string) (string) {
	// Set cost = 12. Higher cost = more secure, but slower.

    fmt.Println("Plain password: ", plainPassword)
    
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), BCRYPT_COST)
	if err != nil {
		return ""
	}

    fmt.Println("Hashed bytes: ", hashedBytes)
	return string(hashedBytes)
}

func (*BcryptHash) Compare(hashedValue, plainText string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(plainText))
    return err == nil
}