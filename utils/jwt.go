package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user *database.User) (string, error) {
    tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":  user.ID,
        "iat": time.Now().Unix(),
        "eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
    })
    return token.SignedString(privateKey)
}

func GetCurrentUser(context *gin.Context) (*database.User, error) {
    err := ValidateJWT(context)
    if err != nil {
        return &database.User{}, err
    }
    token, _ := getToken(context)
    claims, _ := token.Claims.(jwt.MapClaims)
    userId := uint(claims["id"].(float64))

    user, err := database.FindUserById(userId)
    if err != nil {
        return &database.User{}, err
    }
    return user, nil
}

func ValidateJWT(context *gin.Context) error {
    token, err := getToken(context)
    if err != nil {
        return err
    }
    _, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {
        return nil
    }
    return errors.New("invalid token provided")
}

func getToken(context *gin.Context) (*jwt.Token, error) {
    tokenString := getTokenFromRequest(context)
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return privateKey, nil
    })
    return token, err
}

func getTokenFromRequest(context *gin.Context) string {
    bearerToken := context.Request.Header.Get("Authorization")
    splitToken := strings.Split(bearerToken, " ")
    if len(splitToken) == 2 {
        return splitToken[1]
    }
    return ""
}