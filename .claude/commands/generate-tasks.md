Please generate tasks from the PRD using /.claude/rules/generate-tasks.md
If not explicitly told which PRD to use, generate a list of PRDs and ask the user to select one under `/.claude/tasks/` or create a new one using `create-prd.md`:
- assume it's stored under `/.claude/tasks/` and has a filename starting with `prd-` (e.g., `prd-[name].md`)
- it should not already have a corresponding task list in `/.claude/tasks/` (e.g., `tasks-prd-[name].md`)
- **always** ask the user to confirm the PRD file name before proceeding
Make sure to provide options in number lists so I can respond easily (if multiple options).