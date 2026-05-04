package daemon

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/multica-ai/multica/server/internal/daemon/execenv"
)

// BuildPrompt constructs the task prompt for an agent CLI.
// Keep this minimal — detailed instructions live in CLAUDE.md / AGENTS.md
// injected by execenv.InjectRuntimeConfig.
func BuildPrompt(task Task) string {
	basePrompt := ""
	if task.ChatSessionID != "" {
		basePrompt = buildChatPrompt(task)
	} else if task.TriggerCommentID != "" {
		basePrompt = buildCommentPrompt(task)
	} else {
		basePrompt = buildDefaultPrompt(task)
	}

	// Inject skill enforcement block when agent has strict skills.
	skillBlock := ""
	if hasStrictSkills(task) {
		skillBlock = "\n\n" + buildSkillEnforcementInstructions(task)
	}

	// Inject MCP constraint block when issue has specific MCP config.
	mcpBlock := ""
	if len(task.IssueMcpConfig) > 0 && string(task.IssueMcpConfig) != "null" {
		mcpBlock = "\n\n" + buildMcpConstraintInstructions(task.IssueMcpConfig)
	}

	// Add mandatory output format for all tasks
	return basePrompt + skillBlock + mcpBlock + "\n\n" + buildMandatoryFormatInstructions()
}

// buildMandatoryFormatInstructions adds mandatory output format instructions
func buildMandatoryFormatInstructions() string {
	var b strings.Builder
	b.WriteString("=== MANDATORY OUTPUT FORMAT ===\n")
	b.WriteString("You MUST include the following format section at the END of your response.\n")
	b.WriteString("EACH field MUST be on its own separate line. Do NOT put multiple fields on the same line.\n")
	b.WriteString("CRITICAL: Use REAL newlines (line breaks) between each field. Do NOT use \\n as literal text.\n\n")
	b.WriteString("# BÁO CÁO CÔNG VIỆC\n")
	b.WriteString("**Thời gian:** [thời gian xử lý, ví dụ: 22/04/2026 14:30 - 22/04/2026 14:45]\n")
	b.WriteString("**Tóm tắt:** [tóm tắt ngắn gọn những gì đã làm - 1-2 câu]\n")
	b.WriteString("**Kết quả:**\n")
	b.WriteString("- [kết quả 1]\n")
	b.WriteString("- [kết quả 2]\n")
	b.WriteString("**Vấn đề:** [vấn đề gặp phải (nếu có), nếu không có thì viết 'Không có']\n")
	b.WriteString("**Skills sử dụng:** [liệt kê skills đã sử dụng, nếu không có thì để trống]\n")
	b.WriteString("**MCP sử dụng:** [liệt kê MCP đã sử dụng, nếu không có thì để trống]\n\n")
	b.WriteString("IMPORTANT INSTRUCTIONS:\n")
	b.WriteString("1. BEFORE STARTING WORK: State which skills you will use in this task.\n")
	b.WriteString("2. Understand the skill requirements and follow them precisely\n")
	b.WriteString("3. When listing skills in 'Skills sử dụng:', specify which skills you actually used and how\n")
	b.WriteString("4. When posting your response to the issue comment, include the FULL response including this format section. Do not truncate or omit any part.\n")
	b.WriteString("5. CHART RENDERING: When you need to display a chart or data visualization, use an echarts fenced code block with valid JSON:\n")
	b.WriteString("   ```echarts\n")
	b.WriteString("   { \"xAxis\": { \"type\": \"category\", \"data\": [\"Mon\",\"Tue\",\"Wed\"] }, \"yAxis\": { \"type\": \"value\" }, \"series\": [{ \"type\": \"bar\", \"data\": [120, 200, 150] }] }\n")
	b.WriteString("   ```\n")
	b.WriteString("   The system will automatically render this as an interactive chart. Do NOT output raw JSON without the echarts code block wrapper.\n\n")
	b.WriteString("Example (plain answer):\n")
	b.WriteString("5 + 6 = 11\n\n")
	b.WriteString("# BÁO CÁO CÔNG VIỆC \n")
	b.WriteString("**Thời gian:** 22/04/2026 14:30 - 22/04/2026 14:45\n")
	b.WriteString("**Tóm tắt:** Đã tính toán phép cộng 5 + 6 và post kết quả vào issue\n")
	b.WriteString("**Kết quả:**\n")
	b.WriteString("- Tính toán: 5 + 6 = 11\n")
	b.WriteString("- Post kết quả vào issue #123\n")
	b.WriteString("**Vấn đề:** Không có\n")
	b.WriteString("**Skills sử dụng:** simple-report (đã đọc skill content và follow format)\n")
	b.WriteString("**MCP sử dụng:** \n\n")
	return b.String()
}

func hasSkills(task Task) bool {
	if task.Agent == nil || len(task.Agent.Skills) == 0 {
		return false
	}
	return true
}

// hasStrictSkills returns true if the agent has at least one skill with Strict=true.
func hasStrictSkills(task Task) bool {
	if task.Agent == nil {
		return false
	}
	for _, s := range task.Agent.Skills {
		if s.Strict {
			return true
		}
	}
	return false
}

