CREATE TABLE team (
    team_name VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    team_name VARCHAR(255) NOT NULL REFERENCES team(team_name) ON DELETE CASCADE ON UPDATE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE status_enum AS ENUM ('OPEN', 'IN_REVIEW', 'APPROVED', 'REJECTED', 'MERGED', 'CLOSED');

CREATE TABLE pull_request (
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,

    status status_enum NOT NULL DEFAULT 'OPEN', 

    reviewer1_id VARCHAR(255) REFERENCES users(user_id) ON DELETE SET NULL,
    reviewer2_id VARCHAR(255) REFERENCES users(user_id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMPTZ,

    CHECK (author_id <> reviewer1_id),
    CHECK (author_id <> reviewer2_id)
);