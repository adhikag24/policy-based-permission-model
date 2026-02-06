CREATE TABLE
    accounts (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO
    accounts (name)
VALUES
    ('Acme Corp'),
    ('Globex Inc'),
    ('Umbrella Group'),
    ('Wayne Enterprises');