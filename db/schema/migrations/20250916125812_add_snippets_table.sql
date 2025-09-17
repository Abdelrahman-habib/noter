-- +goose Up
CREATE TABLE IF NOT EXISTS snippets (
    id CHAR(36) NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    created_by INTEGER NOT NULL,
    public BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_snippets_created ON snippets(created);
ALTER TABLE snippets ADD CONSTRAINT fk_snippets_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE snippets DROP CONSTRAINT fk_snippets_created_by;
DROP INDEX idx_snippets_created ON snippets;
DROP TABLE IF EXISTS snippets;
