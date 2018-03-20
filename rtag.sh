#!/bin/sh

osascript <<EOF
tell application "iTerm2"
tell application "System Events"
    keystroke "git rebase -i HEAD~2"
    keystroke return
end tell
end tell
EOF