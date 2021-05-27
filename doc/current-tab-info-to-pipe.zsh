#!/bin/zsh

set -e  # exit on error

# ensure workflow cache directory exists
mkdir -p "$alfred_workflow_cache"

# create a named pipe for the workflow to write to
pipe="${alfred_workflow_cache}/tab.$$.json"
mkfifo -m 0600 "$pipe"
# ensure pipe is deleted when script exits
trap "rm -f '$pipe'" EXIT INT TERM

# tell workflow to write tab deets (as Alfred JSON) to pipe
script="tell application id \"com.runningwithcrayons.Alfred\" to run trigger \"current-tab\" in workflow \"net.deanishe.alfred.firefox-assistant\" with argument \"$pipe\""
/usr/bin/osascript -e "$script" &!
# send JSON from pipe and to Alfred
cat "$pipe"
