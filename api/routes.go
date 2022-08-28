package api

import (
	"api-desatanggap/api/admin"
	"api-desatanggap/api/middleware"
	"api-desatanggap/api/user"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	UserController  *user.Controller
	AdminController *admin.Controller
}

func RegistrationPath(e *echo.Echo, controller Controller) {
	e.POST("/sendemail", controller.UserController.SmtpEmail)
	e.GET("/photo/:name", func(c echo.Context) error {
		name := fmt.Sprintf("utils/img/%s", c.Param("name"))
		fmt.Println(name)
		return c.Inline(name, name)
	})
	acc := e.Group("/users")
	acc.GET("/products", controller.UserController.GetProductByIdAccount)
	acc.GET("/products/status", controller.UserController.GetProductByIdAccStatus)
	acc.POST("/products/transaction", controller.UserController.PostProductTransaction)
	acc.GET("/products/transaction", controller.UserController.GetProductTransactionByIDUser)
	acc.POST("/login", controller.UserController.LoginAccount)
	acc.POST("/registrations", controller.UserController.RegisterAccount)
	acc.GET("/role", controller.UserController.GetRole)
	acc.POST("/products", controller.UserController.CreateProduct)
	acc.POST("/verification-account", controller.UserController.VerificationAccount)
	e.POST("/upload", controller.UserController.UploadFileHandle)
	acc.GET("/test", controller.UserController.GetRole, middleware.JWTMiddleware())
	acc.DELETE("/delete", controller.UserController.DeleteUser)
	acc.PUT("/send-verification", controller.UserController.SendVerification)
	// e.POST("/register", controller.UserController.RegisterUser)
	// e.GET("/User", controller.UserController.FindUser)
	admin := e.Group("/administrators")
	admin.GET("/products", controller.AdminController.FindProductByStatus)
	admin.POST("/products/:id/approve", controller.AdminController.UpdateStatusProduct)
	admin.POST("/cooperation", controller.AdminController.CreateCooperation)
	admin.GET("/role", controller.AdminController.GetRole)
	admin.POST("/registrations", controller.AdminController.RegisterAdmin)
	admin.POST("/login", controller.AdminController.LoginAdmin)

}
