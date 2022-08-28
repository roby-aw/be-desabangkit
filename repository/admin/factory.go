package admin

import (
	"api-desatanggap/business/admin"
	"api-desatanggap/utils"
)

func RepositoryFactory(dbCon *utils.DatabaseConnection) admin.Repository {
	adminRepo := NewMongoRepository(dbCon.MongoDB)
	return adminRepo
}
