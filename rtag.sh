#!/bin/sh

osascript <<EOF
tell application "iTerm2"
tell application "System Events"
    keystroke "git rebase -i HEAD~2"
    keystroke return
end tell
end tell

# # tell application "System Events"
# 	set t to application "iTerm2"
#   	tell t 
#   		# key down command
#   		keystroke "t" using command down

#   		# key up command
#   	end tell
# # end tell

# tell application "iTerm2"
#     activate
#     set t to (make new terminal)
#     tell t
#         tell (make new session at the end of sessions)
#             write text "pwd && ls"
#         end tell
#     end tell
# end tell

tell application "iTerm2"
	tell current window
		tell current session
			write text "ezhome nadesico"
		end tell
	end tell
end tell
EOF