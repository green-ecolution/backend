-- name: CreateAuthUser :one
INSERT INTO auth_users (name, username, email, email_verified, picture) VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: GetAuthUserById :one
SELECT * FROM auth_users WHERE id = $1;

-- name: GetAuthUserByAccountId :one
SELECT u.* FROM auth_users AS u INNER JOIN auth_accounts AS a ON u.id = a.user_id WHERE a.provider = $1 AND a.provider_id = $2;

-- name: UpdateAuthUser :exec
UPDATE auth_users SET name = $2, username = $3, email = $4, email_verified = $5, picture = $6 WHERE id = $1;


-- name: CreateAuthAccount :one
INSERT INTO auth_accounts (
  provider,
  provider_id,
  email,
  email_verified,
  access_token,
  expiry,
  refresh_token,
  refresh_expiry,
  token_type,
  id_token,
  name,
  preferred_username,
  nickname,
  picture,
  profile,
  user_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15,
  $16
) RETURNING id;

-- name: GetAuthAccountById :one
SELECT * FROM auth_accounts WHERE provider = $1 AND provider_id = $2;

-- name: GetAuthAccountsByUserId :many
SELECT * FROM auth_accounts WHERE user_id = $1;

-- name: GetAuthAccountByUserIdAndProvider :one
SELECT * FROM auth_accounts WHERE provider = $1 AND user_id = $2;

-- name: UpdateAuthAccount :exec
UPDATE auth_accounts SET
  email = $3,
  email_verified = $4,
  access_token = $5,
  expiry = $6,
  refresh_token = $7,
  refresh_expiry = $8,
  token_type = $9,
  id_token = $10,
  name = $11,
  preferred_username = $12,
  nickname = $13,
  picture = $14,
  profile = $15
WHERE
  provider = $1 AND
  provider_id = $2;

-- name: DeleteAuthAccount :exec
DELETE FROM auth_accounts WHERE provider = $1 AND provider_id = $2;

