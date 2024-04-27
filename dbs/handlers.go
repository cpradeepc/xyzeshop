package dbs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"xyzeshop/constval"
	"xyzeshop/payloads"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type dbManager struct {
	conection *mongo.Client
	ctx       context.Context
	cancel    context.CancelFunc
}

// signatures of regdUser,cart,address,varification,product
type Manager interface {
	InsertOne(data interface{}, collectionName string) (interface{}, error)
	GetSingleRecordByEmail(string, string) (payloads.Verification, error)
	UpdateVerification(payloads.Verification, string) error
	UpdateEmailVerifiedStatus(payloads.Verification, string) error
	GetSingleRecordByEmailForUser(string, string) payloads.RegdUser
	GetListProduct(int, int, int, string) ([]payloads.Product, int64, error)
	SearchProduct(int, int, int, string, string) ([]payloads.Product, int64, error)
	GetSingleProductById(primitive.ObjectID, string) (payloads.Product, error)
	UpdateProduct(payloads.Product, string) error
	DeleteProduct(primitive.ObjectID, string) error
	GetSingleAddress(primitive.ObjectID, string) (payloads.RegdUserAddress, error)
	GetSingleUserByUserId(primitive.ObjectID, string) payloads.RegdUser
	UpdateUser(payloads.RegdUser, string) error
	GetCartObjectById(primitive.ObjectID, string) (payloads.Cart, error)
	UpdateCartToCheckout(payloads.Cart, string) error
}

var Mgr Manager

func ConDb() (clnt *mongo.Client, contx context.Context, cancelFunc context.CancelFunc) {
	// uri := os.Getenv("DB_HOST")
	var uri string
	uriH := os.Getenv("DB_HOST")
	uriP := os.Getenv("DB_PORT")
	uri = fmt.Sprintf("%s%s%s", uriH, ":", uriP)

	if uriP == "" || uriH == "" {
		uri = fmt.Sprintf("%s%s%s", constval.MDBhost, ":", constval.MDBport)
		//uri = "host.docker.internal"
	}
	fmt.Println("db connection initializing")
	uriAdr := fmt.Sprintf("%s%s", "mongodb://", uri)
	opt := options.Client().ApplyURI(uriAdr)
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	log.Println("cancel :", cancelFn)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Println("Unable to initialize database connectors, Retrying...")
		ConDb()
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Unable to connect to the database, Retrying...")
		ConDb()
	}
	log.Printf("succeed connection to database at %s\n", uri)
	Mgr = &dbManager{conection: client, ctx: ctx, cancel: cancelFn}

	return client, ctx, cancelFn
}

func CloseCon(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	//cancelFunc to cancel to context
	defer cancel()

	//client provides a method to close
	//a mongoDB connection
	defer func() {
		//client.Disconnect method also has deadline.
		//return err if any
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("connection is being closed...")
			panic(err)
		}
	}()

}

func (mgr *dbManager) InsertOne(data interface{}, collectionName string) (interface{}, error) {
	log.Println("data: ", data)
	coll := mgr.conection.Database(constval.DbName).Collection(collectionName)

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (mgr *dbManager) GetSingleRecordByEmail(mail string, collName string) (payloads.Verification, error) {
	resp := payloads.Verification{}
	//filter := bson.M{"email": mail}
	// filter := bson.D{bson.E{"exists", "true"},
	// 	{"email", "cpradeepc@zohomail.in"},
	// }
	coll := mgr.conection.Database(constval.DbName).Collection(collName)
	fmt.Println("resp mail>> :", resp.Email)
	fmt.Println("resp >> :", resp)
	//var resp2 any
	fmt.Println(">> 1")
	err := coll.FindOne(context.TODO(), bson.M{"email": "cpradeepc@zohomail.in"}).Decode(&resp)
	// err := coll.FindOne(context.TODO(), filter).Decode(&resp2)
	//fmt.Println("resp2 :", resp)
	if err != nil {
		log.Println("error during finding :", err)
		return resp, err
	}
	return resp, nil
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
	cur, er := coll.Find(context.TODO(), bson.M{}, findOpt)
	er = cur.All(context.TODO(), &products)
	if er != nil {
		log.Println("error was invoked during find :", er)
	}
	log.Println("cursor: ", cur)
	itemCount, er := coll.CountDocuments(context.TODO(), bson.M{})
	return products, itemCount, er

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

	cur, er := coll.Find(context.TODO(), searchFilter, findOpt)
	er = cur.All(context.TODO(), &products)
	if er != nil {
		log.Println("error was invoked during find :", er)
	}
	log.Println("cursor: ", cur, er)
	count, er = coll.CountDocuments(context.TODO(), searchFilter)
	return products, count, er

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
