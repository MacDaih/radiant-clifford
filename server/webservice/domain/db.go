package domain

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	u "webservice/utils"
)

var authOpt = options.Credential{
	AuthMechanism: os.Getenv("AUTH"),
	AuthSource:    os.Getenv("DB_NAME"),
	Username:      os.Getenv("DB_USER"),
	Password:      os.Getenv("DB_PWD"),
}

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
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientopt := options.Client().SetAuth(authOpt).ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientopt)
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
	var filter bson.M
	col := cli.Database(os.Getenv("DB_NAME")).Collection("report")
	filter = bson.M{"report_time": bson.M{"$gte": elapse}}

	c, err := col.Find(ctx, filter, nil)

	if u.ErrLog("Get Report Err : ", err) {
		return []Report{}, err
	}

	for c.Next(ctx) {
		var r Report
		u.ErrLog("error fetching report : ", c.Decode(&r))
		reports = append(reports, r)
	}
	defer cli.Disconnect(ctx)
	return reports, nil
}

func InsertReport(r Report) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cli, err := connectDB()

	if err != nil {
		log.Println("db conn error : ", err)
		return err
	}
	col := cli.Database("radiant").Collection("report")

	_, err = col.InsertOne(ctx, bson.M{
		"report_time": time.Now().Unix(),
		"temp":        r.Temp,
		"hum":         r.Hum,
		"light":       r.Light,
	})

	if err != nil {
		log.Println("db query error : ", err)
		return err
	}

	return cli.Disconnect(ctx)
}
