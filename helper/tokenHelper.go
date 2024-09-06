package helper

import (
	"booking-app/database"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Firstname string
	Lastname  string
	Email     string
	Uid       string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(firstname string, lastname string, email string, uid string) (signedToken string, signedRefreshToken string) {
	claims := &SignedDetails{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken
}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		/* msg = err.Error() */
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		/* msg = err.Error() */
		return
	}
	return claims, msg
}
func UpdateAllTokens(signedToken string, signedRefreshToken string, Uid string) {
	var updateObj primitive.D
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	upsert := true
	filter := bson.M{"Password": Uid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := database.UserCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	/* return */
}
