package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strconv"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/entity"
	"github.com/sirupsen/logrus"
)

const (
	tokenURL    = "https://oauth.yandex.ru/token"
	userInfoURL = "https://login.yandex.ru/info"
)

type YandexAPI struct {
	cfg *config.Config
}

func NewYandexAPI(cfg *config.Config) *YandexAPI {
	return &YandexAPI{
		cfg: cfg,
	}
}

func (ya *YandexAPI) GetAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", ya.cfg.ClientID)
	data.Set("client_secret", ya.cfg.ClientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {

		return "", fmt.Errorf("get access token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp map[string]interface{}
	json.Unmarshal(body, &tokenResp)
	logrus.Info(tokenResp)
	_, ok := tokenResp["access_token"]
	if ok {
		accessToken := tokenResp["access_token"].(string)
		return accessToken, nil
	}
	erro := tokenResp["error"].(string)
	return "", fmt.Errorf("yandex code: %s", erro)
}

func (ya *YandexAPI) GetUserInfo(accessToken string) (*entity.User, error) {
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("get user info: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo map[string]string
	json.Unmarshal(body, &userInfo)
	id, err := strconv.Atoi(userInfo["id"])
	if err != nil {
		return nil, fmt.Errorf("id to int: %w", err)
	}
	user := &entity.User{
		ID:          id,
		LoginYandex: userInfo["login"],
		FirstName:   userInfo["first_name"],
		LastName:    userInfo["last_name"],
	}
	return user, nil
}
