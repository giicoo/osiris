package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string `json:"photo_400_orig"`
	City      City   `json:"city"`
}
type City struct {
	Title string `json:"title"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{},
		Endpoint:     vk.Endpoint,
	}
	r.GET("/", func(c *gin.Context) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		// получаем URL для редиректа на OAuth API VK и передаем его в темплейт
		c.HTML(http.StatusOK, "index.html", gin.H{
			"authUrl": url,
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		ctx := context.Background()
		// получаем код от API VK из квери стринга
		authCode := c.Request.URL.Query()["code"]
		// меняем код на access токен
		tok, err := conf.Exchange(ctx, authCode[0])
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, tok)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
