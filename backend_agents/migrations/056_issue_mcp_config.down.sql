-- server/migrations/056_issue_mcp_config.down.sql
ALTER TABLE issue DROP COLUMN IF EXISTS mcp_config;
