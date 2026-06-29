UPDATE feature_flags
SET name = $1, description = $2, status = $3, environment = $4
WHERE id = $5;