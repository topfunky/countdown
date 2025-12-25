# Countdown app

This simple command line application displays a spinner, a title, and a number which counts down in place by one every second until it reaches the final number (defaults to starting at `100` and ending at `0`).

```plain
🌙 Liftoff in 99
```

## Installation

### Homebrew (macOS/Linux)

```sh
brew install OWNER/tap/countdown
```

### Go install

```sh
go install github.com/OWNER/countdown@latest
```

### Binary releases

Download pre-built binaries from the [Releases](https://github.com/OWNER/countdown/releases) page.

## Usage

```
Usage: countdown [flags]

Display spinner while displaying a number which counts downward

Flags:
  -h, --help                  Show context-sensitive help.
  -v, --version               Print the version number

  -s, --spinner="dot"         Spinner type ($COUNTDOWN_SPINNER)
                              dot|line|minidot|jump|pulse|points|globe|moon|monkey|meter|hamburger|none
      --title="Liftoff in"    Text to display to user while counting
                              ($COUNTDOWN_TITLE)
  -r, --range="100..0"        Numbers to count from and to
  -t, --time-interval=1       Number of seconds between each iteration
  -d, --decrement=1           Number subtracted from current count at
                              each iteration
  -f, --final-phase=5         Number at which the final phase starts. At this
                              number, the foreground and background colors are
                              swapped. Can be a number such as `5` or a
                              percentage such as `10%` 
  
Style Flags
  --spinner.foreground="212"    Foreground Color ($COUNTDOWN_FOREGROUND)
  --spinner.background=""       Background Color ($COUNTDOWN_BACKGROUND)
  --title.foreground=""         Foreground Color ($COUNTDOWN_FOREGROUND)
  --title.background=""         Background Color ($COUNTDOWN_BACKGROUND)
  --padding="0 0"               Padding ($COUNTDOWN_PADDING)

```

## Development

This CLI application is written with [Bubbletea](https://github.com/charmbracelet/bubbletea).

A Makefile implements many development tasks.

```sh
make test

make install-deps

make format

make snapshot    # Build snapshot release locally
```

### Releasing

Releases are automated via GitHub Actions. To create a new release:

1. Tag a commit with a version: `git tag v1.0.0`
1. Push the tag: `git push origin v1.0.0`

The workflow builds binaries for Linux, macOS, and Windows (amd64/arm64) and creates a GitHub release.

### Homebrew Tap Setup

To enable Homebrew installation, create a separate repository named `homebrew-tap` and set the `HOMEBREW_TAP_TOKEN` secret in this repository with a GitHub token that has write access to the tap repository. GoReleaser will automatically update the formula on release.
