package dbs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"xyzeshop/constval"
	"xyzeshop/payloads"

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
	InsertOne(interface{}, string) (interface{}, error)
	GetSingleRecordByEmail(string, string) *payloads.Verification
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

func ConDb() {
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = constval.MDBuri
	}
	opt := options.Client().ApplyURI(fmt.Sprintf("%s%s", "mongodb://", uri))
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
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
	//Mgr = &dbManager{conection: client, ctx: ctx, cancel: cancelFn}

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
			panic(err)
		}
	}()
}
