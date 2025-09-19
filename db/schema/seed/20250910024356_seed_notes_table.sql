-- +goose Up

-- Insert user first to get the ID for foreign key reference
INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 09:18:24'
);

-- Insert notes with UUID, public field, and created_by reference
INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    '[seed] An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY),
    TRUE,
    1
);

INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440001',
    '[seed] Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY),
    TRUE,
    1
);

INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440002',
    '[seed] First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY),
    FALSE,
    1
);
-- +goose Down

DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440000';
DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440001';
DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440002';
DELETE FROM users WHERE name = 'Alice Jones';

