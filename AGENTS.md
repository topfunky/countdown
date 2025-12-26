# Agent Protocols

Follow best practices that a principal level programmer would use when implementing code and infrastructure.

## Coding

- **Terraform**: Split `main.tf`, `variables.tf`, `outputs.tf`. Always run `terraform fmt`.
- **Idempotency**: Ensure initialization scripts verify state before acting. Fail fast on error.
- **Tests**: Write tests wherever possible and verify that tests pass before adding new features.
- **Languages**: Avoid JavaScript/Node.js unless absolutely necessary. Prefer compiled or high-performance tools (Rust, Go, etc.).
- **Scripts**: Implement scripts in separate files (bash, python, etc.) rather than embedding them as inline strings in YAML/TOML/JSON.
- **Dependencies**: Use current versions of dependencies, including GitHub Actions. Update dependencies to latest versions when adding or modifying code.

## Source Control (`jj`)

- Use `jj` exclusively for version control.
- **Initialization**: If the current repo is a git repo but not a jj project, initialize it with `jj git init` before any version control operations. This ensures no git commands will ever be run.
- **Confirmation**: Show changed files and stats (e.g. `jj diff --stat`) and ALWAYS ask the user for confirmation before creating a commit.
- **Commit**: `jj commit -m "type(scope): description"` (Atomic: describes current change and creates new one)
- **Sync**: NEVER sync with remote repository. Allow the user to sync manually outside of the agent.

## Commits

- **Style**: Conventional Commits.
- **Types**: `feat`, `fix`, `docs`, `chore`, `refactor`.
- **Format**: `type(scope): concise imperative description`.

## Documentation

- **Style**: Terse. Bullet points preferred.
- **Comments**: Explain **WHY** (design decisions, complexity), not **WHAT** (code behavior).

## Markdown

- Follow Markdown best practices.
- Ordered lists must all start with `1.`
