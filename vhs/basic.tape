# Configuration

Output vhs/basic.gif

Require countdown

Set Shell "bash"
Set FontSize 32
Set Width 1200
Set Height 300

Env PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/vhs"

# Script for recording

# Type './countdown --spinner "bomb" --range "10..0"      --title "Implosion in"' Sleep 1s
Type './countdown --spinner "bomb" --range "10..0" \' 
Enter
Type '  --title "Implosion in"' Sleep 1s
Enter Sleep 14s
