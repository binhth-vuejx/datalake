-- Add strict mode flag to agent_skill for skill enforcement
-- When strict=true, agent must explicitly load and confirm skill usage before task execution

ALTER TABLE agent_skill ADD COLUMN IF NOT EXISTS strict BOOLEAN NOT NULL DEFAULT false;
