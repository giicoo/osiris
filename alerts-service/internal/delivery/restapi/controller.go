package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	"github.com/giicoo/osiris/alerts-service/internal/entity/models"
	"github.com/giicoo/osiris/alerts-service/internal/services"
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

// Osiris godoc
//
//	@Summary	create alert
//	@Schemes
//	@Description	create alert
//	@Tags			alerts
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.CreateAlert	true	"Write Alert"
//	@Success		200		{object}	entity.Alert
//	@Router			/create/alert [post]
func (cont *Controller) CreateAlert(c *gin.Context) {
	var json models.CreateAlert
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert := &entity.Alert{
		UserID:      json.UserID,
		Title:       json.Title,
		Description: json.Description,
		TypeID:      json.TypeID,
		Location:    json.Location,
		Radius:      json.Radius,
		Status:      json.Status,
	}
	alertID, err := cont.services.CreateAlert(alert)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, alertID)
	return
}

// Osiris godoc
//
//	@Summary	create type
//	@Schemes
//	@Description	create type
//	@Tags			types
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.CreateType	true	"Write Type"
//	@Success		200		{object}	entity.Type
//	@Router			/create/type [post]
func (cont *Controller) CreateType(c *gin.Context) {
	var json models.CreateType
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	typeModel := &entity.Type{
		Title: json.Title,
	}
	typeID, err := cont.services.CreateType(typeModel)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, typeID)
	return
}

// Osiris godoc
//
//		@Summary	get alert
//		@Schemes
//		@Description	get alert by id
//		@Tags			alerts
//		@Accept			json
//		@Produce		json
//	 @Param          id   path      int  true  "Alert ID"
//		@Success		200		{object}	entity.Alert
//		@Router			/get/alert/{id} [get]
func (cont *Controller) GetAlert(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := cont.services.GetAlert(id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// Osiris godoc
//
//		@Summary	get type
//		@Schemes
//		@Description	get type by id
//		@Tags			types
//		@Accept			json
//		@Produce		json
//	 @Param          id   path      int  true  "Type ID"
//		@Success		200		{object}	entity.Type
//		@Router			/get/type/{id} [get]
func (cont *Controller) GetType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := cont.services.GetType(id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	get alerts
//	@Schemes
//	@Description	get alerts
//	@Tags			alerts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]entity.Alert
//	@Router			/get/alerts [get]
func (cont *Controller) GetAlerts(c *gin.Context) {
	response, err := cont.services.GetAlerts()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	get types
//	@Schemes
//	@Description	get types
//	@Tags			types
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]entity.Type
//	@Router			/get/types [get]
func (cont *Controller) GetTypes(c *gin.Context) {
	response, err := cont.services.GetTypes()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	delete type
//	@Schemes
//	@Description	delete type
//	@Tags			types
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.DeleteType	true	"Write ID"
//	@Success		200		{object}	string
//	@Router			/delete/type [delete]
func (cont *Controller) DeleteType(c *gin.Context) {
	var json models.DeleteType
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cont.services.DeleteType(json.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "successful deleted",
	})
	return
}
