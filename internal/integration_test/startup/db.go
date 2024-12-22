package startup

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/client"
	erpinvoicefile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/erp_invoice_file"
	erporderfile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/erp_order_file"
	invoice "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/invoice"
	merchantconfig "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/merchant_config"
	orderfileuploadlog "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/order_file_upload_log"
	saleorder "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/sale_order"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao/mongodb/staff"
)

var mongoDB *mongo.Database
var mongoClient *mongo.Client

func InitMongoDBClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if mongoClient == nil {

		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context,
				startedEvent *event.CommandStartedEvent) {
				// fmt.Println(startedEvent.Command)
			},
		}
		opts := options.Client().
			ApplyURI("mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=rs0").
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mongoClient = client
	}
	return mongoClient
}

func InitMongoDB(mongoClient *mongo.Client) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if mongoDB == nil {
		err := mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(err)
		}
		mongoDB = mongoClient.Database("cst_wstyle_integration")
	}

	staff.InitCollections(mongoDB)
	merchantconfig.InitCollections(mongoDB)
	orderfileuploadlog.InitCollections(mongoDB)
	saleorder.InitCollections(mongoDB)
	erporderfile.InitCollections(mongoDB)
	erpinvoicefile.InitCollections(mongoDB)
	client.InitCollections(mongoDB)
	invoice.InitCollections(mongoDB)
	return mongoDB
}
