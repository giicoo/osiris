package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

	// Получение access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", cont.cfg.ClientID)
	data.Set("client_secret", cont.cfg.ClientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResp map[string]interface{}
	json.Unmarshal(body, &tokenResp)
	accessToken := tokenResp["access_token"].(string)

	// Получение информации о пользователе
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)
	c.JSON(http.StatusOK, userInfo)
}
