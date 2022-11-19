package service

import (
	"context"
	"webservice/internal/domain"
	"webservice/internal/repository"
)

type radiantService struct {
	repository repository.Repository
}

func NewRadiantService(r repository.Repository) RadiantService {
	return &radiantService{
		repository: r,
	}
}

type RadiantService interface {
	GetReportsFrom(ctx context.Context, from string) (domain.ReportSample, error)
	GetReportsByDate(ctx context.Context, date string) ([]domain.Report, error)
}

func (sr *radiantService) GetReportsFrom(ctx context.Context, from string) (domain.ReportSample, error) {
	t := domain.ToStamp(from)
	reports, err := sr.repository.GetReports(ctx, t)

	if err != nil {
		return domain.ReportSample{}, err
	}

	return domain.FormatSample(reports), nil
}

func (sr *radiantService) GetReportsByDate(ctx context.Context, date string) ([]domain.Report, error) {
	return nil, nil
}
