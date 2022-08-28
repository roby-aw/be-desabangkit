package admin

import (
	adminBusiness "api-desatanggap/business/admin"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service adminBusiness.Service
}

func NewController(service adminBusiness.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (Controller *Controller) RegisterAdmin(c echo.Context) error {
	Data := adminBusiness.RegAdmin{}
	c.Bind(&Data)
	_, err := Controller.service.CreateAdmin(&Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success register account",
	})
}

func (Controller *Controller) LoginAdmin(c echo.Context) error {
	Data := adminBusiness.AuthLogin{}
	c.Bind(&Data)
	result, err := Controller.service.LoginAdmin(&Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success login",
		"data":     result,
	})
}
func (Controller Controller) GetRole(c echo.Context) error {
	result, err := Controller.service.GetRole()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get role",
		"data":     result,
	})
}

func (Controller *Controller) CreateCooperation(c echo.Context) error {
	Data := adminBusiness.RegCooperation{}
	c.Bind(&Data)
	_, err := Controller.service.CreateCooperation(&Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success create cooperation",
	})
}

func (Controller *Controller) FindProductByStatus(c echo.Context) error {
	var preorder *bool
	if c.QueryParam("preorder") != "" {
		preorder1, _ := strconv.ParseBool(c.QueryParam("preorder"))
		preorder = &preorder1
	}
	result, err := Controller.service.GetProductByStatus(preorder)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get product",
		"data":     result,
	})
}

func (Controller *Controller) UpdateStatusProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": "id is required",
		})
	}
	err := Controller.service.UpdateStatusProduct(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success update status product",
	})
}
