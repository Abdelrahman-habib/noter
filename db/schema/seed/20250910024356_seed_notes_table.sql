-- +goose Up
INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 09:18:24 UTC'
);

-- Look up the seed user by email because its ID changes when reseeding
INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    '[seed] An old silent pond',
    E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + 365 * INTERVAL '24 hours',
    TRUE,
    (SELECT id FROM users WHERE email = 'alice@example.com')
);

INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440001',
    '[seed] Over the wintry forest',
    E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + 365 * INTERVAL '24 hours',
    TRUE,
    (SELECT id FROM users WHERE email = 'alice@example.com')
);

INSERT INTO notes (id, title, content, created, expires, public, created_by) VALUES (
    '550e8400-e29b-41d4-a716-446655440002',
    '[seed] First autumn morning',
    E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + 7 * INTERVAL '24 hours',
    FALSE,
    (SELECT id FROM users WHERE email = 'alice@example.com')
);
-- +goose Down

DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440000';
DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440001';
DELETE FROM notes WHERE id = '550e8400-e29b-41d4-a716-446655440002';
DELETE FROM users WHERE email = 'alice@example.com';
