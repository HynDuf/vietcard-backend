package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestHandler interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	RefreshToken(c *gin.Context)
	CreateCard(c *gin.Context)
	CreateDeck(c *gin.Context)
	GetDeckWithReviewCards(c *gin.Context)
}

func GetLoggedInUserID(c *gin.Context) (string, error) {
	uID, isExisted := c.Get("x-user-id")
	if !isExisted {
		return "", errors.New("Missing x-user-id (set at middleware) in Gin Context")
	}

	return uID.(string), nil
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupRequest struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateCardRequest struct {
	UserID       primitive.ObjectID `json:"user_id" swaggerignore:"true"`
	DeckID       primitive.ObjectID `json:"deck_id" binding:"required"`
	Question     string             `json:"question" binding:"required"`
	Answer       string             `json:"answer" binding:"required"`
	WrongAnswers []string           `json:"wrong_answers" binding:"required"`
}

type CreateCardResponse struct {
	Success bool `json:"success"`
}

type CreateDeckRequest struct {
	UserID primitive.ObjectID `json:"user_id" swaggerignore:"true"`
	Name   string             `json:"name" binding:"required"`
}

type CreateDeckResponse struct {
	Success bool `json:"success"`
}
