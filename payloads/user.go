package payloads

import "go.mongodb.org/mongo-driver/bson/primitive"

//unregiter user,"users" table in db
type GuestUser struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
}

//Regd. user
type RegdUser struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	UserType  string             `json:"user_type" bson:"user_type"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
}

//Regd. user update
type RegdUserUpdate struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password" bson:"password"`
}

//Regd. user address,"addresses" table in db
type RegdUserAddress struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Address1 string             `json:"address_1" bson:"address_1"`
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	City     string             `json:"city" bson:"city"`
	Country  string             `json:"country" bson:"country"`
}

//Regd. user address add new, update
type RegdUserPostalAddress struct {
	Address1 string             `json:"address_1" bson:"address_1"`
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	City     string             `json:"city" bson:"city"`
	Country  string             `json:"country" bson:"country"`
}

//unregister user verification,"verification" table in db
type Verification struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"_id"`
	OTP       int64              `json:"otp" bson:"otp"`
	Status    bool               `json:"status" bson:"status"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}

//Regd. user login
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
