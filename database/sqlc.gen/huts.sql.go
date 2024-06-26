// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: huts.sql

package sqlc

import (
	"context"
)

const upsertHuts = `-- name: UpsertHuts :exec
INSERT INTO doc.hut (asset_id, status, name, region, x, y)
SELECT
    UNNEST($1::integer[]),
    UNNEST($2::text[]),
    UNNEST($3::text[]),
    UNNEST($4::text[]),
    UNNEST($5::float[]),
    UNNEST($6::float[])
ON CONFLICT (asset_id)
    DO UPDATE SET
        asset_id = EXCLUDED.asset_id,
        status = EXCLUDED.status,
        name = EXCLUDED.name,
        region = EXCLUDED.region,
        x = EXCLUDED.x,
        y = EXCLUDED.y
`

type UpsertHutsParams struct {
	AssetIds []int32
	Statuses []string
	Names    []string
	Regions  []string
	Xs       []float64
	Ys       []float64
}

func (q *Queries) UpsertHuts(ctx context.Context, arg UpsertHutsParams) error {
	_, err := q.db.Exec(ctx, upsertHuts,
		arg.AssetIds,
		arg.Statuses,
		arg.Names,
		arg.Regions,
		arg.Xs,
		arg.Ys,
	)
	return err
}
