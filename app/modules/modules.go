package modules

import (
	"api-desatanggap/api"
	adminApi "api-desatanggap/api/admin"
	userApi "api-desatanggap/api/user"
	adminBusiness "api-desatanggap/business/admin"
	userBusiness "api-desatanggap/business/user"
	"api-desatanggap/config"
	adminRepo "api-desatanggap/repository/admin"
	userRepo "api-desatanggap/repository/user"
	"api-desatanggap/utils"
)

func RegistrationModules(dbCon *utils.DatabaseConnection, _ *config.AppConfig) api.Controller {
	userPermitRepository := userRepo.RepositoryFactory(dbCon)
	userPermitService := userBusiness.NewService(userPermitRepository)
	userPermitController := userApi.NewController(userPermitService)

	adminPermitRepository := adminRepo.RepositoryFactory(dbCon)
	adminPermitService := adminBusiness.NewService(adminPermitRepository)
	adminPermitController := adminApi.NewController(adminPermitService)
	controller := api.Controller{
		UserController:  userPermitController,
		AdminController: adminPermitController,
	}
	return controller
}
