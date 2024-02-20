package dbs

import (
	"context"
	"log"
	"xyzeshop/constval"
	"xyzeshop/payloads"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (mgr *dbManager) InsertOne(data interface{}, collectionName string) (interface{}, error) {
	log.Println("data: ", data)
	coll := mgr.conection.Database(constval.DbName).Collection(collectionName)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (mgr *dbManager) GetSingleRecordByEmail(mail string, collName string) *payloads.Verification {
	resp := &payloads.Verification{}
	filter := bson.D{{"email", mail}}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	err := coll.FindOne(context.TODO(), filter).Decode(resp)
	if err != nil {
		log.Println("error during finding :", err)
		return nil
	}
	return resp
}

func (mgr *dbManager) UpdateVerification(data payloads.Verification) error                     {}
func (mgr *dbManager) UpdateEmailVerifiedStatus(payloads.Verification, string) error           {}
func (mgr *dbManager) GetSingleRecordByEmailForUser(string, string) payloads.RegdUser          {}
func (mgr *dbManager) GetListProduct(int, int, int, string) ([]payloads.Product, int64, error) {}
func (mgr *dbManager) SearchProduct(int, int, int, string, string) ([]payloads.Product, int64, error) {
}
func (mgr *dbManager) GetSignleProductById(primitive.ObjectID, string) (payloads.Product, error) {}
func (mgr *dbManager) UpdateProduct(payloads.Product, string) error                              {}
func (mgr *dbManager) DeleteProduct(primitive.ObjectID, string) error                            {}
func (mgr *dbManager) GetSingleAddress(primitive.ObjectID, string) (payloads.RegdUserAddress, error) {
}
func (mgr *dbManager) GetSingleUserByUserId(primitive.ObjectID, string) payloads.RegdUser  {}
func (mgr *dbManager) UpdateUser(payloads.RegdUser, string) error                          {}
func (mgr *dbManager) GetCartObjectById(primitive.ObjectID, string) (payloads.Cart, error) {}
func (mgr *dbManager) UpdateCartToCheckout(payloads.Cart, string)                          {}
