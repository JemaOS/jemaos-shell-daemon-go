#!/bin/bash
# Ensure at least two arguments are provided, otherwise exit with status 1
[ $# -lt 2 ] && exit 1

# Source the shell library for shared functions
. /usr/share/jemaos_shell/shell_lib.sh

# Extract the level and message from the arguments
level=$1
msg=$2

# Emit a system event with the provided level and message
emit_event $NOTIFY_TYPE_SYSTEM -1 $level "${msg}"