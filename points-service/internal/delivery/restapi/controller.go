package restapi

import (
	"net/http"
	"strconv"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	"github.com/giicoo/osiris/points-service/internal/entity/models"
	"github.com/giicoo/osiris/points-service/internal/services"
	"github.com/giicoo/osiris/points-service/pkg/apiError"
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	var json models.CreatePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	point := &entity.Point{
		UserID:   user.ID,
		Title:    json.Title,
		Location: json.Location,
		Radius:   json.Radius,
	}
	pointDB, aerr := cont.services.CreatePoint(point)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, aerr := cont.services.GetPoint(id, user.ID)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	response, aerr := cont.services.GetPoints(user.ID)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	//TODO: изменение точки только принадлежащему юзеру
	var json models.UpdateTitlePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, aerr := cont.services.UpdateTitle(json.ID, json.Title, user.ID)
	if aerr != nil {
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	var json models.UpdateLocationPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, aerr := cont.services.UpdateLocation(json.ID, json.Location, user.ID)
	if aerr != nil {
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}
	var json models.UpdateRadiusPoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, aerr := cont.services.UpdateRadius(json.ID, json.Radius, user.ID)
	if aerr != nil {
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
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
	user, aerr := GetUser(c)
	if aerr != nil {
		logrus.Error(aerr)
		c.JSON(http.StatusBadRequest, gin.H{"error": aerr.Error()})
		return
	}
	var json models.DeletePoint
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if aerr := cont.services.DeletePoint(json.ID, user.ID); aerr != nil {
		c.JSON(aerr.Code(), gin.H{"error": aerr.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "successful deleted",
	})
	return
}