func buildDefaultPrompt(task Task) string {
	var b strings.Builder
	b.WriteString("You are running as a local coding agent for a Multica workspace.\n\n")
	fmt.Fprintf(&b, "Your assigned issue ID is: %s\n\n", task.IssueID)
	fmt.Fprintf(&b, "Start by running `multica issue get %s --output json` to understand your task, then complete it.\n", task.IssueID)
	return b.String()
}

// buildSkillEnforcementInstructions adds instructions to enforce skill usage
func buildSkillEnforcementInstructions(task Task) string {
	var b strings.Builder

	strictSkills := []SkillData{}
	for _, s := range task.Agent.Skills {
		if s.Strict {
			strictSkills = append(strictSkills, s)
		}
	}

	b.WriteString("=== STRICT SKILL ENFORCEMENT — MANDATORY ===\n\n")
	b.WriteString("This agent has STRICT MODE enabled. The following skills have been loaded directly from the database.\n")
	b.WriteString("You MUST follow these skills exactly. Do NOT search for skill files on disk — the content is already provided below.\n\n")

	b.WriteString("ASSIGNED STRICT SKILLS:\n")
	for _, skill := range strictSkills {
		b.WriteString(fmt.Sprintf("  - %s\n", skill.Name))
	}

	b.WriteString("\nABSOLUTE RULES:\n")
	b.WriteString("1. You MUST follow EVERY assigned skill's instructions before starting any work.\n")
	b.WriteString("2. You are ONLY permitted to perform actions described in the assigned skills.\n")
	b.WriteString("3. You MUST NOT perform any action not explicitly covered by the assigned skills.\n")
	b.WriteString("4. Do NOT use Bash, file system, or any tool to search for skill files — they are embedded below.\n")
	b.WriteString("5. If the task requires something not covered by the assigned skills, stop and report it.\n\n")

	// Embed full skill content directly — no file reading needed
	for _, skill := range strictSkills {
		fmt.Fprintf(&b, "=== SKILL CONTENT: %s ===\n\n", skill.Name)
		b.WriteString(skill.Content)
		b.WriteString("\n\n")
		if len(skill.Files) > 0 {
			b.WriteString("Supporting files:\n")
			for _, f := range skill.Files {
				fmt.Fprintf(&b, "--- %s ---\n%s\n\n", f.Path, f.Content)
			}
		}
		fmt.Fprintf(&b, "=== END SKILL: %s ===\n\n", skill.Name)
	}

	b.WriteString("=== END STRICT SKILL ENFORCEMENT ===\n")

	return b.String()
}

