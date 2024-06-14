-- name: UpsertHuts :exec
INSERT INTO doc.hut (asset_id, status, name, region, x, y)
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

