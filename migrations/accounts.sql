CREATE TABLE
    accounts (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO
    accounts (name, created_at)
VALUES
    ('Acme Corp', NOW ()),
    ('Globex Inc', NOW ()),
    ('Umbrella Group', NOW ());
    ('Wayne Enterprises', NOW ());