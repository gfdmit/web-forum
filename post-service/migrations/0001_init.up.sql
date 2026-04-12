CREATE SCHEMA IF NOT EXISTS forum;

CREATE TABLE IF NOT EXISTS forum.users (
    id              BIGSERIAL PRIMARY KEY,
    login           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS forum.profiles (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL UNIQUE REFERENCES forum.users(id) ON DELETE CASCADE,
    university_id   TEXT NOT NULL UNIQUE,
    firstname       TEXT NOT NULL,
    lastname        TEXT NOT NULL,
    middlename      TEXT,
    birthday        DATE NOT NULL,
    faculty         TEXT NOT NULL,
    grade           TEXT NOT NULL,
    "group"         TEXT NOT NULL,
    status          TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON forum.profiles (user_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS forum.boards (
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT UNIQUE NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS forum.posts (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES forum.users(id) ON DELETE SET NULL,
    board_id        BIGINT NOT NULL REFERENCES forum.boards(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    text            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON forum.posts (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_posts_board_id ON forum.posts (board_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS forum.comments (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES forum.users(id) ON DELETE SET NULL,
    post_id         BIGINT NOT NULL REFERENCES forum.posts(id) ON DELETE CASCADE,
    text            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON forum.comments (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON forum.comments (post_id) WHERE deleted_at IS NULL;
