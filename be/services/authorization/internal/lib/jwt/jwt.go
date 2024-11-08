package jwt

import (
	"fmt"
	"lib/generator"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=30"`
	LastName    string `json:"last_name" validate:"required,min=2,max=30"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	IP          string `json:"ip"`
	Jti         int    `json:"jti"`
	jwt.StandardClaims
}

// TokenGenerator генерирует JWT с указанием параметров id и ip адреса пользователя
func GenerateToken(
	firstName string,
	lastName string,
	email string,
	phoneNumber string,
) (accessToken string, refreshToken string, err error) {

	claims := &SignedDetails{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Jti:         generator.GenerateRandomInt(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	err = WriteInMap(accessToken, ts{
		accessToken,
		refreshToken,
		true,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// ValidateToken валидирует и декодирует информацию из токена
func ValidateToken(accessToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("the token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token is already expired")
	}

	return claims, nil
}

// RequireTokens проверяет валидность токена
func RequireTokens(accessToken, refreshToken string) error {
	ts, err := ReadFromMap(accessToken)
	if err != nil {
		return err
	}
	if !ts.Valid {
		return fmt.Errorf("token already used")
	}
	if refreshToken != ts.refreshToken {
		return fmt.Errorf("token connectivity is broken")
	}
	if accessToken != ts.accessToken {
		return fmt.Errorf("token connectivity is broken")
	}

	ts.Valid = false
	return nil
}
