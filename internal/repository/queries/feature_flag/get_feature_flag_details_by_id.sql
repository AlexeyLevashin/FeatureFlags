SELECT
    f.id,
    f.name,
    f.description,
    f.environment,
    f.status,
    t.name AS owner_team,
    creator.name AS creator_name,
    creator.surname AS creator_surname,
    f.created_at,
    updater.name AS updater_name,
    updater.surname AS updater_surname,
    updates.updated_at

FROM feature_flags f

         LEFT JOIN teams t ON f.owner_team_id = t.id
         LEFT JOIN users creator ON f.owner_user_id = creator.id

         LEFT JOIN (
    SELECT DISTINCT ON (feature_flag_id) updated_at, user_id, feature_flag_id
    FROM feature_flag_updates
    ORDER BY feature_flag_id, updated_at DESC
) updates ON f.id = updates.feature_flag_id

         LEFT JOIN users updater ON updates.user_id = updater.id

WHERE f.id = $1;