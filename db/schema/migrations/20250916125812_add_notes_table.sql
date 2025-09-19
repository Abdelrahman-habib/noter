-- +goose Up
CREATE TABLE IF NOT EXISTS notes (
    id CHAR(36) NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    created_by INTEGER NOT NULL,
    public BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_notes_created ON notes(created);
ALTER TABLE notes ADD CONSTRAINT fk_notes_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE notes DROP CONSTRAINT fk_notes_created_by;
DROP INDEX idx_notes_created ON notes;
DROP TABLE IF EXISTS notes;
