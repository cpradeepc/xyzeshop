package dbs

import (
	"context"
	"fmt"
	"log"
	"xyzeshop/constval"
	"xyzeshop/payloads"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	filter := bson.D{primitive.E{Key: "email", Value: mail}}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	err := coll.FindOne(context.TODO(), filter).Decode(resp)
	if err != nil {
		log.Println("error during finding :", err)
		return nil
	}
	return resp
}

func (mgr *dbManager) UpdateVerification(data payloads.Verification, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "email", Value: data.Email}}
	update := bson.D{primitive.E{Key: "$set", Value: data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}
func (mgr *dbManager) UpdateEmailVerifiedStatus(req payloads.Verification, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "email", Value: req.Email}}
	update := bson.D{primitive.E{Key: "$set", Value: req}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}

// get signle user
func (mgr *dbManager) GetSingleRecordByEmailForUser(mail string, collName string) payloads.RegdUser {
	var resp payloads.RegdUser
	filter := bson.D{primitive.E{Key: "email", Value: mail}}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	coll.FindOne(context.TODO(), filter).Decode(&resp)
	fmt.Println("regduser :", resp)
	return resp
}
func (mgr *dbManager) GetListProduct(page int, limit int, offset int, collName string) (products []payloads.Product, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	findOpt := options.Find()
	findOpt.SetSkip(int64(skip))
	findOpt.SetLimit(int64(limit))
	cur, err := coll.Find(context.TODO(), bson.M{}, findOpt)
	err = cur.All(context.TODO(), &products)
	itemCount, err := coll.CountDocuments(context.TODO(), bson.M{})
	return products, itemCount, err

}
func (mgr *dbManager) SearchProduct(page int, limit int, offset int, search string, collName string) (products []payloads.Product, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	findOpt := options.Find()
	findOpt.SetSkip(int64(skip))
	findOpt.SetLimit(int64(limit))

	searchFilter := bson.M{}
	if len(search) >= 3 {
		searchFilter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			{"description": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		}
	}

	cur, err := coll.Find(context.TODO(), searchFilter, findOpt)
	cur.All(context.TODO(), &products)
	count, err = coll.CountDocuments(context.TODO(), searchFilter)
	return products, count, err

}
func (mgr *dbManager) GetSingleProductById(id primitive.ObjectID, collName string) (product payloads.Product, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	err = coll.FindOne(context.TODO(), filter).Decode(&product)
	return product, err
}
func (mgr *dbManager) UpdateProduct(product payloads.Product, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: product.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: product}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}
func (mgr *dbManager) DeleteProduct(id primitive.ObjectID, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	return err
}
func (mgr *dbManager) GetSingleAddress(id primitive.ObjectID, collName string) (addr payloads.RegdUserAddress, err error) {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err = coll.FindOne(context.TODO(), filter).Decode(&addr)
	return addr, err
}
func (mgr *dbManager) GetSingleUserByUserId(id primitive.ObjectID, collName string) (user payloads.RegdUser) {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	coll.FindOne(context.TODO(), filter).Decode(&user)
	return user
}
func (mgr *dbManager) UpdateUser(user payloads.RegdUser, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.D{primitive.E{Key: "$set", Value: user}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}
func (mgr *dbManager) GetCartObjectById(id primitive.ObjectID, collName string) (cart payloads.Cart, err error) {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err = coll.FindOne(context.TODO(), filter).Decode(&cart)
	return cart, err
}
func (mgr *dbManager) UpdateCartToCheckout(cart payloads.Cart, collName string) error {
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	filter := bson.D{primitive.E{Key: "_id", Value: cart.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: cart}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}
