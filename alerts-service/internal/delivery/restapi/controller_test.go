package restapi

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	mock_rabbitmq "github.com/giicoo/osiris/alerts-service/internal/infrastructure/mock"
	mock_repository "github.com/giicoo/osiris/alerts-service/internal/repository/mock"
	"github.com/giicoo/osiris/alerts-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePoint(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)
	type mockBehaviorRabbitMq func(r *mock_rabbitmq.MockAlertProducing)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehaviorRepo     mockBehaviorRepo
		mockBehaviorRabbitMq mockBehaviorRabbitMq
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
			}`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateAlert(gomock.Any()).Return(1, nil).AnyTimes()
				r.EXPECT().GetAlert(1).Return(&entity.Alert{ID: 1, UserID: 1}, nil)
			},
			mockBehaviorRabbitMq: func(r *mock_rabbitmq.MockAlertProducing) {
				r.EXPECT().PublicMessage(gomock.Any()).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":1,"title":"","description":"","type_id":0,"location":"","radius":0,"status":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name: "JSON err",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
			`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateAlert(gomock.Any()).Return(1, nil).AnyTimes()
				r.EXPECT().GetAlert(1).Return(&entity.Alert{ID: 1, UserID: 1}, nil).AnyTimes()
			},
			mockBehaviorRabbitMq: func(r *mock_rabbitmq.MockAlertProducing) {
				r.EXPECT().PublicMessage(gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name: "DB err",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
				}
			`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateAlert(gomock.Any()).Return(0, errors.New("test")).AnyTimes()
				r.EXPECT().GetAlert(1).Return(&entity.Alert{ID: 1, UserID: 1}, nil).AnyTimes()
			},
			mockBehaviorRabbitMq: func(r *mock_rabbitmq.MockAlertProducing) {
				r.EXPECT().PublicMessage(gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service create alert: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)
			test.mockBehaviorRabbitMq(rabbitmq)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create/alert",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.POST("/create/alert", controller.CreateAlert)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestCreateType(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
			}`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateType(gomock.Any()).Return(1, nil).AnyTimes()
				r.EXPECT().GetType(1).Return(&entity.Type{ID: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":""}`,
		},
		{
			name: "JSON err",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
			`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateType(gomock.Any()).Return(1, nil).AnyTimes()
				r.EXPECT().GetType(1).Return(&entity.Type{ID: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name: "DB err",
			inputBody: `{
				"title":"test",
				"description":"test",
				"type_id":1,
				"location":"test",
				"radius":1,
				"status":true
			}`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().CreateType(gomock.Any()).Return(1, errors.New("test")).AnyTimes()
				r.EXPECT().GetType(1).Return(&entity.Type{ID: 1}, nil).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service create type: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create/type",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.POST("/create/type", controller.CreateType)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetAlert(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetAlert(1).Return(&entity.Alert{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"user_id":0,"title":"","description":"","type_id":0,"location":"","radius":0,"status":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name: "DB err",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetAlert(1).Return(&entity.Alert{}, errors.New("test"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service get alert: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/alert/1", nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/alert/:id", controller.GetAlert)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetType(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetType(1).Return(&entity.Type{}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"title":""}`,
		},
		{
			name: "DB err",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetType(1).Return(&entity.Type{}, errors.New("test"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service get type: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/type/1", nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/type/:id", controller.GetType)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetAlerts(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetAlerts().Return([]*entity.Alert{&entity.Alert{}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"user_id":0,"title":"","description":"","type_id":0,"location":"","radius":0,"status":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "DB err",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetAlerts().Return([]*entity.Alert{&entity.Alert{}}, errors.New("test"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service get alerts: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/alerts", nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/alerts", controller.GetAlerts)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestGetTypes(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetTypes().Return([]*entity.Type{&entity.Type{}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":""}]`,
		},
		{
			name: "DB err",
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().GetTypes().Return([]*entity.Type{&entity.Type{}}, errors.New("test"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service get types: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/types", nil)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.GET("/get/types", controller.GetTypes)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}

func TestDeleteType(t *testing.T) {
	type mockBehaviorRepo func(r *mock_repository.MockRepo)

	tests := []struct {
		name                 string
		inputBody            string
		mockBehaviorRepo     mockBehaviorRepo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"id":1
			}`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().DeleteType(1).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"successful deleted"}`,
		},
		{
			name: "JSON err",
			inputBody: `{
				"id":1
			`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().DeleteType(1).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
		},
		{
			name: "DB err",
			inputBody: `{
				"id":1
			}`,
			mockBehaviorRepo: func(r *mock_repository.MockRepo) {
				r.EXPECT().DeleteType(1).Return(errors.New("test")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service delete type: test"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepo(c)
			test.mockBehaviorRepo(repo)

			rabbitmq := mock_rabbitmq.NewMockAlertProducing(c)

			cfg := &config.Config{}
			services := services.NewServices(cfg, repo, rabbitmq)
			controller := NewController(cfg, services)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete/type",
				bytes.NewBufferString(test.inputBody))

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				user := &entity.User{ID: 1}
				ctx.Set("user", user)
				ctx.Next()
			})
			r.DELETE("/delete/type", controller.DeleteType)
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}
