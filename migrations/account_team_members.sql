CREATE TABLE
    account_team_members (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        account_id BIGINT UNSIGNED NOT NULL,
        team_member_id BIGINT UNSIGNED NOT NULL,
        UNIQUE KEY uniq_member (account_id, team_member_id),
        FOREIGN KEY (account_id) REFERENCES accounts (id),
        FOREIGN KEY (team_member_id) REFERENCES team_members (id)
    );

INSERT INTO
    account_team_members (account_id, team_member_id)
VALUES
    (1, 1);

INSERT INTO
    account_team_members (account_id, team_member_id)
VALUES
    (1, 2);

INSERT INTO
    account_team_members (account_id, team_member_id)
VALUES
    (2, 3);