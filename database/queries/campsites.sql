-- name: GetCampsite :one
SELECT
    *
FROM
    doc.campsite
WHERE
    id = @id;

-- name: GetCampsites :many
SELECT
    *
FROM
    doc.campsite
WHERE
    id = @id
ORDER BY
    id
LIMIT @lim;

-- name: UpsertCampsites :exec
INSERT INTO doc.campsite (asset_id, status, name, region, x, y)
SELECT
    UNNEST(@asset_ids::integer[]),
    UNNEST(@statuses::text[]),
    UNNEST(@names::text[]),
    UNNEST(@regions::text[]),
    UNNEST(@xs::float[]),
    UNNEST(@ys::float[])
ON CONFLICT (asset_id)
    DO UPDATE SET
        asset_id = EXCLUDED.asset_id,
        status = EXCLUDED.status,
        name = EXCLUDED.name,
        region = EXCLUDED.region,
        x = EXCLUDED.x,
        y = EXCLUDED.y;

