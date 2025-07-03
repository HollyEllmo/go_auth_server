INSERT INTO apps (id, name, secret)
VALUES (1, 'auth', 'test-secret')
ON CONFLICT DO NOTHING;