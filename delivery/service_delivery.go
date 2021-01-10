package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/interfaces/service"
	"github.com/naysudes/technopark-db-forum/tools"
	"net/http"
)

type ServiceDelivery struct {
	serviceUseCase service.Usecase
}

func NewServiceHandler(e *echo.Echo, usecase service.Usecase) *ServiceDelivery {
	sh := &ServiceDelivery{
		serviceUseCase: usecase,
	}
	e.GET("/api/service/status", sh.GetStatus())
	e.POST("/api/service/clear", sh.Clear())
	return sh
}

func (delivery *ServiceDelivery) GetStatus() echo.HandlerFunc {
	return func(context echo.Context) error {
		stat, err := delivery.serviceUseCase.GetStatus()
		if err != nil {
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, stat)
	}
}

func (delivery *ServiceDelivery) Clear() echo.HandlerFunc {
	return func(context echo.Context) error {
		err := delivery.serviceUseCase.DeleteAll()
		if err != nil {
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.NoContent(http.StatusOK)
	}
}
