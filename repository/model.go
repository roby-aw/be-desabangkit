package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Customer struct {
// 	ID            uint      `gorm:"primarykey"`
// 	Nama          string    `json:"nama"`
// 	Tanggal_lahir time.Time `json:"tanggal_lahir"`
// 	Gender        int       `json:"gender"`
// 	ID_Hobi       int       `json:"id_hobi"`
// 	ID_Jurusan    int       `json:"id_jurusan"`
// }

type Account struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email,omitempty" binding:"required"`
	Fullname   string             `bson:"fullname,omitempty" binding:"required"`
	Password   string             `bson:"password,omitempty" binding:"required"`
	Role_id    primitive.ObjectID `bson:"role_id,omitempty" binding:"required"`
	IsVerified bool               `bson:"isverified"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty" binding:"required"`
	Fullname string             `bson:"fullname,omitempty" binding:"required"`
	Password string             `bson:"password,omitempty" binding:"required"`
	Role_id  primitive.ObjectID `bson:"role_id,omitempty" binding:"required"`
	// Roles     Role               `bson:"roles"`
}
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Rolename    string             `bson:"rolename,omitempty" binding:"required" json:"rolename"`
	Rolelabel   string             `bson:"rolelabel,omitempty" binding:"required" json:"rolelabel"`
	Description string             `bson:"description,omitempty" binding:"required" json:"description"`
}

type CodeOtp struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
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
	Latitude       string             `json:"latitude" bson:"latitude,omitempty"`
	Longitude      string             `json:"longitude" bson:"longitude,omitempty"`
	UserID         primitive.ObjectID `json:"userid" bson:"userid,omitempty"`
	UserAddress    string             `json:"user_address" bson:"user_address,omitempty"`
	IsPreorder     bool               `json:"is_preorder" bson:"is_preorder"`
	IsApproved     bool               `json:"is_approved" bson:"is_approved"`
	Created_at     time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type InputProductTransaction struct {
	Status     string             `json:"status" bson:"status,omitempty"`
	Productid  primitive.ObjectID `json:"productid" bson:"productid,omitempty"`
	Userid     primitive.ObjectID `json:"userid" bson:"userid,omitempty"`
	Amount     int                `json:"amount" bson:"amount,omitempty"`
	Created_at time.Time          `json:"created_at" bson:"created_at,omitempty"`
}
