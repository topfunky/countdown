# Countdown

A terminal countdown timer with animated spinners.

![Screen recording of liftoff](https://vhs.charm.sh/vhs-22kW7Is8pCwf5CXnk3oMLE.gif)

```
🌙 Liftoff in 99
```

Includes an option to render big numbers.

![Screen recording of big numbers](https://vhs.charm.sh/vhs-3UPyKQ0dKOOUpwHoBVW678.gif)

## Installation (in development)

### Homebrew (macOS/Linux)

```sh
brew install topfunky/tap/countdown
```

### Go install

```sh
go install github.com/topfunky/countdown@latest
```

### Binary releases

Download pre-built binaries from the [Releases](https://github.com/topfunky/countdown/releases) page for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### Build from source

```sh
git clone https://github.com/topfunky/countdown.git
cd countdown
make install-deps
make
```

## Usage

```sh
countdown [flags]
```

### Examples

```sh
# Default countdown from 100 to 0
countdown

# Quick 10-second countdown
countdown -r 10..0

# Countdown with moon spinner
countdown --spinner moon

# Custom title and range
countdown --title "Launch in" -r 60..0

# Slow countdown (every five seconds)
countdown -r 30..0 -t 5

# Count up instead of down
countdown -r 0..100

# Decrement by 5 each step
countdown -r 100..0 -d 5

# Custom colors
countdown --spinner.foreground 201 --title.foreground 39

# With padding
countdown --padding "1 2"

# Big ASCII art numbers
countdown -b -r 10..0
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-h, --help` | | Show help |
| `-v, --version` | | Print version |
| `-s, --spinner` | `dot` | Spinner animation type |
| `--title` | `Liftoff in` | Text displayed before the number |
| `-r, --range` | `100..0` | Start and end numbers (e.g., `10..0` or `0..100`) |
| `-t, --time-interval` | `1` | Seconds between each tick |
| `-d, --decrement` | `1` | Amount to change count each tick |
| `-f, --final-phase` | `5` | Threshold for final phase styling (number or percentage like `10%`) |
| `-b, --big` | `false` | Display numbers using large ASCII art digits |

### Style Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--spinner.foreground` | `212` | Spinner color (ANSI 0-255 or hex) |
| `--spinner.background` | | Spinner background color |
| `--title.foreground` | | Title text color |
| `--title.background` | | Title background color |
| `--padding` | `0 0` | Vertical and horizontal padding |

### Spinner Types

| Type | Description |
|------|-------------|
| `dot` | Braille dot pattern (default) |
| `line` | Rotating line `\|/-\` |
| `minidot` | Small dots |
| `jump` | Bouncing dot |
| `pulse` | Pulsing dot |
| `points` | Expanding points |
| `globe` | Rotating globe 🌍 |
| `moon` | Moon phases 🌙 |
| `monkey` | See no evil monkey 🙈 |
| `meter` | Progress meter |
| `hamburger` | Hamburger menu animation |
| `bomb` | Bomb and explosion 💣💥 |
| `none` | No spinner |

### Environment Variables

Some flags can be set via environment variables:

| Variable | Flag |
|----------|------|
| `COUNTDOWN_SPINNER` | `--spinner` |
| `COUNTDOWN_TITLE` | `--title` |
| `COUNTDOWN_SPINNER_FOREGROUND` | `--spinner.foreground` |
| `COUNTDOWN_SPINNER_BACKGROUND` | `--spinner.background` |
| `COUNTDOWN_TITLE_FOREGROUND` | `--title.foreground` |
| `COUNTDOWN_TITLE_BACKGROUND` | `--title.background` |
| `COUNTDOWN_PADDING` | `--padding` |

### Final Phase

When the countdown reaches the final phase threshold, colors are inverted to create visual emphasis. Set with `-f` or `--final-phase`:

- Absolute number: `-f 5` (triggers at 5, the default)
- Percentage: `-f 10%` (triggers at 10% of total range)

### Controls

- `q`, `Esc`, or `Ctrl+C` to quit early

## Development

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) from Charm.

### Make targets

```sh
make build        # Build binary
make test         # Run tests
make format       # Format code
make lint         # Run linter
make install-deps # Install dev dependencies
make snapshot     # Build snapshot release locally
make clean        # Remove build artifacts
```

### Releasing

Releases are automated via GitHub Actions. To create a new release, first commit with `jj` and `push` the `main` branch to GitHub.

```sh
jj bookmark move main --to @-
jj git push
```

Then release with `gh`.

```sh
gh release create v1.2.3 --title "v1.2.3" --notes "this is a public release"
```

Or use the short version, with interactive prompts.

```sh
gh release create v4.5.6
```

The workflow builds binaries for all platforms and creates a GitHub release.

Optionally fetch changes with `jj` to see the new tag on `main`.

### Record demo video with VHS

```nushell
docker run --rm -v ($env.PWD):/vhs ghcr.io/charmbracelet/vhs vhs/basic.tape
```
Or with custom Dockerfile (installs `go` for building the binary dynamically):

```nushell
docker run --rm -v ($env.PWD):/vhs (docker build -q /vhs) vhs/big.tape
```

## License

MIT
