package repository

import (
	"context"
	"log"
	"webservice/internal/domain"
	"webservice/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	collection = "report"
)

type Repository interface {
	GetReports(context.Context, int64) ([]domain.Report, error)
	InsertReport(context.Context, domain.Report) error
}

type reportsRepo struct {
	dbname string
	dbHost string
	dbPort string
}

func NewReportRepository(name, dbHost, dbPort string) Repository {
	return &reportsRepo{
		dbname: name,
		dbHost: dbHost,
		dbPort: dbPort,
	}
}

func (r *reportsRepo) GetReports(ctx context.Context, elapse int64) ([]domain.Report, error) {
	client, err := database.ConnectDB(r.dbHost, r.dbPort)

	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	var reports []domain.Report

	filter := bson.M{"report_time": bson.M{"$gte": elapse}}
	coll := client.Database(r.dbname).Collection(collection)
	res, err := coll.Find(ctx, filter)

	if err != nil {
		log.Println("read err : ", err)
		return nil, err
	}
	defer res.Close(ctx)

	for res.Next(ctx) {
		var r domain.Report
		if err = res.Decode(&r); err != nil {
			log.Println("decoding err ", err)
			continue
		}
		reports = append(reports, r)
	}
	return reports, err
}

func (r *reportsRepo) InsertReport(ctx context.Context, report domain.Report) (err error) {
	return database.Write(ctx, r.dbname, collection, report)
}
