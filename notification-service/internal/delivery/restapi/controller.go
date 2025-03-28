package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/notification-service/internal/config"
	"github.com/giicoo/osiris/notification-service/internal/entity"
	"github.com/giicoo/osiris/notification-service/internal/services"
	"github.com/giicoo/osiris/notification-service/pkg/apiError"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (cont *Controller) WsConnect(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Info("try to connect")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cont.services.CreateSession(id, conn)
}
