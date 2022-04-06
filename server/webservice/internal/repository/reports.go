package repository

import (
	"context"
	"log"
	"os"
	"time"
	"webservice/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetReports(context.Context, int64) ([]domain.Report, error)
	InsertReport(context.Context, domain.Report) error
}

type reportsRepo struct {
	client *mongo.Client
}

func NewReportRepository(c *mongo.Client) Repository {
	return &reportsRepo{
		client: c,
	}
}

func (r *reportsRepo) GetReports(ctx context.Context, elapse int64) ([]domain.Report, error) {
	var reports []domain.Report

	var filter bson.M
	col := r.client.Database(os.Getenv("DB_NAME")).Collection("report")
	filter = bson.M{"report_time": bson.M{"$gte": elapse}}

	c, err := col.Find(ctx, filter, nil)

	if err != nil {
		log.Println(err)
		return []domain.Report{}, err
	}

	for c.Next(ctx) {
		var r domain.Report
		if err := c.Decode(&r); err != nil {
			log.Println(err)
			continue
		}
		reports = append(reports, r)
	}
	defer r.client.Disconnect(ctx)
	return reports, nil
}

func (r *reportsRepo) InsertReport(ctx context.Context, report domain.Report) error {
	col := r.client.Database("radiant").Collection("report")

	_, err := col.InsertOne(ctx, bson.M{
		"report_time": time.Now().Unix(),
		"temp":        report.Temp,
		"hum":         report.Hum,
		"light":       report.Light,
	})

	if err != nil {
		log.Println("db query error : ", err)
		return err
	}
	return r.client.Disconnect(ctx)
}
