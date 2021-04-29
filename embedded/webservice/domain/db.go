package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Report struct {
	RptAt int64   `bson:"report_time" json:"time"`
	Temp  float64 `bson:"temp" json:"t"`
	Hum   float64 `bson:"hum" json:"h"`
	Light int32   `bson:"ligth" json:"l"`
}

type Reports []Report

func connectDB() (*mongo.Client, error) {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func GetReports() (Reports, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var reports Reports
	cli, err := connectDB()
	if err != nil {
		return Reports{}, err
	}
	col := cli.Database("radiant_clifford").Collection("report")
	c, err := col.Find(ctx, bson.D{})
	if err != nil {
		return Reports{}, err
	}
	for c.Next(ctx) {
		var r Report
		c.Decode(&r)
		reports = append(reports, r)
	}
	return reports, nil
}
