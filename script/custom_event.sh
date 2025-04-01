#!/bin/bash
# Ensure at least three arguments are provided, otherwise exit with status 1
[ $# -lt 3 ] && exit 1

# Source the shell library for shared functions
. /usr/share/jemaos_shell/shell_lib.sh

# Extract the data, extra data, and message from the arguments
data=$1
exdata=$2
msg=$3

# Emit a custom event with the provided data, extra data, and message
emit_event $NOTIFY_TYPE_CUSTOM $data $exdata "${msg}"