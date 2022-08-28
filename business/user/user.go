package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID            uint   `json:"id"`
	Nama          string `json:"nama"`
	Tanggal_lahir string `json:"tanggal_lahir"`
	Gender        string `json:"gender"`
	ID_Hobi       int    `json:"id_hobi"`
	ID_Gender     int    `json:"id_gender"`
}

type Regcustomer struct {
	Nama          string `json:"nama"`
	Tanggal_lahir string `json:"tanggal_lahir"`
	Gender        int    `json:"gender"`
}

type RegAccount struct {
	Fullname  string `json:"fullname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Role_id   string `json:"role_id,omitempty" binding:"required"`
	Url_photo string `json:"url_photo"`
}

type Account struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email,omitempty" binding:"required"`
	Fullname   string             `bson:"fullname,omitempty" binding:"required"`
	Password   string             `bson:"password,omitempty" binding:"required"`
	Role_id    primitive.ObjectID `bson:"role_id,omitempty" binding:"required"`
	IsVerified bool               `bson:"isverified,omitempty"`
	Roles      []Role             `bson:"roles" json:"roles"`
}

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name,omitempty" binding:"required" json:"rolename"`
	Label       string             `bson:"label,omitempty" binding:"required" json:"rolelabel"`
	Description string             `bson:"description,omitempty" binding:"required" json:"description"`
}

type AuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	ID    int
	Email string
	Role  string
	jwt.StandardClaims
}

type ResLogin struct {
	Account Account `json:"account"`
	Token   string  `json:"token"`
}

type Hobi struct {
	ID   int    `gorm:"primarykey" json:"id"`
	Nama string `json:"nama"`
}

type CodeOtp struct {
	Email      string    `json:"email" bson:"email,omitempty"`
	Code       string    `json:"code" bson:"code,omitempty"`
	Expired_at time.Time `json:"expired_at" bson:"expired_at,omitempty"`
}

type Product struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SKU            string             `json:"sku" bson:"sku,omitempty"`
	Photo_url      string             `json:"photo_url" bson:"photo_url,omitempty"`
	Name           string             `json:"name" bson:"name,omitempty"`
	PriceExpected  int                `json:"price_expected" bson:"price_expected,omitempty"`
	Quantity       int                `json:"Quantity" bson:"Quantity,omitempty"`
	Category       string             `json:"category" bson:"category,omitempty"`
	DeliveryOption string             `json:"delivery_option" bson:"delivery_option,omitempty"`
	Cooperationid  primitive.ObjectID `json:"cooperationid" bson:"cooperationid,omitempty"`
	Cooperations   []Cooperation      `json:"cooperations" bson:"cooperations,omitempty"`
	Latitude       string             `json:"latitude" bson:"latitude,omitempty"`
	Longitude      string             `json:"longitude" bson:"longitude,omitempty"`
	UserID         primitive.ObjectID `json:"userid" bson:"userid,omitempty"`
	Account        Account            `json:"account" bson:"account"`
	UserAddress    string             `json:"user_address" bson:"user_address,omitempty"`
	IsPreorder     bool               `json:"is_preorder" bson:"is_preorder,omitempty"`
	IsApproved     bool               `json:"is_approved" bson:"is_approved,omitempty"`
	Created_at     time.Time          `json:"created_at" bson:"created_at,omitempty"`
}
type ProductForTrans struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SKU            string             `json:"sku" bson:"sku,omitempty"`
	Photo_url      string             `json:"photo_url" bson:"photo_url,omitempty"`
	Name           string             `json:"name" bson:"name,omitempty"`
	PriceExpected  int                `json:"price_expected" bson:"price_expected,omitempty"`
	Quantity       int                `json:"Quantity" bson:"Quantity,omitempty"`
	Category       string             `json:"category" bson:"category,omitempty"`
	DeliveryOption string             `json:"delivery_option" bson:"delivery_option,omitempty"`
	Cooperationid  primitive.ObjectID `json:"cooperationid" bson:"cooperationid,omitempty"`
	Latitude       string             `json:"latitude" bson:"latitude,omitempty"`
	Longitude      string             `json:"longitude" bson:"longitude,omitempty"`
	UserID         primitive.ObjectID `json:"userid" bson:"userid,omitempty"`
	UserAddress    string             `json:"user_address" bson:"user_address,omitempty"`
	IsPreorder     bool               `json:"is_preorder" bson:"is_preorder,omitempty"`
	IsApproved     bool               `json:"is_approved" bson:"is_approved,omitempty"`
	Created_at     time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type InputProduct struct {
	SKU            string    `json:"sku" bson:"sku,omitempty"`
	Photo_url      string    `json:"photo_url" bson:"photo_url,omitempty"`
	Name           string    `json:"name" bson:"name,omitempty"`
	PriceExpected  int       `json:"price_expected" bson:"price_expected,omitempty"`
	Quantity       int       `json:"Quantity" bson:"Quantity,omitempty"`
	Category       string    `json:"category" bson:"category,omitempty"`
	DeliveryOption string    `json:"delivery_option" bson:"delivery_option,omitempty"`
	Cooperationid  string    `json:"cooperationid" bson:"cooperationid,omitempty"`
	Latitude       string    `json:"latitude" bson:"latitude,omitempty"`
	Longitude      string    `json:"longitude" bson:"longitude,omitempty"`
	UserID         string    `json:"userid" bson:"userid,omitempty"`
	UserAddress    string    `json:"user_address" bson:"user_address,omitempty"`
	Created_at     time.Time `json:"created_at" bson:"created_at,omitempty"`
}

type Cooperation struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name,omitempty"`
	Address    string             `json:"address" bson:"address,omitempty"`
	City       string             `json:"city" bson:"city,omitempty"`
	Province   string             `json:"province" bson:"province,omitempty"`
	Latitude   string             `json:"latitude" bson:"latitude,omitempty"`
	Longitude  string             `json:"longitude" bson:"longitude,omitempty"`
	Photo_url  string             `json:"photo_url" bson:"photo_url,omitempty"`
	Email      string             `json:"email" bson:"email,omitempty"`
	Username   string             `json:"username" bson:"username,omitempty"`
	Password   string             `json:"password" bson:"password,omitempty"`
	Created_at time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type ProductTransaction struct {
	Status     string            `json:"status" bson:"status,omitempty"`
	Productid  string            `json:"productid" bson:"productid,omitempty"`
	Product    []ProductForTrans `json:"product" bson:"product,omitempty"`
	Userid     string            `json:"userid" bson:"userid,omitempty"`
	Account    []Account         `json:"account" bson:"account,omitempty"`
	Amount     int               `json:"amount" bson:"amount,omitempty"`
	Created_at time.Time         `json:"created_at" bson:"created_at,omitempty"`
}

type InputProductTransaction struct {
	Status     string    `json:"status" bson:"status,omitempty"`
	Productid  string    `json:"productid" bson:"productid,omitempty"`
	Userid     string    `json:"userid" bson:"userid,omitempty"`
	Amount     int       `json:"amount" bson:"amount,omitempty"`
	Created_at time.Time `json:"created_at" bson:"created_at,omitempty"`
}
