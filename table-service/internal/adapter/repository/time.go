package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
	"time"
)

func (i *TableRepoImplement) GetCurrentTime(ctx context.Context) (result domain.TestTime, customError *exceptions.CustomError) {
	ts, err := i.repository.GetTimeNow(ctx)
	if err != nil {
		return domain.TestTime{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch time: %w", err),
		}
	}

	if !ts.Valid {
		return domain.TestTime{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch time: %w", err),
		}
	}

	t := ts.Time

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return domain.TestTime{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch time: %w", err),
		}
	}

	now := time.Now()
	iso8601 := now.Format(time.RFC3339)
	timeFromDB := t.In(loc).Format(time.RFC3339)

	result.TimeFromDB = timeFromDB
	result.TimeFromAPI = iso8601

	return result, nil
}
