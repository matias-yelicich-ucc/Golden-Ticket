# Skill Registry

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When writing Go tests, using teatest, or adding test coverage. | go-testing | C:/Users/Usuario/.gemini/config/skills/go-testing/SKILL.md |
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | C:/Users/Usuario/.gemini/config/skills/branch-pr/SKILL.md |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | C:/Users/Usuario/.gemini/config/skills/issue-creation/SKILL.md |
| When user says "judgment day", "review adversarial", etc. | judgment-day | C:/Users/Usuario/.gemini/config/skills/judgment-day/SKILL.md |
| When creating a new AI agent skill. | skill-creator | C:/Users/Usuario/.gemini/config/skills/skill-creator/SKILL.md |

## Compact Rules

### go-testing
- Use Gomega matchers for assertions where applicable.
- For Bubbletea TUI, use the teatest package.
- Ensure all tests run cleanly with no race conditions using `go test -race`.

### branch-pr
- Ensure conventional commit messages are used.
- Make PR description clear and follow issue references.

### issue-creation
- Title issues clearly with prefix like bug:, feat:, refactor:.
- Describe reproduction steps for bugs.

### judgment-day
- Run dual blind reviews in parallel to get independent feedback.
- Do not bypass verification criteria.

### skill-creator
- Create skills in folders containing `SKILL.md` with proper YAML frontmatter.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| README.md | d:/programacion/Golden-Ticket/README.md | Project guidelines and documentation |
