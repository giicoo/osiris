package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/giicoo/osiris/notification-service/internal/entity"
	"github.com/gin-gonic/gin"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("get header: %w", err).Error(),
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, " ")
		if len(idTokenHeader) < 2 {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Must provide Authorization header with format `Bearer {token}`",
			})
			c.Abort()
			return
		}
		url := fmt.Sprintf("http://auth-service:8080/auth/%s", idTokenHeader[1])

		r, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("get user info: %w", err).Error(),
			})
			c.Abort()
			return
		}
		defer r.Body.Close()

		user := &entity.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("decode user info: %w", err).Error(),
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
