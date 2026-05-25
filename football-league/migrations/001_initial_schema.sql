-- =============================================================================
-- 001_initial_schema.sql
-- Creates all tables, enforces constraints, seeds teams, fixtures, and standings.
-- Mounted into the postgres container at /docker-entrypoint-initdb.d/ so it runs
-- automatically on first startup. Safe to re-run (all statements are idempotent).
-- =============================================================================

-- BEGIN;

-- +goose Up

-- ── TEAMS ─────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS teams (
    id       SERIAL      PRIMARY KEY,
    name     VARCHAR(100) NOT NULL UNIQUE,
    attack   DECIMAL(4,2) NOT NULL DEFAULT 1.0 CHECK (attack  > 0),
    defense  DECIMAL(4,2) NOT NULL DEFAULT 1.0 CHECK (defense > 0)
);

-- ── MATCHES ───────────────────────────────────────────────────────────────────
-- home_score / away_score are NULL until the match is simulated.
-- The CHECK constraints prevent impossible data (e.g. week 7, negative scores).
CREATE TABLE IF NOT EXISTS matches (
    id           SERIAL  PRIMARY KEY,
    week         INTEGER NOT NULL CHECK (week >= 1 AND week <= 6),
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_score   INTEGER          CHECK (home_score >= 0),
    away_score   INTEGER          CHECK (away_score >= 0),
    played       BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT no_self_match  CHECK (home_team_id != away_team_id),
    CONSTRAINT unique_fixture UNIQUE (week, home_team_id, away_team_id)
);

-- ── STANDINGS ─────────────────────────────────────────────────────────────────
-- team_id is both PK and FK — one row per team, always.
-- updated_at is set automatically by standingRepo.Recalculate() via NOW().
CREATE TABLE IF NOT EXISTS standings (
    team_id       INTEGER     PRIMARY KEY REFERENCES teams(id),
    played        INTEGER     NOT NULL DEFAULT 0,
    won           INTEGER     NOT NULL DEFAULT 0,
    drawn         INTEGER     NOT NULL DEFAULT 0,
    lost          INTEGER     NOT NULL DEFAULT 0,
    goals_for     INTEGER     NOT NULL DEFAULT 0,
    goals_against INTEGER     NOT NULL DEFAULT 0,
    goal_diff     INTEGER     NOT NULL DEFAULT 0,
    points        INTEGER     NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── SEED: TEAMS ───────────────────────────────────────────────────────────────
-- Attack/Defense multipliers relative to league average (1.0).
-- Higher attack → more goals scored. Higher defense → fewer goals conceded.
-- These values feed directly into the Poisson lambda formula in the match engine.
INSERT INTO teams (name, attack, defense) VALUES
    ('Manchester City', 1.80, 1.30),   -- elite on both ends
    ('Liverpool',       1.60, 1.20),   -- strong attack and defence
    ('Arsenal',         1.40, 1.10),   -- above average
    ('Chelsea',         1.20, 1.00)    -- league average baseline
ON CONFLICT (name) DO NOTHING;

-- ── SEED: FIXTURES ────────────────────────────────────────────────────────────
-- Double round-robin: each team plays every other team once at home, once away.
-- 4 teams × 3 opponents × 2 legs = 12 matches total, 2 per week over 6 weeks.
--
-- Fixture grid:
--   Week 1: Man City (H) vs Liverpool,  Arsenal (H) vs Chelsea
--   Week 2: Man City (H) vs Arsenal,    Liverpool (H) vs Chelsea
--   Week 3: Man City (H) vs Chelsea,    Liverpool (H) vs Arsenal
--   Week 4: Liverpool (H) vs Man City,  Chelsea (H) vs Arsenal      ← reverse legs begin
--   Week 5: Arsenal (H)  vs Man City,   Chelsea (H) vs Liverpool
--   Week 6: Chelsea (H)  vs Man City,   Arsenal (H) vs Liverpool
--
-- IDs:  1=Manchester City  2=Liverpool  3=Arsenal  4=Chelsea
INSERT INTO matches (week, home_team_id, away_team_id) VALUES
    (1, 1, 2), (1, 3, 4),
    (2, 1, 3), (2, 2, 4),
    (3, 1, 4), (3, 2, 3),
    (4, 2, 1), (4, 4, 3),
    (5, 3, 1), (5, 4, 2),
    (6, 4, 1), (6, 3, 2)
ON CONFLICT DO NOTHING;

-- ── SEED: STANDINGS (all zeros) ───────────────────────────────────────────────
-- One row per team initialised to zero. standingRepo.Recalculate() will populate
-- actual values after each week is played.
INSERT INTO standings (team_id) VALUES (1), (2), (3), (4)
ON CONFLICT (team_id) DO NOTHING;

-- COMMIT;