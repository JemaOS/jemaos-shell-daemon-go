#!/bin/bash
NOTIFY_TYPE_SYSTEM=0  # Notification type for system events
NOTIFY_TYPE_COMMAND=1 # Notification type for command events
NOTIFY_TYPE_CUSTOM=2  # Notification type for custom events

# Function to emit an event via D-Bus
emit_event() {
  # Ensure exactly four arguments are provided, otherwise return with status 1
  [ $# -ne 4 ] && return 1
  local type=$1      # Type of the notification
  local handler=$2   # Handler for the notification
  local state=$3     # State of the notification
  local msg=$4       # Message to include in the notification

  # Send the event using dbus-send
  dbus-send --system --type=method_call \
   --dest=io.jemaos.ShellDaemon \
   /io/jemaos/ShellDaemon \
   io.jemaos.ShellInterface.EmitNotification \
   int32:$type \
   int32:$handler \
   int32:$state \
   string:"${msg}"
}