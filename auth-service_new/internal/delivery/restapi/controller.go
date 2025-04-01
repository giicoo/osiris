package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/entity"
	"github.com/giicoo/osiris/auth-service/internal/entity/models"
	"github.com/giicoo/osiris/auth-service/internal/services"
	"github.com/giicoo/osiris/auth-service/pkg/apiError"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func GetUser(c *gin.Context) (*entity.User, apiError.AErr) {
	usr, _ := c.Get("user")
	user, ok := usr.(*entity.User)
	if !ok {
		return nil, apiError.ErrDontHaveUser
	}
	return user, nil
}

func (cont *Controller) CreateUser(c *gin.Context) {
	var json models.UserCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		Email:    json.Email,
		Password: json.Password,
	}
	userDB, aerr := cont.services.CreateUser(user)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}

	response := models.UserIdRequest{UserID: userDB.ID}
	c.JSON(200, response)
	return
}

func (cont *Controller) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, aerr := cont.services.GetUserByID(id)
	if aerr != nil {
		logrus.Error(err)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, response)
	return
}

func (cont *Controller) Login(c *gin.Context) {
	var json models.UserCheckRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &entity.User{
		Email:    json.Email,
		Password: json.Password,
	}
	sessionID, aerr := cont.services.CheckUser(user)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}

	response := models.SessionResponse{ID: sessionID}
	c.JSON(200, response)
	return
}

func (cont *Controller) Auth(c *gin.Context) {
	sessionID := c.Param("session_id")

	user, err := cont.services.GetUserFromSession(sessionID)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := models.UserResponse{UserID: user.ID, Email: user.Email}
	c.JSON(200, response)
}

func (cont *Controller) Logout(c *gin.Context) {
	var json models.SessionRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cont.services.DeleteSession(json.ID); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}
