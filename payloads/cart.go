package payloads

import "go.mongodb.org/mongo-driver/bson/primitive"

//user Cart,"my_cart" table in db
type Cart struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Checkout  bool               `json:"checkout,omitempty" bson:"checkout"`
}

type CartUser struct {
	UserID    string `json:"user_id" bson:"user_id"`
	ProductID string `json:"product_id" bson:"product_id"`
}
