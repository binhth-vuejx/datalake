-- server/migrations/056_issue_mcp_config.up.sql
ALTER TABLE issue ADD COLUMN mcp_config jsonb DEFAULT NULL;
