package user

import (
	"api-desatanggap/business/user"
	"api-desatanggap/repository"
	"api-desatanggap/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository struct {
	col          *mongo.Collection
	colRole      *mongo.Collection
	colCode      *mongo.Collection
	colProd      *mongo.Collection
	colProdTrans *mongo.Collection
}

func NewMongoRepository(col *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		col:          col.Collection("users"),
		colRole:      col.Collection("roles_user"),
		colCode:      col.Collection("code_otp"),
		colProd:      col.Collection("products"),
		colProdTrans: col.Collection("products_transaction"),
	}
}

func (repo *MongoDBRepository) FindAccountByEmail(email string) (*user.Account, error) {
	filter := bson.A{
		bson.M{
			"$match": bson.M{
				"email": bson.M{"$regex": email},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "roles_user",
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
	var tmpAccount []user.Account
	if err = cur.All(context.Background(), &tmpAccount); err != nil {
		return nil, err
	}
	if len(tmpAccount) == 0 {
		return nil, errors.New("Data Not Found")
	}
	return &tmpAccount[0], nil
}

func (repo *MongoDBRepository) CreateAccount(Data *user.RegAccount) (*user.Account, error) {
	hashpw, _ := utils.Hash(Data.Password)
	ObjId_userid, err := primitive.ObjectIDFromHex(Data.Role_id)
	if err != nil {
		return nil, err
	}
	InsertData := &repository.Account{
		Email:      Data.Email,
		Fullname:   Data.Fullname,
		Password:   string(hashpw),
		Role_id:    ObjId_userid,
		IsVerified: false,
	}
	result, err := repo.col.InsertOne(context.Background(), InsertData)
	if err != nil {
		return nil, err
	}
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", result.InsertedID))

	InsertData.ID = id

	ResponseAccount := &user.Account{
		ID:         id,
		Email:      InsertData.Email,
		Fullname:   InsertData.Fullname,
		Password:   InsertData.Password,
		Role_id:    InsertData.Role_id,
		IsVerified: InsertData.IsVerified,
	}
	return ResponseAccount, nil
}

func (repo *MongoDBRepository) CreateToken(Data *user.Account) (*string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &user.Claims{
		Email: Data.Email,
		Role:  Data.Roles[0].Name,
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

func (repo *MongoDBRepository) Createcustomer(Data *user.Regcustomer) (*user.Regcustomer, error) {
	return nil, nil
}

func (repo *MongoDBRepository) Findcustomer() ([]user.Customer, error) {
	return nil, nil
}

func (repo *MongoDBRepository) GetRole() ([]*user.Role, error) {
	var Role []*user.Role
	cur, err := repo.colRole.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	cur.All(context.Background(), &Role)
	return Role, err
}

func (repo *MongoDBRepository) SendVerification(email string) error {
	var tmpAcc repository.Account
	err := repo.col.FindOne(context.Background(), bson.M{"email": email}).Decode(&tmpAcc)
	if err != nil {
		return errors.New("Email Not Registered")
	}
	codeotp, err := utils.InitEmail(email)
	if err != nil {
		return err
	}
	err = repo.CreateCodeOtp(email, codeotp)
	return nil
}

func (repo *MongoDBRepository) ValidationEmail(Data string) error {
	return nil
}

func (repo *MongoDBRepository) CreateCodeOtp(email string, codeotp string) error {
	timeExpired := time.Now().Add(24 * time.Hour)
	InsertCode := &repository.CodeOtp{
		Email:      email,
		Code:       codeotp,
		Expired_at: timeExpired,
	}
	var tmpAcc repository.CodeOtp
	repo.colCode.FindOne(context.Background(), bson.M{"email": email}).Decode(&tmpAcc)
	if tmpAcc.Email != "" {
		filter := bson.M{"email": email}
		update := bson.M{
			"$set": bson.M{
				"code":       codeotp,
				"expired_at": timeExpired,
			},
		}
		repo.colCode.UpdateOne(context.Background(), filter, update)
		return nil
	}
	_, err := repo.colCode.InsertOne(context.Background(), InsertCode)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoDBRepository) VerificationAccount(code string) error {
	var codeotp repository.CodeOtp
	err := repo.colCode.FindOne(context.Background(), bson.M{"code": code}).Decode(&codeotp)
	if err != nil {
		return errors.New("Code Not Found")
	}
	if codeotp.Expired_at.Before(time.Now()) {
		return errors.New("Code Expired")
	}
	repo.colCode.DeleteOne(context.Background(), bson.M{"code": code})
	filter := bson.M{"email": codeotp.Email}
	update := bson.M{"isverified": true}
	_, err = repo.col.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoDBRepository) DeleteUser(email string) error {
	filter := bson.M{"email": email}
	err := repo.col.FindOneAndDelete(context.Background(), filter).Err()
	if err != nil {
		return errors.New("Data Not Found")
	}
	return nil
}

func (repo *MongoDBRepository) InputProduct(Data *user.InputProduct, preorder string) error {
	CooperationID, err := primitive.ObjectIDFromHex(Data.Cooperationid)
	if err != nil {
		return err
	}
	UserID, err := primitive.ObjectIDFromHex(Data.UserID)
	if err != nil {
		return err
	}
	var IsPreorder bool
	if preorder == "true" {
		IsPreorder = true
	} else {
		IsPreorder = false
	}

	random := utils.RandomCapitalNumber(14)

	insertProd := &repository.Product{
		SKU:            *random,
		Photo_url:      Data.Photo_url,
		Name:           Data.Name,
		PriceExpected:  Data.PriceExpected,
		Quantity:       Data.Quantity,
		Latitude:       Data.Latitude,
		Longitude:      Data.Longitude,
		Category:       Data.Category,
		DeliveryOption: Data.DeliveryOption,
		Cooperationid:  CooperationID,
		UserID:         UserID,
		UserAddress:    Data.UserAddress,
		IsPreorder:     IsPreorder,
		IsApproved:     false,
		Created_at:     time.Now(),
	}
	repo.colProd.InsertOne(context.Background(), &insertProd)
	return nil
}

func (repo *MongoDBRepository) GetProductByIdAccount(id string) ([]user.Product, error) {
	var Product []user.Product
	UserID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.A{
		bson.M{
			"$match": bson.M{
				"userid": UserID,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userid",
				"foreignField": "_id",
				"as":           "account",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "cooperation",
				"localField":   "cooperationid",
				"foreignField": "_id",
				"as":           "cooperations",
			},
		},
		bson.M{
			"$unwind": "$account",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "roles_user",
				"localField":   "account.role_id",
				"foreignField": "_id",
				"as":           "account.roles",
			},
		},
	}
	fmt.Println(id, UserID)

	cur, err := repo.colProd.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem user.Product
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		Product = append(Product, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return Product, err
}

func (repo *MongoDBRepository) GetProductByAccStatus(approved *bool, verified *bool, id string) ([]user.Product, error) {
	var Product []user.Product
	UserID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.A{
		bson.M{
			"$match": bson.M{
				"userid": UserID,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userid",
				"foreignField": "_id",
				"as":           "account",
			},
		},
		bson.M{
			"$unwind": "$account",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "roles_user",
				"localField":   "account.role_id",
				"foreignField": "_id",
				"as":           "account.roles",
			},
		},
	}
	if approved != nil {
		filter = append(filter, bson.M{
			"$match": bson.M{
				"is_approved": *approved,
			},
		})
	}
	if verified != nil {
		filter = append(filter, bson.M{
			"$match": bson.M{
				"is_preorder": *verified,
			},
		})
	}

	cur, err := repo.colProd.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem user.Product
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		Product = append(Product, elem)
	}
	return Product, nil
}

func (repo *MongoDBRepository) InsertProductTransaction(InsertProduct *user.InputProductTransaction) error {
	UserID, err := primitive.ObjectIDFromHex(InsertProduct.Userid)
	if err != nil {
		return err
	}
	ProductID, err := primitive.ObjectIDFromHex(InsertProduct.Productid)
	if err != nil {
		return err
	}
	insertProd := &repository.InputProductTransaction{
		Userid:     UserID,
		Productid:  ProductID,
		Status:     InsertProduct.Status,
		Amount:     InsertProduct.Amount,
		Created_at: time.Now(),
	}
	repo.colProdTrans.InsertOne(context.Background(), &insertProd)
	return nil
}

func (repo *MongoDBRepository) GetProductTranscationByIDUser(id string) ([]user.ProductTransaction, error) {
	var Product []user.ProductTransaction
	UserID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.A{
		bson.M{
			"$match": bson.M{
				"userid": UserID,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "products",
				"localField":   "productid",
				"foreignField": "_id",
				"as":           "product",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userid",
				"foreignField": "_id",
				"as":           "account",
			},
		},
	}
	cur, err := repo.colProdTrans.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem user.ProductTransaction
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		Product = append(Product, elem)
	}
	return Product, nil
}
