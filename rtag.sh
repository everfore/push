#!/bin/sh

git add .
git commit -m "auto"

osascript <<EOF
tell application "iTerm"
exec command "git status"
end tell
EOF