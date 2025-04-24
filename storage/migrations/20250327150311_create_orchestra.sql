-- +goose Up
-- +goose StatementBegin
CREATE TABLE orchestra_info
(
--- Например Привелегии, правила членства, история оркестра, ссылка на вступлнение в чат/приватный тг-канал, чтобы все боты отсюда подтягивали, кароч весь текст сюда
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Таблица друзей
CREATE TABLE club_members
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name  TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE CONSTRAINT email_syntax CHECK (email ~ '^[^@]+@[^@]+\.[^@]+$'),
    phone      TEXT NOT NULL UNIQUE
        CONSTRAINT phone_number CHECK (phone ~ '^7\d{10}$'),
    status     TEXT NOT NULL    DEFAULT 'pending'
        CONSTRAINT status_variants CHECK (status IN ('pending', 'approved', 'declined')),
    created_at TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP
);

-- Таблица видов событий--
CREATE TABLE event_types
(
    id          SERIAL PRIMARY KEY,
    name        TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL
);

-- Таблица локаций--
CREATE TABLE locations
(
    id    SERIAL PRIMARY KEY,
    name  TEXT UNIQUE NOT NULL,
    route TEXT NOT NULL,
    features TEXT NOT NULL
);

-- Таблица событий
CREATE TABLE events
(
    id          SERIAL PRIMARY KEY,
    title       TEXT                                NOT NULL,
    description TEXT,
    event_type  INTEGER REFERENCES event_types(id),
    event_date  TIMESTAMPTZ                         NOT NULL
        CONSTRAINT event_not_in_future CHECK (event_date >= CURRENT_TIMESTAMP),
    location    INTEGER REFERENCES locations (id)   NOT NULL,
    capacity    INTEGER NOT NULL CHECK (capacity > 0),
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_timestamp_club
    BEFORE UPDATE ON club_members
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- Таблица регистраций на события
CREATE TABLE registrations
(
    id                  SERIAL PRIMARY KEY,
    user_id             UUID    NOT NULL REFERENCES club_members (id) ON DELETE CASCADE,
    event_id            INTEGER NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    registration_status TEXT    NOT NULL DEFAULT 'registered' CHECK (registration_status IN ('registered', 'cancelled')),
    created_at          TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, event_id)
);

CREATE TRIGGER trigger_update_timestamp_reg
    BEFORE UPDATE ON registrations
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orchestra_info;
DROP TABLE IF EXISTS registrations;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS club_members;
DROP TABLE IF EXISTS event_types;
-- +goose StatementEnd

