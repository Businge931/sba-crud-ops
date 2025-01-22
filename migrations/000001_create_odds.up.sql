-- Create odds table
CREATE TABLE IF NOT EXISTS odds (
    id BIGSERIAL PRIMARY KEY,
    league VARCHAR(255) NOT NULL,
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    home_team_win_odds DECIMAL(10,2) NOT NULL,
    away_team_win_odds DECIMAL(10,2) NOT NULL,
    draw_odds DECIMAL(10,2) NOT NULL,
    game_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(league, home_team, away_team, game_date)
);
