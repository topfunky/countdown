# Agent Protocols

Follow best practices that a principal level programmer would use when implementing code and infrastructure.

## Coding

- **Terraform**: Split `main.tf`, `variables.tf`, `outputs.tf`. Always run `terraform fmt`.
- **Idempotency**: Ensure initialization scripts verify state before acting. Fail fast on error.
- **Tests**: Write tests wherever possible and verify that tests pass before adding new features.
- **Languages**: Avoid JavaScript/Node.js unless absolutely necessary. Prefer compiled or high-performance tools (Rust, Go, etc.).
- **Scripts**: Implement scripts in separate files (bash, python, etc.) rather than embedding them as inline strings in YAML/TOML/JSON.
- **Dependencies**: Use current versions of dependencies, including GitHub Actions. Update dependencies to latest versions when adding or modifying code.

## Version Control

See [`.ai/agents/VERSION_CONTROL.md`](.ai/agents/VERSION_CONTROL.md) for comprehensive version control guidelines using Jujutsu (jj).

## Documentation

- **Style**: Terse. Bullet points preferred.
- **Comments**: Explain **WHY** (design decisions, complexity), not **WHAT** (code behavior).

## Markdown

- Follow Markdown best practices.
- Ordered lists must all start with `1.`
