SELECT id, email, password_hash, name, surname, team_id
FROM users
WHERE email = $1