package jobs

import (
	"context"
	"time"

	"doc-app/config"
	"doc-app/database"
	"doc-app/database/sqlc.gen"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func workers() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &DocWorker[sqlc.DocCampsite]{})
	river.AddWorker(workers, &DocWorker[sqlc.DocHut]{})
	return workers
}

func jobs(config config.Config) []*river.PeriodicJob {
	return []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return DocArgs[sqlc.DocCampsite]{
					Config: config,
					API:    "https://api.doc.govt.nz/v2/campsites",
				}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return DocArgs[sqlc.DocHut]{
					Config: config,
					API:    "https://api.doc.govt.nz/v2/huts",
				}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	}
}

func Work(ctx context.Context, config config.Config, database database.Database) (*river.Client[pgx.Tx], error) {
	client, err := river.NewClient(riverpgxv5.New(database.Pool), &river.Config{
		Workers:      workers(),
		PeriodicJobs: jobs(config),
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
