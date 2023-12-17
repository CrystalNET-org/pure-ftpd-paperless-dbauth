CREATE TABLE auth_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    is_superuser BOOLEAN NOT NULL,
    first_name VARCHAR(30),
    last_name VARCHAR(30),
    email VARCHAR(255),
    is_staff BOOLEAN NOT NULL,
    is_active BOOLEAN NOT NULL,
    date_joined TIMESTAMP NOT NULL
);
INSERT INTO auth_user (
    username,
    password,
    is_superuser,
    is_staff,
    is_active,
    date_joined
) VALUES (
    'test',
    'pbkdf2_sha256$600000$QPX4iFALW395lWmtuqJFpR$bEnT40kR4SrWm+ZsY5LB+G6JkMGa7YiCX2+pdwIdnGM=',
    'false',
    'false',
    'true',
    NOW()
)
-- password == testpassword