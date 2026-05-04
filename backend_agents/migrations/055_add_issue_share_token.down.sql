-- server/migrations/055_add_issue_share_token.down.sql
ALTER TABLE issue DROP COLUMN share_token;
