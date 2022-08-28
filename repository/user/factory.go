package user

import (
	"api-desatanggap/business/user"
	"api-desatanggap/utils"
)

func RepositoryFactory(dbCon *utils.DatabaseConnection) user.Repository {
	userRepo := NewMongoRepository(dbCon.MongoDB)
	return userRepo
}
