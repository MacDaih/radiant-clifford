package actions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	d "../domain"
	u "../utils"
)

func connectDB() *mongo.Client {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		u.ErrLog("DB Err : ", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		u.ErrLog("DB Err : ", err)
	}
	return client
}

func InsertReports(arr []d.Report) {
	var res []interface{}
	for _, j := range arr {
		res = append(res, j)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cli := connectDB()
	col := cli.Database("radiant_clifford").Collection("report")
	_, err := col.InsertMany(ctx, res)
	if err != nil {
		u.ErrLog("Data Err : ", err)
	}
	defer cli.Disconnect(ctx)
}
