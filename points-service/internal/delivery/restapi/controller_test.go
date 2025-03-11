package restapi

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	mock_repository "github.com/giicoo/osiris/points-service/internal/repository/mock"
	"github.com/giicoo/osiris/points-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePoint(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"title":"test",
				"location":"POINT(0 0)",
				"radius":1
			}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreatePoint(gomock.Any()).Return(1, nil).AnyTimes()
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name: "DB err",
			inputBody: `{
				"title":"test",
				"location":"POINT(0 0)",
				"radius":1
			}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreatePoint(gomock.Any()).Return(0, errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service create point: test"}`,
		},
		{
			name:      "JSON err",
			inputBody: `{}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreatePoint(gomock.Any()).Return(0, errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Key: 'CreatePoint.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'CreatePoint.Location' Error:Field validation for 'Location' failed on the 'required' tag\nKey: 'CreatePoint.Radius' Error:Field validation for 'Radius' failed on the 'required' tag"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create/point",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.POST("/create/point", controller.CreatePoint)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetPoint(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: 1,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "Point not by user",
			inputBody: 1,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 0, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"point not by user with 1 id"}`,
		},
		{
			name:      "DB err",
			inputBody: 1,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(nil, errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service get point: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/get/point/%d", test.inputBody), nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/point/:id", controller.GetPoint)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetPoints(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoints(gomock.Any()).Return([]*entity.Point{&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, &entity.Point{ID: 2, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":2,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/points", nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/points", controller.GetPoints)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestDeletePoint(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().DeletePoint(gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"successful deleted"}`,
		},
		{
			name:      "JSON err",
			inputBody: `{"id":1`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().DeletePoint(gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name:      "DB err",
			inputBody: `{"id":1}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().DeletePoint(gomock.Any()).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service delete point: test"}`,
		},
		{
			name:      "point not by user",
			inputBody: `{"id":1}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 0, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().DeletePoint(gomock.Any()).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"point not by user with 1 id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete/point",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.DELETE("/delete/point", controller.DeletePoint)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestUpdateTitle(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"title":"test"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateTitle(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "point no by user",
			inputBody: `{"id":1,"title":"test"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 0, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateTitle(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"point not by user with 1 id"}`,
		},
		{
			name:      "JSON err",
			inputBody: `{"id":1,"title":"test"`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateTitle(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name:      "DB err",
			inputBody: `{"id":1,"title":"test"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateTitle(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service update point: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/update/point/title",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.PUT("/update/point/title", controller.UpdateTitle)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestUpdateLocation(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"location":"POINT(0 0)"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "point no by user",
			inputBody: `{"id":1,"location":"POINT(0 0)"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 0, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"point not by user with 1 id"}`,
		},
		{
			name:      "JSON err",
			inputBody: `{"id":1,"location":"POINT(0 0)"`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name:      "DB err",
			inputBody: `{"id":1,"location":"POINT(0 0)"}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service update point: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/update/point/location",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.PUT("/update/point/location", controller.UpdateLocation)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestUpdateRadius(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"radius":2}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateRadius(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"test","location":"POINT(0 0)","radius":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "point no by user",
			inputBody: `{"id":1,"radius":2}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 0, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateRadius(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"point not by user with 1 id"}`,
		},
		{
			name:      "JSON err",
			inputBody: `{"id":1,"radius":2`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateRadius(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name:      "DB err",
			inputBody: `{"id":1,"radius":2}`,
			mockBehavior: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetPoint(gomock.Any()).Return(&entity.Point{ID: 1, UserID: 1, Title: "test", Location: "POINT(0 0)", Radius: 1}, nil).AnyTimes()
				r.EXPECT().UpdateRadius(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service update point: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehavior(repo)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/update/point/radius",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.PUT("/update/point/radius", controller.UpdateRadius)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}
