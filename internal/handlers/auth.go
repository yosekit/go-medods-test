package handlers

import (
	"encoding/base64"
	"net/http"

	repos "goauth/internal/repos/postgres"
	"goauth/internal/services"

	"github.com/gin-gonic/gin"
)

// POST /auth/tokens
// Query Parameters:
// - user_guid (required): string
var GenerateTokens gin.HandlerFunc = func(ctx *gin.Context) {
	userID := ctx.Query("user_guid")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}
	clientIP := ctx.ClientIP()

	accessToken, err := services.GenerateAccessToken(userID, clientIP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken := services.GenerateRefreshToken()

	err = repos.SaveRefreshToken(userID, clientIP, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
		return
	}

	encodedRefreshToken := encodeBase64(refreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"refreshToken": encodedRefreshToken,
	})
}


// POST /auth/refresh
// Body Raw (JSON) Parameters:
// - accessToken (required): string
// - refreshToken (requiered): base64 string
var RefreshTokens gin.HandlerFunc = func(ctx *gin.Context) {
	var request struct {
		AccessToken string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tokenClaims, err := services.ParseAccessToken(request.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID := tokenClaims["sub"].(string)
	originalIP := tokenClaims["ip"].(string)

	decodedRefreshToken, err := decodeBase64(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token format"})
		return
	}

	err = repos.ValidateRefreshToken(userID, decodedRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	currentIP := ctx.ClientIP()
	if originalIP != currentIP {
		// TODO отправить email-уведомелние об измененном IP-адресе
		// go sendEmailWarning(userID, originalIP, currentIP)
	} 

	accessToken, err := services.GenerateAccessToken(userID, currentIP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken := services.GenerateRefreshToken()

	err = repos.UpdateRefreshToken(userID, currentIP, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update refresh token"})
		return
	}

	encodedRefreshToken := encodeBase64(refreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"refreshToken": encodedRefreshToken,
	})
}


func encodeBase64(token []byte) string {
	return base64.StdEncoding.EncodeToString(token)
}

func decodeBase64(token string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(token)
}