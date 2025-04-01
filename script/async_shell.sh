#!/bin/bash
# Ensure at least two arguments are provided, otherwise exit with status 1
[ $# -lt 2 ] && exit 1

# Source the shell library for shared functions
. /usr/share/jemaos_shell/shell_lib.sh

# Define constants for event codes
ON_CLOSED=2  # Event code for a successful command execution
ON_ERROR=3   # Event code for a failed command execution

# Function to convert event codes to string representations
code_to_str() {
  if [ $1 -eq $ON_CLOSED ]; then
    echo "OnClosed"  # Return "OnClosed" for a successful execution
  else
    echo "OnError"   # Return "OnError" for a failed execution
  fi
}

# Extract the key and command from the arguments
key=$1
shift
command=$1
shift

# Execute the command with the remaining arguments
$command "$@"

# Check the exit status of the command
if [ $? -ne 0 ]; then
  returncode=$ON_ERROR  # Set return code to ON_ERROR if the command failed
else
  returncode=$ON_CLOSED # Set return code to ON_CLOSED if the command succeeded
fi

# Emit an event with the appropriate return code and its string representation
emit_event $NOTIFY_TYPE_COMMAND $key $returncode "$(code_to_str $returncode)"