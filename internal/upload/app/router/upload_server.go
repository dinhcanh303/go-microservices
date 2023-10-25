package router

import (
	"github.com/dinhcanh303/go-microservices/internal/upload/app/handlers"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func NewConfigureRoutes(uc uploads.UseCase, echo *echo.Echo) error {
	uploadHandler := handlers.NewUploadHandler(uc)
	echo.GET("/attachments/:id", uploadHandler.GetAttachment)
	echo.POST("/upload", uploadHandler.UploadFile)
	echo.DELETE("/attachments/:id", uploadHandler.DeleteAttachment)
	echo.PUT("/attachments/:id", uploadHandler.UpdateAttachment)
	return nil
}

var ConfigureRoutesSet = wire.NewSet(NewConfigureRoutes)
