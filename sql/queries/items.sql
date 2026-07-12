-- name: AddItem :one
INSERT INTO items(
    name, created_at, updated_at
) VALUES(
@name, @time, @time
) RETURNING *;
