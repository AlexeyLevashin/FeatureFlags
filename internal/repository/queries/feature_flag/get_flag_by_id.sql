SELECT id, name, description, status, environment, owner_user_id, owner_team_id, updated_at
FROM feature_flags
WHERE id = $1
