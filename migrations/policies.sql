CREATE TABLE
    policies (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        account_id BIGINT UNSIGNED NOT NULL,
        team_member_id BIGINT UNSIGNED NOT NULL,
        resource VARCHAR(255) NOT NULL,
        action VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_policies_account_id_team_member_id_action ON policies (account_id, team_member_id, action);

INSERT INTO
    policies (account_id, team_member_id, resource, action)
VALUES
    (1, 2, '/blogs/*', 'read'),
    (1, 2, '/funnels/page/12', 'write'),