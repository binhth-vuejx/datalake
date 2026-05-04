-- server/migrations/055_add_issue_share_token.up.sql
ALTER TABLE issue ADD COLUMN share_token UUID UNIQUE;
