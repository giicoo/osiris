package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	"github.com/giicoo/osiris/alerts-service/internal/entity/models"
	"github.com/giicoo/osiris/alerts-service/internal/services"
	"github.com/giicoo/osiris/alerts-service/pkg/apiError"
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

// @Title Create Alert
// @Param  alert  body  models.CreateAlert  true  "Info of a alert"
// @Success  200  object  entity.Alert           "Alert JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource alerts
// @Route /create/alert [post]
func (cont *Controller) CreateAlert(c *gin.Context) {
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	var json models.CreateAlert
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert := &entity.Alert{
		UserID:      user.ID,
		Title:       json.Title,
		Description: json.Description,
		TypeID:      json.TypeID,
		Location:    json.Location,
		Radius:      json.Radius,
		Status:      json.Status,
	}
	alertID, aerr := cont.services.CreateAlert(alert)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, alertID)
	return
}

// @Title Stop Alert
// @Param  alert  body  models.StopAlert  true  "Info of a alert"
// @Success  200  object  entity.Response      "Success JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource alerts
// @Route /stop/alert [post]
func (cont *Controller) StopAlert(c *gin.Context) {
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	var json models.StopAlert
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aerr = cont.services.StopAlert(json.ID, user.ID)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "successful stop",
	})
	return
}

// @Title Create Type
// @Param  type  body  models.CreateType  true  "Info of a type"
// @Success  200  object  entity.Type           "Type JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource types
// @Route /create/type [post]
func (cont *Controller) CreateType(c *gin.Context) {
	// user, aerr := GetUser(c)
	// if aerr != nil {
	// 	logrus.Error(aerr)
	// 	c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
	// 	return
	// }
	var json models.CreateType
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	typeModel := &entity.Type{
		Title: json.Title,
	}
	typeID, aerr := cont.services.CreateType(typeModel)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, typeID)
	return
}

// @Title Get Alert
// @Param  alertID  path  int  true  "ID alert" "1"
// @Success  200  object  entity.Alert           "Type JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource alerts
// @Route /get/alert/{alertID} [get]
func (cont *Controller) GetAlert(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, aerr := cont.services.GetAlert(id)
	if aerr != nil {
		logrus.Error(err)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// @Title Get Type
// @Param  typeID  path  int  true  "ID type" "1"
// @Success  200  object  entity.Type           "Type JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource types
// @Route /get/type/{typeID} [get]
func (cont *Controller) GetType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, aerr := cont.services.GetType(id)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// @Title Get Alerts
// @Success  200  object  []entity.Alert           "Alerts JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource alerts
// @Route /get/alerts [get]
func (cont *Controller) GetAlerts(c *gin.Context) {
	response, aerr := cont.services.GetAlerts()
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// @Title Get Types
// @Success  200  object  []entity.Type           "Types JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource types
// @Route /get/types [get]
func (cont *Controller) GetTypes(c *gin.Context) {
	response, aerr := cont.services.GetTypes()
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// @Title Delete Type
// @Param  type  body  models.DeleteType  true  "Info of a type"
// @Success  200  object  entity.Response           "]JSON"
// @Failure  400  object  entity.Error  "ErrorResponse JSON"
// @Resource types
// @Route /delete/type [delete]
func (cont *Controller) DeleteType(c *gin.Context) {
	var json models.DeleteType
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if aerr := cont.services.DeleteType(json.ID); aerr != nil {
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "successful deleted",
	})
	return
}
