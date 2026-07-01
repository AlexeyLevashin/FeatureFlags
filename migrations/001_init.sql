CREATE TABLE teams(
    id SERIAL PRIMARY KEY ,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    team_id INT NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    CONSTRAINT fk_users_teams
        FOREIGN KEY(team_id)
        REFERENCES teams(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE feature_flags(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    environment VARCHAR(50) NOT NULL,
    owner_user_id INT NOT NULL,
    owner_team_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_flags_users
        FOREIGN KEY(owner_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT fk_flags_teams
        FOREIGN KEY(owner_team_id)
        REFERENCES teams(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE feature_flag_updates(
    id SERIAL PRIMARY KEY,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id INT,
    feature_flag_id INT NOT NULL,

    CONSTRAINT fk_users_update
       FOREIGN KEY(user_id)
       REFERENCES users(id)
       ON DELETE SET NULL
       ON UPDATE CASCADE,

   CONSTRAINT fk_flags_update
       FOREIGN KEY(feature_flag_id)
       REFERENCES feature_flags(id)
       ON DELETE CASCADE
       ON UPDATE CASCADE
);