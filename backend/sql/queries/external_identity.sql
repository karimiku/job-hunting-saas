-- name: InsertExternalIdentity :exec
INSERT INTO external_identities (id, user_id, provider, subject, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: FindExternalIdentityByProviderAndSubject :one
SELECT id, user_id, provider, subject, created_at
FROM external_identities
WHERE provider = $1 AND subject = $2;
