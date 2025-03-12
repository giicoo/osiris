package restapi

import (
	"fmt"
	"net/http"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	oauthURL    = "https://oauth.yandex.ru/authorize"
	tokenURL    = "https://oauth.yandex.ru/token"
	userInfoURL = "https://login.yandex.ru/info"
)

type Controller struct {
	cfg      *config.Config
	services *services.Services
}

func NewController(cfg *config.Config, services *services.Services) *Controller {
	return &Controller{
		cfg:      cfg,
		services: services,
	}
}

func (cont *Controller) Auth(c *gin.Context) {
	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s", oauthURL, cont.cfg.ClientID, cont.cfg.RedirectURL)
	c.Redirect(http.StatusFound, authURL)
}

func (cont *Controller) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}
	accessToken, err := cont.services.CreateUser(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (cont *Controller) CheckUser(c *gin.Context) {
	accessToken := c.Query("access_token")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token in request"})
		return
	}
	user, err := cont.services.GetUser(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)
}
