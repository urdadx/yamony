-- name: CreatePreferences :one
INSERT INTO preferences (
    page_id,
    user_id,
    social_icons_position,
    hide_shop,
    hide_link_branding,
    hide_share_button,
    shop_layout,
    button_style,
    theme
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPreferencesByID :one
SELECT * FROM preferences
WHERE id = $1;

-- name: GetPreferencesByPageID :one
SELECT * FROM preferences
WHERE page_id = $1;

-- name: GetPreferencesByUserID :many
SELECT * FROM preferences
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdatePreferences :one
UPDATE preferences
SET
    social_icons_position = COALESCE(sqlc.narg('social_icons_position'), social_icons_position),
    hide_shop = COALESCE(sqlc.narg('hide_shop'), hide_shop),
    hide_link_branding = COALESCE(sqlc.narg('hide_link_branding'), hide_link_branding),
    hide_share_button = COALESCE(sqlc.narg('hide_share_button'), hide_share_button),
    shop_layout = COALESCE(sqlc.narg('shop_layout'), shop_layout),
    button_style = COALESCE(sqlc.narg('button_style'), button_style),
    theme = COALESCE(sqlc.narg('theme'), theme),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdatePreferencesByPageID :one
UPDATE preferences
SET
    social_icons_position = COALESCE(sqlc.narg('social_icons_position'), social_icons_position),
    hide_shop = COALESCE(sqlc.narg('hide_shop'), hide_shop),
    hide_link_branding = COALESCE(sqlc.narg('hide_link_branding'), hide_link_branding),
    hide_share_button = COALESCE(sqlc.narg('hide_share_button'), hide_share_button),
    shop_layout = COALESCE(sqlc.narg('shop_layout'), shop_layout),
    button_style = COALESCE(sqlc.narg('button_style'), button_style),
    theme = COALESCE(sqlc.narg('theme'), theme),
    updated_at = NOW()
WHERE page_id = sqlc.arg('page_id')
RETURNING *;

-- name: DeletePreferences :exec
DELETE FROM preferences
WHERE id = $1;

-- name: DeletePreferencesByPageID :exec
DELETE FROM preferences
WHERE page_id = $1;

-- name: UpsertPreferences :one
INSERT INTO preferences (
    page_id,
    user_id,
    social_icons_position,
    hide_shop,
    hide_link_branding,
    hide_share_button,
    shop_layout,
    button_style,
    theme
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (page_id) DO UPDATE SET
    social_icons_position = EXCLUDED.social_icons_position,
    hide_shop = EXCLUDED.hide_shop,
    hide_link_branding = EXCLUDED.hide_link_branding,
    hide_share_button = EXCLUDED.hide_share_button,
    shop_layout = EXCLUDED.shop_layout,
    button_style = EXCLUDED.button_style,
    theme = EXCLUDED.theme,
    updated_at = NOW()
RETURNING *;
