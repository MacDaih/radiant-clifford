package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	u "webservice/utils"
)

type Report struct {
	RptAt int64   `bson:"report_time" json:"time"`
	Temp  float64 `bson:"temp" json:"temperature"`
	Hum   float64 `bson:"hum" json:"humidity"`
	Light int32   `bson:"ligth" json:"light_lvl"`
}

type Overview struct {
	TempAverage float64 `json:"temp_av"`
	HumAverage  float64 `json:"hum_av"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
	MaxHum      float64 `json:"max_hum"`
	MinHum      float64 `json:"min_hum"`
}

type ReportSample struct {
	Metrics Overview
	Reports []Report
}

func connectDB() (*mongo.Client, error) {
	//URI & credentials to be specified
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

func GetReports(elapse int64) ([]Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var reports []Report
	cli, err := connectDB()
	defer cancel()
	if err != nil {
		return []Report{}, err
	}
	col := cli.Database("radiant_clifford").Collection("report")
	//Commented for dev purpose
	// filter := bson.M{"report_time": bson.M{"$gte": elapse}}
	// filter
	opt := options.Find()
	opt.SetSort(bson.M{"report_time": -1})
	c, err := col.Find(ctx, bson.D{}, opt)

	if u.ErrLog("Get Report Err : ", err) {
		return []Report{}, err
	}
	for c.Next(ctx) {
		var r Report
		c.Decode(&r)
		reports = append(reports, r)
	}
	return reports, nil
}
