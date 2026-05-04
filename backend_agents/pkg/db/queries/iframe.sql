-- name: ListIframes :many
SELECT * FROM iframe
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: GetIframe :one
SELECT * FROM iframe
WHERE id = $1;

-- name: GetIframeInWorkspace :one
SELECT * FROM iframe
WHERE id = $1 AND workspace_id = $2;

-- name: CreateIframe :one
INSERT INTO iframe (
    workspace_id, title, description, icon, iframe_url, iframe_script
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateIframe :one
UPDATE iframe SET
    title = COALESCE(sqlc.narg('title'), title),
    description = sqlc.narg('description'),
    icon = sqlc.narg('icon'),
    iframe_url = sqlc.narg('iframe_url'),
    iframe_script = sqlc.narg('iframe_script'),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteIframe :exec
DELETE FROM iframe WHERE id = $1;
