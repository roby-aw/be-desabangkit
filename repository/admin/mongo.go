package admin

import (
	"api-desatanggap/business/admin"
	"api-desatanggap/repository"
	"api-desatanggap/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository struct {
	col     *mongo.Collection
	colRole *mongo.Collection
	colCoop *mongo.Collection
	colProd *mongo.Collection
}

func NewMongoRepository(col *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		col:     col.Collection("admin"),
		colRole: col.Collection("roles_admin"),
		colCoop: col.Collection("cooperation"),
		colProd: col.Collection("products"),
	}
}

func (repo *MongoDBRepository) CreateAdmin(Data *admin.RegAdmin) (*admin.Admin, error) {
	hashpw, _ := utils.Hash(Data.Password)
	ObjId_userid, _ := primitive.ObjectIDFromHex(Data.Role_id)
	InsertData := &repository.Admin{
		Username: Data.Username,
		Fullname: Data.Fullname,
		Password: string(hashpw),
		Role_id:  ObjId_userid,
	}
	result, err := repo.col.InsertOne(context.Background(), InsertData)
	var tmpAdmin admin.Admin
	data, _ := json.Marshal(Data)
	err = json.Unmarshal(data, &tmpAdmin)
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", result.InsertedID))
	tmpAdmin.ID = id
	return &tmpAdmin, nil
}

func (repo *MongoDBRepository) FindAdminByUsername(username string) (*admin.Admin, error) {
	filter := bson.A{
		bson.M{
			"$match": bson.M{
				"username": bson.M{"$regex": username},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "roles_admin",
				"localField":   "role_id",
				"foreignField": "_id",
				"as":           "roles",
			},
		},
	}
	cur, err := repo.col.Aggregate(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var tmpAdmin []admin.Admin
	if err = cur.All(context.Background(), &tmpAdmin); err != nil {
		return nil, err
	}
	if len(tmpAdmin) < 1 {
		return nil, errors.New("Data Not Found")
	}
	return &tmpAdmin[0], nil
}

func (repo *MongoDBRepository) CreateToken(Data *admin.Admin) (*string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &admin.Claims{
		Username: Data.Username,
		Role:     Data.Roles[0].Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	SECRET_KEY := os.Getenv("SECRET_JWT")
	token_jwt, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return nil, err
	}
	return &token_jwt, err
}

func (repo *MongoDBRepository) GetRole() ([]*admin.Role, error) {
	var Role []*admin.Role
	cur, err := repo.colRole.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	cur.All(context.Background(), &Role)
	return Role, nil
}

func (repo *MongoDBRepository) CreateCooperation(Data *admin.RegCooperation) (*admin.Cooperation, error) {
	InsertData := &admin.Cooperation{
		Name:       Data.Name,
		Address:    Data.Address,
		City:       Data.City,
		Province:   Data.Province,
		Latitude:   Data.Latitude,
		Longitude:  Data.Longitude,
		Photo_url:  Data.Photo_url,
		Email:      Data.Email,
		Username:   Data.Username,
		Password:   Data.Password,
		Created_at: time.Now(),
	}
	result, err := repo.colCoop.InsertOne(context.Background(), InsertData)
	var tmpCooperation admin.Cooperation
	data, _ := json.Marshal(Data)
	err = json.Unmarshal(data, &tmpCooperation)
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", result.InsertedID))
	tmpCooperation.ID = id
	return &tmpCooperation, nil
}

func (repo *MongoDBRepository) GetProductByStatus(preorder *bool) ([]admin.Product, error) {
	var Product []admin.Product
	filter := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userid",
				"foreignField": "_id",
				"as":           "account",
			},
		},
	}
	if preorder != nil {
		fmt.Println("preorder not nil")
		filter1 := bson.M{
			"$match": bson.M{
				"is_preorder": preorder,
			},
		}
		filter = append(filter, filter1)
	}
	fmt.Println(filter)

	cur, err := repo.colProd.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem admin.Product
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		Product = append(Product, elem)
	}
	return Product, nil
}

func (repo *MongoDBRepository) UpdateStatusProduct(id string) error {
	ObjId_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": ObjId_id,
	}
	update := bson.M{
		"$set": bson.M{
			"is_approved": true,
		},
	}
	_, err := repo.colProd.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