// extractFormatSection extracts the format section from skill content
func extractFormatSection(content string) string {
	lines := strings.Split(content, "\n")
	inFormatSection := false
	var formatLines []string
	foundFormat := false
	skipNextEmptyLine := false

	for _, line := range lines {
		lineLower := strings.ToLower(line)

		// Start of format section - look for various patterns
		if !inFormatSection {
			if strings.Contains(lineLower, "format chuẩn") ||
				strings.Contains(lineLower, "format section") ||
				strings.Contains(lineLower, "required format") ||
				strings.Contains(lineLower, "### format") ||
				(strings.Contains(lineLower, "format") && strings.Contains(lineLower, "chuẩn")) ||
				strings.HasPrefix(lineLower, "## format") ||
				strings.HasPrefix(lineLower, "### format") {
				inFormatSection = true
				skipNextEmptyLine = true
				continue
			}
		}

		if inFormatSection {
			// End of format section (next major section at same or higher level)
			trimmedLine := strings.TrimSpace(line)

			// Check for section headers that end the format section
			if strings.HasPrefix(line, "###") && !strings.Contains(lineLower, "format") {
				break
			}
			if strings.HasPrefix(line, "##") && !strings.Contains(lineLower, "format") {
				break
			}
			if strings.HasPrefix(line, "#") && !strings.Contains(lineLower, "format") && !strings.Contains(lineLower, "format") {
				break
			}
			if trimmedLine == "---" {
				break
			}

			// Skip empty lines right after the header
			if skipNextEmptyLine && trimmedLine == "" {
				continue
			}
			skipNextEmptyLine = false

			// Collect format lines
			// Include code block content but skip the markers themselves
			if trimmedLine != "```" {
				formatLines = append(formatLines, line)
				foundFormat = true
			}
		}
	}

	// If format section too long (like entire guide), truncate to first 30 lines
	if foundFormat && len(formatLines) > 30 {
		formatLines = formatLines[:30]
	}

	// Debug logging
	if foundFormat {
		fmt.Printf("DEBUG: Found format section with %d lines\n", len(formatLines))
	} else {
		fmt.Printf("DEBUG: No format section found in content (length: %d)\n", len(content))
		fmt.Printf("DEBUG: Content preview: %s\n", content[:min(200, len(content))])
	}

	if foundFormat && len(formatLines) > 0 {
		return strings.Join(formatLines, "\n")
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildCommentPrompt constructs a prompt for comment-triggered tasks.
// The triggering comment content is embedded directly so the agent cannot
// miss it, even when stale output files exist in a reused workdir.
// The reply instructions (including the current TriggerCommentID as --parent)
// are re-emitted on every turn so resumed sessions cannot carry forward a
// previous turn's --parent UUID.
func buildCommentPrompt(task Task) string {
	var b strings.Builder
	b.WriteString("You are running as a local coding agent for a Multica workspace.\n\n")
	fmt.Fprintf(&b, "Your assigned issue ID is: %s\n\n", task.IssueID)
	if task.TriggerCommentContent != "" {
		b.WriteString("[NEW COMMENT] A user just left a new comment that triggered this task. You MUST respond to THIS comment, not any previous ones:\n\n")
		fmt.Fprintf(&b, "> %s\n\n", task.TriggerCommentContent)
	}
	fmt.Fprintf(&b, "Start by running `multica issue get %s --output json` to understand your task, then complete it.\n\n", task.IssueID)
	b.WriteString(execenv.BuildCommentReplyInstructions(task.IssueID, task.TriggerCommentID))
	return b.String()
}

// buildChatPrompt constructs a prompt for interactive chat tasks.
func buildChatPrompt(task Task) string {
	var b strings.Builder
	b.WriteString("You are running as a chat assistant for a Multica workspace.\n")
	b.WriteString("A user is chatting with you directly. Respond to their message.\n\n")
	fmt.Fprintf(&b, "User message:\n%s\n", task.ChatMessage)
	return b.String()
}

// buildMcpConstraintInstructions injects a strict data-source constraint block
// when the issue has a specific MCP config selected. The agent MUST only use
// the listed MCP servers as data/tool sources and is forbidden from accessing
// any other external data sources.
func buildMcpConstraintInstructions(issueMcpConfig json.RawMessage) string {
	// Parse server names from the config
	var cfg struct {
		McpServers map[string]json.RawMessage `json:"mcpServers"`
	}
	serverNames := []string{}
	if err := json.Unmarshal(issueMcpConfig, &cfg); err == nil {
		for name := range cfg.McpServers {
			serverNames = append(serverNames, name)
		}
	}

	var b strings.Builder
	b.WriteString("=== MCP DATA SOURCE CONSTRAINT — MANDATORY ===\n\n")
	b.WriteString("This issue has been configured with specific MCP server(s) as the ONLY permitted data sources.\n\n")

	if len(serverNames) > 0 {
		b.WriteString("PERMITTED MCP SERVERS:\n")
		for _, name := range serverNames {
			fmt.Fprintf(&b, "  - %s\n", name)
		}
		b.WriteString("\n")
	}

	b.WriteString("ABSOLUTE RULES — YOU MUST FOLLOW THESE WITHOUT EXCEPTION:\n")
	b.WriteString("1. ALL data, information, and tool calls MUST come exclusively from the MCP server(s) listed above.\n")
	b.WriteString("2. You are STRICTLY FORBIDDEN from accessing any other data source, including:\n")
	b.WriteString("   - Reading local files outside the working directory\n")
	b.WriteString("   - Making HTTP/API calls to any endpoint NOT provided by the listed MCP servers\n")
	b.WriteString("   - Using your own training knowledge as a data source for facts, figures, or records\n")
	b.WriteString("   - Querying databases, APIs, or services not exposed through the listed MCP servers\n")
	b.WriteString("3. If the required data is NOT available through the listed MCP servers, you MUST:\n")
	b.WriteString("   - Stop immediately\n")
	b.WriteString("   - Report clearly that the data cannot be retrieved from the permitted MCP sources\n")
	b.WriteString("   - Do NOT attempt to substitute with data from any other source\n")
	b.WriteString("4. Every piece of data you use in your response MUST be traceable to a call made to one of the listed MCP servers.\n\n")
	b.WriteString("VIOLATION of these rules is not permitted under any circumstances.\n")
	b.WriteString("=== END MCP CONSTRAINT ===\n")

	return b.String()
}

// buildStrictSkillSystemPrompt builds a system prompt block containing the
// full content of all strict skills assigned to the agent.
// Returns empty string if no strict skills exist.
func buildStrictSkillSystemPrompt(skills []SkillData) string {
	var strictSkills []SkillData
	for _, s := range skills {
		if s.Strict {
			strictSkills = append(strictSkills, s)
		}
	}
	if len(strictSkills) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("=== STRICT SKILL MODE ===\n")
	b.WriteString("You have been assigned the following skills with STRICT MODE enabled.\n")
	b.WriteString("You MUST follow these skills exactly. Their content is provided below — do NOT search for skill files on disk.\n\n")

	for _, skill := range strictSkills {
		fmt.Fprintf(&b, "--- SKILL: %s ---\n", skill.Name)
		b.WriteString(skill.Content)
		b.WriteString("\n")
		for _, f := range skill.Files {
			fmt.Fprintf(&b, "[File: %s]\n%s\n", f.Path, f.Content)
		}
		fmt.Fprintf(&b, "--- END SKILL: %s ---\n\n", skill.Name)
	}

	b.WriteString("You MUST only perform actions described in the above skills.\n")
	b.WriteString("=== END STRICT SKILL MODE ===\n")
	return b.String()
}
