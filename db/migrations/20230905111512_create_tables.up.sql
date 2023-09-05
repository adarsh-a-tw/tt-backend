-- migrations/create_tables/up.sql

-- Match Table
CREATE TABLE IF NOT EXISTS match (
    id SERIAL PRIMARY KEY NOT NULL,
    stage TEXT NOT NULL,
    format TEXT NOT NULL,
    game_point INT NOT NULL,
    set_count INT NOT NULL,
    status TEXT NOT NULL
);

-- Set Table
CREATE TABLE IF NOT EXISTS set (
    id SERIAL PRIMARY KEY NOT NULL,
    set_number INT NOT NULL,
    match_id INT NOT NULL,
    opp_a_score INT NOT NULL,
    opp_b_score INT NOT NULL,
    is_completed BOOLEAN NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match(id)
);

-- Team Table
CREATE TABLE IF NOT EXISTS team (
    id SERIAL PRIMARY KEY NOT NULL,
    player_a TEXT NOT NULL,
    player_b TEXT NOT NULL
);

-- Player Table
CREATE TABLE IF NOT EXISTS player (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);

-- Team Match Mapping Table
CREATE TABLE IF NOT EXISTS team_match_mapping (
    match_id INT NOT NULL,
    team_id INT NOT NULL,
    is_opp_a BOOLEAN NOT NULL,
    is_winner BOOLEAN NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match(id),
    FOREIGN KEY (team_id) REFERENCES team(id)
);

-- Player Match Mapping Table
CREATE TABLE IF NOT EXISTS player_match_mapping (
    match_id INT NOT NULL,
    player_id INT NOT NULL,
    is_opp_a BOOLEAN NOT NULL,
    is_winner BOOLEAN NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match(id),
    FOREIGN KEY (player_id) REFERENCES player(id)
);
