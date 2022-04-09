package repository

import (
	"context"
	"log"
	"os"
	"webservice/internal/database"
	"webservice/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	collection = "report"
)

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
)

type Repository interface {
	GetReports(context.Context, int64) ([]domain.Report, error)
	InsertReport(context.Context, domain.Report) error
}

type reportsRepo struct {
	dbname string
}

func NewReportRepository(name string) Repository {
	return &reportsRepo{
		dbname: name,
	}
}

func (r *reportsRepo) GetReports(ctx context.Context, elapse int64) ([]domain.Report, error) {
	client, err := database.ConnectDB(host, port)

	if err != nil {
		return nil, err
	}

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
