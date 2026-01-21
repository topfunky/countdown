# Configuration

Output vhs/big.gif

Set Shell "bash"
Set FontSize 32
Set Width 1200
Set Height 500

Env PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/vhs"

# Setup
Hide
Type 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv bash)"'
Enter
Type "go build -o countdown ."
Enter
Wait
Type "clear"
Enter
Show

# Script for recording

Type './countdown --big --range 10..0 --spinner="bomb"' Sleep 1s
Enter Sleep 14s
