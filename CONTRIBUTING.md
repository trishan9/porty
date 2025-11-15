# Contributing to **porty**

Thank you for your interest in contributing!
Porty is an open-source project, and contributions of all kinds are welcome, from fixing bugs to adding features, improving documentation, or testing on different Linux distributions.

This document describes the contribution process, coding style, and development workflow.

## 🧭 Ways to Contribute

You can help improve porty by:

- Reporting bugs
- Suggesting features
- Improving the TUI/UX
- Adding packaging
- Optimizing scanning or reducing system calls
- Enhancing process detection
- Writing documentation
- Testing across distros (Ubuntu, Arch, Fedora, etc.)

## 🐛 Reporting Issues

Please open a GitHub issue if you encounter:

- Incorrect port detection
- Wrong user/PID mapping
- TUI rendering issues
- Incorrect colors or theme behavior
- Crash / panic
- Build failure
- Incorrect package installation (AUR, installer, etc.)

When submitting a bug report, include:

- OS / distribution
- Terminal emulator + shell
- Go version (if built from source)
- Porty version (`porty version`)
- Steps to reproduce
- Screenshots (if relevant)
- Logs / error messages

The more detail, the easier it is to fix.

## 🌱 Feature Requests

If you’d like to propose an idea:

1. Check if the feature already exists or has been requested.
2. Open a **Feature Request** issue.
3. Describe:

   - What problem it solves
   - Why it's useful
   - UI/UX expectations
   - Examples or mockups

All constructive feature requests are welcome.

## 🛠️ Development Setup

### Requirements

- Go 1.21+
- Linux system (supports `/proc`)
- Make (optional)
- A terminal that supports truecolor

### Clone the repository

```bash
git clone https://github.com/trishan9/porty
cd porty
```

### Build locally

```bash
go build -o porty .
```

Run:

```bash
./porty list
```

## 📦 Project Structure

```
cmd/       → CLI using Cobra
internal/        → core logic (scanning, parsing, killing)
tui/             → BubbleTea TUI code
install.sh       → universal Linux installer
```

## Coding Guidelines

### General

- Follow idiomatic Go style (gofmt, go vet, staticcheck).
- Avoid unnecessary dependencies.
- Keep functions small and focused.
- Prefer readability over cleverness.
- Comment complex logic (especially /proc parsing).

### TUI Code

- Maintain consistent styling.
- Avoid hardcoding full RGB sequences; use LipGloss.
- Keep UI responsive, minimize blocking operations.

### Error Handling

- Use informative errors (`fmt.Errorf("failed to X: %w", err)`).
- Avoid panics except in unreachable conditions.
- Ensure TUI does not crash on unusual system states.

## Testing

Before submitting a PR:

1. Test functionality manually:

   ```bash
   porty list
   porty kill --port 3000
   ```

2. Test on multiple distros if possible (Arch, Ubuntu, Fedora).
3. Verify that the TUI works in:

   - Kitty
   - Alacritty
   - GNOME Terminal
   - Konsole
   - WezTerm

4. Run `go vet`:

   ```bash
   go vet ./...
   ```

Automated tests (unit tests) can be added soon. PRs welcome.

## Pull Requests

### Before opening a PR

- Sync your branch with `main`.
- Ensure the code builds.
- Follow naming conventions (e.g., `feature/tcp-filter`, `fix/kernel-detection`).

### PR checklist

- [ ] Code compiles
- [ ] No breaking changes
- [ ] No unused imports or dead code
- [ ] Meaningful commit messages
- [ ] Tests added (if applicable)
- [ ] Documentation updated (if applicable)

### PR review

- Every PR gets reviewed for correctness, readability, and design.
- Maintainers may ask for small changes before merging.

## Release Process (for maintainers)

1. Bump version in Go binary:

   ```bash
   porty version
   ```

2. Push a new git tag:

   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

3. GitHub Actions builds & uploads:

   - Linux, macOS binaries
   - AUR updates (manual push)
   - `.deb` and `.rpm` packages (if configured)

## Final Notes

Contributions of all types are important: code, docs, packaging, UX, testing, ideas.

Thank you for contributing to **porty**.
