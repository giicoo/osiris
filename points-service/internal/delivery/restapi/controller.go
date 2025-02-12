package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	"github.com/giicoo/osiris/points-service/internal/entity/models"
	"github.com/giicoo/osiris/points-service/internal/services"
	"github.com/gin-gonic/gin"
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

func (cont *Controller) CreatePoint(c *gin.Context) {
	var json models.CreatePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, json)
}

func (cont *Controller) GetPoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	response := entity.Point{ID: id}
	c.JSON(200, response)
}

func (cont *Controller) GetPoints(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Param("user_id"))

	response := entity.Point{UserID: user_id}
	c.JSON(200, response)
}

func (cont *Controller) UpdateTitle(c *gin.Context) {
	var json models.UpdateTitlePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, json)
}

func (cont *Controller) UpdateLocation(c *gin.Context) {
	var json models.UpdateLocationPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, json)
}

func (cont *Controller) UpdateRadius(c *gin.Context) {
	var json models.UpdateRadiusPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, json)
}

func (cont *Controller) DeletePoint(c *gin.Context) {
	var json models.DeletePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, json)
}
