-- +goose Up
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);


-- +goose Down
DROP INDEX sessions_expiry_idx ON sessions;
DROP TABLE sessions;
