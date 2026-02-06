CREATE TABLE
    team_members (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(150) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO
    team_members (email)
VALUES
    ('alice@example.com'),
    ('bob@example.org'),
    ('carol@example.net');