package execenv

import "fmt"

// BuildCommentReplyInstructions returns the canonical block telling an agent
// how to post its reply for a comment-triggered task. Both the per-turn
// prompt (daemon.buildCommentPrompt) and the CLAUDE.md workflow
// (InjectRuntimeConfig) call this so the trigger comment ID and the
// --parent value cannot drift between surfaces.
//
// The explicit "do not reuse --parent from previous turns" wording exists
// because resumed Claude sessions keep prior turns' tool calls in context
// and will otherwise copy the old --parent UUID forward.
func BuildCommentReplyInstructions(issueID, triggerCommentID string) string {
	if triggerCommentID == "" {
		return ""
	}
	return fmt.Sprintf(
		"Reply by writing your full response to a temp file and piping it — this preserves newlines and formatting:\n\n"+
			"    cat > /tmp/reply.md << 'REPLY_EOF'\n"+
			"    <your full response here, with real newlines>\n"+
			"    REPLY_EOF\n"+
			"    multica issue comment add %s --parent %s --content-stdin < /tmp/reply.md\n\n"+
			"IMPORTANT: Always use --content-stdin (not --content) so that newlines in your response are preserved correctly.\n"+
			"Do NOT use --content \"...\" for multi-line responses — shell quoting will collapse newlines.\n"+
			"Always use the trigger comment ID shown above as --parent; do NOT reuse --parent values from previous turns.\n",
		issueID, triggerCommentID,
	)
}
