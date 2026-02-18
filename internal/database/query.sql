-- name: GetShortURL :one
SELECT id, url FROM short_urls 
WHERE code = ? LIMIT 1;

-- name: GetShortURLById :one
SELECT id, code, url FROM short_urls 
WHERE id = ? LIMIT 1;

-- name: CreateShortURL :execresult
INSERT INTO short_urls (
    code, url
) VALUES (
    ?, ?
);

-- name: UpdateShortURL :exec
UPDATE short_urls
SET url = ?
WHERE id = ?;

-- name: DeleteShortURL :exec
DELETE FROM short_urls
WHERE id = ?;