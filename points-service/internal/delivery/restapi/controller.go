package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	"github.com/giicoo/osiris/points-service/internal/entity/models"
	"github.com/giicoo/osiris/points-service/internal/services"
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
//	@Summary	create point
//	@Schemes
//	@Description	create point
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.CreatePoint	true	"Write Point"
//	@Success		200		{object}	entity.Point
//	@Router			/create/point [post]
func (cont *Controller) CreatePoint(c *gin.Context) {
	var json models.CreatePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	point := &entity.Point{
		UserID:   json.UserID,
		Title:    json.Title,
		Location: json.Location,
		Radius:   json.Radius,
	}
	pointDB, err := cont.services.CreatePoint(point)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, pointDB)
	return
}

// Osiris godoc
//
//		@Summary	get point
//		@Schemes
//		@Description	get point by id
//		@Tags			points
//		@Accept			json
//		@Produce		json
//	 @Param          id   path      int  true  "Point ID"
//		@Success		200		{object}	entity.Point
//		@Router			/get/point/{id} [get]
func (cont *Controller) GetPoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Warn(id)
	response, err := cont.services.GetPoint(id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Warn(response)
	c.JSON(200, response)
	return
}

// Osiris godoc
//
//		@Summary	get points
//		@Schemes
//		@Description	get points by user_id
//		@Tags			points
//		@Accept			json
//		@Produce		json
//	 @Param          user_id   path      int  true  "User ID"
//		@Success		200		{object}	[]entity.Point
//		@Router			/get/points/{user_id} [get]
func (cont *Controller) GetPoints(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := cont.services.GetPoints(user_id)
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
//	@Summary	update title
//	@Schemes
//	@Description	update title
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.UpdateTitlePoint	true	"Write Title"
//	@Success		200		{object}	entity.Point
//	@Router			/update/point/title [put]
func (cont *Controller) UpdateTitle(c *gin.Context) {
	var json models.UpdateTitlePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := cont.services.UpdateTitle(json.ID, json.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	update location
//	@Schemes
//	@Description	update location
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.UpdateLocationPoint	true	"Write Location like '0 30'"
//	@Success		200		{object}	entity.Point
//	@Router			/update/point/location [put]
func (cont *Controller) UpdateLocation(c *gin.Context) {
	var json models.UpdateLocationPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := cont.services.UpdateLocation(json.ID, json.Location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	update radius
//	@Schemes
//	@Description	update radius
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.UpdateRadiusPoint	true	"Write Radius"
//	@Success		200		{object}	entity.Point
//	@Router			/update/point/radius [put]
func (cont *Controller) UpdateRadius(c *gin.Context) {
	var json models.UpdateRadiusPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := cont.services.UpdateRadius(json.ID, json.Radius)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
	return
}

// Osiris godoc
//
//	@Summary	delete point
//	@Schemes
//	@Description	delete point
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Param			point	body		models.DeletePoint	true	"Write ID"
//	@Success		200		{object}	string
//	@Router			/delete/point [delete]
func (cont *Controller) DeletePoint(c *gin.Context) {
	var json models.DeletePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cont.services.DeletePoint(json.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "successful deleted",
	})
	return
}
