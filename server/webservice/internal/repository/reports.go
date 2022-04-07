package repository

import (
	"context"
	"log"
	"time"
	"webservice/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collection = "report"
)

type Repository interface {
	GetReports(context.Context, int64) ([]domain.Report, error)
	InsertReport(context.Context, domain.Report) error
}

type reportsRepo struct {
	db *mongo.Database
}

func NewReportRepository(c *mongo.Database) Repository {
	return &reportsRepo{
		db: c,
	}
}

func (r *reportsRepo) GetReports(ctx context.Context, elapse int64) ([]domain.Report, error) {
	var reports []domain.Report

	var filter bson.M
	col := r.db.Collection(collection)
	filter = bson.M{"report_time": bson.M{"$gte": elapse}}

	c, err := col.Find(ctx, filter, nil)

	if err != nil {
		log.Println(err)
		return []domain.Report{}, err
	}

	for c.Next(ctx) {
		var r domain.Report
		if err := c.Decode(&r); err != nil {
			log.Println("reading report error : ", err)
			continue
		}
		reports = append(reports, r)
	}
	defer r.db.Client().Disconnect(ctx)
	return reports, nil
}

func (r *reportsRepo) InsertReport(ctx context.Context, report domain.Report) error {
	col := r.db.Collection(collection)

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
	return r.db.Client().Disconnect(ctx)
}
