-- migrations/create_tables/down.sql

-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS player_match_mapping;
DROP TABLE IF EXISTS team_match_mapping;
DROP TABLE IF EXISTS player;
DROP TABLE IF EXISTS team;
DROP TABLE IF EXISTS set;
DROP TABLE IF EXISTS match;