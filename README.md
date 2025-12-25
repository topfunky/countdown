# Countdown app

This simple command line application displays a spinner, a title, and a number which counts down in place by one every second until it reaches the final number (defaults to starting at `100` and ending at `0`).

```plain
:moon: Liftoff in 99
```

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
```
