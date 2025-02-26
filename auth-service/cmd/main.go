package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type VKResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/auth/vk", func(c *gin.Context) {
		clientID := os.Getenv("VK_CLIENT_ID")
		redirectURI := "http://localhost:8080/auth/vk/callback"
		authURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&display=page&redirect_uri=%s&scope=email&response_type=code&v=5.131", clientID, redirectURI)
		c.Redirect(http.StatusFound, authURL)
	})

	r.GET("/auth/vk/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		clientID := os.Getenv("VK_CLIENT_ID")
		clientSecret := os.Getenv("VK_CLIENT_SECRET")
		redirectURI := "http://localhost:8080/auth/vk/callback"

		accessTokenURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s", clientID, clientSecret, redirectURI, code)
		resp, err := http.Get(accessTokenURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var vkResp VKResponse
		if err := json.Unmarshal(body, &vkResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": vkResp.AccessToken,
			"user_id":      vkResp.UserID,
		})
	})

	r.Run(":8080")
}
