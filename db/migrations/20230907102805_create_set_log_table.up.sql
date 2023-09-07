-- Set Log table
CREATE TABLE IF NOT EXISTS set_log (
    id SERIAL PRIMARY KEY NOT NULL,
    set_id INT NOT NULL,
    opp_a_score INT NOT NULL,
    opp_b_score INT NOT NULL,
    scored_by_a BOOL NOT NULL,
    FOREIGN KEY (set_id) REFERENCES set(id)
);