package admin

import (
	"api-desatanggap/utils"
	"errors"

	"github.com/go-playground/validator/v10"
)

type Repository interface {
	CreateAdmin(Data *RegAdmin) (*Admin, error)
	FindAdminByUsername(username string) (*Admin, error)
	CreateToken(Data *Admin) (*string, error)
	GetRole() ([]*Role, error)
	CreateCooperation(Data *RegCooperation) (*Cooperation, error)
	GetProductByStatus(preorder *bool) ([]Product, error)
	UpdateStatusProduct(id string) error
}
type Service interface {
	CreateAdmin(Data *RegAdmin) (*Admin, error)
	FindAdminByUsername(username string) (*Admin, error)
	LoginAdmin(auth *AuthLogin) (*ResponseLogin, error)
	GetRole() ([]*Role, error)
	CreateCooperation(Data *RegCooperation) (*Cooperation, error)
	GetProductByStatus(preorder *bool) ([]Product, error)
	UpdateStatusProduct(id string) error
}

type service struct {
	repository Repository
	validate   *validator.Validate
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
		validate:   validator.New(),
	}
}

func (s *service) CreateAdmin(Data *RegAdmin) (*Admin, error) {
	return s.repository.CreateAdmin(Data)
}

func (s *service) FindAdminByUsername(username string) (*Admin, error) {
	return s.repository.FindAdminByUsername(username)
}

func (s *service) LoginAdmin(auth *AuthLogin) (*ResponseLogin, error) {
	err := s.validate.Struct(auth)
	if err != nil {
		return nil, err
	}
	data, err := s.repository.FindAdminByUsername(auth.Username)
	if err != nil {
		return nil, errors.New("wrong username")
	}
	err = utils.VerifyPassword(data.Password, auth.Password)
	if err != nil {
		return nil, errors.New("wrong password")
	}
	token, _ := s.repository.CreateToken(data)
	res := &ResponseLogin{
		Admin: *data,
		Token: *token,
	}
	return res, nil
}

func (s *service) GetRole() ([]*Role, error) {
	return s.repository.GetRole()
}

func (s *service) CreateCooperation(Data *RegCooperation) (*Cooperation, error) {
	return s.repository.CreateCooperation(Data)
}

func (s *service) GetProductByStatus(preorder *bool) ([]Product, error) {
	return s.repository.GetProductByStatus(preorder)
}

func (s *service) UpdateStatusProduct(id string) error {
	return s.repository.UpdateStatusProduct(id)
}
