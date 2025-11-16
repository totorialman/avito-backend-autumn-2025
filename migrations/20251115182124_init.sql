-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    user_id   TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    status            TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at         TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id         TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_reviewers_user_id ON reviewers(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_reviewers_user_id;
DROP INDEX IF EXISTS idx_pr_author_id;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_team_name;

DROP TABLE IF EXISTS reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
