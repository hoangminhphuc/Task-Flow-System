package jwt

import (
	"first-proj/common"
	"first-proj/plugin/tokenprovider"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtProvider struct {
	prefix string
	secret string
}

func NewTokenJWTProvider(prefix string, secret string) *jwtProvider {
	return &jwtProvider{prefix: prefix, secret: secret}
}

type myClaims struct {
	Payload common.TokenPayLoad `json:"payload"`
	jwt.StandardClaims
}


//Created and Expiry can be extracted from JWT, not neccesary here
type token struct {
	Token string `json:"token"`
	Created time.Time `json:"created"`
	Expiry int `json:"expiry"`
}

func (t* token) GetToken() string {
	return t.Token
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayLoad, expiry int) (tokenprovider.Token, error) {
	
	now := time.Now()
	
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		common.TokenPayLoad{
			UId: data.UserId(),
			URole: data.Role(),
		},
		jwt.StandardClaims{
			ExpiresAt: now.Local().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt: now.Local().Unix(),
			Id: fmt.Sprintf("%d", now.Local().UnixNano()),
		},
	})


	//return final token with 3 parts
	tokenString, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return nil, err
	}

	return &token{
		Token: tokenString,
		Expiry: expiry,
		Created: now,
	}, nil
}


/* 
* Validate a JWT, ensure its signature is correct,
* check its validity (including expiration), 
* and return the embedded custom payload if all checks pass.
*/
func (j *jwtProvider) Validate(myToken string) (tokenprovider.TokenPayLoad, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Return a secret key (convert to byte slice)
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	//Type assertion
	claims, ok := res.Claims.(*myClaims)

	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	return claims.Payload, nil
}