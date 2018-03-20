#!/bin/sh

git add .
git commit -m "auto"
echo "======="
osascript <<EOF
tell application "iTerm2"
    activate
    set t to (make new terminal)
    tell t
        tell (make new session at the end of sessions)
            exec command "cd Downloads"
            exec command "clear"
            exec command "pwd"
        end tell
    end tell
end tell
EOF
echo "======="