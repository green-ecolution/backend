-- name: GetSessionById :one
SELECT * FROM token_sessions WHERE session_id = $1;

-- name: UpdateSession :exec
UPDATE token_sessions SET expiry = $2, refresh_jti = $3 WHERE session_id = $1;

-- name: CreateSession :one
INSERT INTO token_sessions (session_id, expiry, refresh_jti) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteSession :exec
DELETE FROM token_sessions WHERE session_id = $1;

-- name: DeleteSessionAfterExpiry :exec
DELETE FROM token_sessions WHERE expiry < NOW;

