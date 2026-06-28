INSERT INTO feature_flags (
        name, description, status, environment, owner_user_id, owner_team_id
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id;