package jobs

import (
	"context"
	"doc-app/config"
	"doc-app/database"
	"doc-app/database/sqlc.gen"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/riverqueue/river"
)

type DocArgs[T any] struct {
	API string `json:"api"`
}

var _ river.JobArgs = DocArgs[struct{}]{}

func (d DocArgs[T]) Kind() string {
	return fmt.Sprintf("doc-%T-fetcher", *new(T))
}

type DocWorker[T any] struct {
	river.WorkerDefaults[DocArgs[T]]
	config   config.Config
	database database.Database
}

var _ river.Worker[DocArgs[struct{}]] = &DocWorker[struct{}]{}

func (d *DocWorker[T]) Work(ctx context.Context, job *river.Job[DocArgs[T]]) error {
	logger := slog.With("job", job.Kind)

	logger.With("api", job.Args.API).InfoContext(ctx, "fetching data")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.Args.API, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", d.config.API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data []T
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	queries := sqlc.New(d.database)
	switch d := any(data).(type) {
	case []sqlc.DocCampsite:
		params := sqlc.UpsertCampsitesParams{}
		for _, campsite := range d {
			params.AssetIds = append(params.AssetIds, campsite.AssetID)
			params.Statuses = append(params.Statuses, campsite.Status)
			params.Names = append(params.Names, campsite.Name)
			params.Regions = append(params.Regions, campsite.Region)
			params.Xs = append(params.Xs, campsite.X)
			params.Ys = append(params.Ys, campsite.Y)
		}

		if err := queries.UpsertCampsites(ctx, params); err != nil {
			return err
		}
	case []sqlc.DocHut:
		params := sqlc.UpsertHutsParams{}
		for _, hut := range d {
			params.AssetIds = append(params.AssetIds, hut.AssetID)
			params.Statuses = append(params.Statuses, hut.Status)
			params.Names = append(params.Names, hut.Name)
			params.Regions = append(params.Regions, hut.Region)
			params.Xs = append(params.Xs, hut.X)
			params.Ys = append(params.Ys, hut.Y)
		}

		if err := queries.UpsertHuts(ctx, params); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported data type: %T", d)
	}

	logger.With("count", len(data)).InfoContext(ctx, "data fetched and upserted")

	return nil
}
