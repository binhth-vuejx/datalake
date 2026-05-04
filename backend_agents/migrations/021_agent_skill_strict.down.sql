-- Remove strict mode flag from agent_skill

ALTER TABLE agent_skill DROP COLUMN IF EXISTS strict;
