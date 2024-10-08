package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var usedTokens = make(map[string]bool)

func createTokens(username string) (string, string, error) {
	accessToken, err := GenerateToken(username, 15*time.Second) // Access Token 15 วินาที
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateToken(username, 30*time.Second) // Refresh Token 30 วินาที
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func validateCredentials(loginDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) bool {
	return loginDetails.Username == "user" && loginDetails.Password == "pass"
}

func loginHandler(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !validateCredentials(loginDetails) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := createTokens(loginDetails.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	usedTokens[accessToken] = true

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func refreshHandler(c *gin.Context) {
	var tokenDetails struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&tokenDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if usedTokens[tokenDetails.RefreshToken] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: Token cannot be used as refresh token"})
		return
	}

	username, err := ValidateToken(tokenDetails.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, refreshToken, err := createTokens(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	usedTokens[accessToken] = true

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
