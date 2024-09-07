-- name: CreateDummy :exec
INSERT INTO dummys(id, status, description)
VALUES (@id, @status, @description)
RETURNING (id, status, description, updated_at, created_at);
