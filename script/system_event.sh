#!/bin/bash
[ $# -lt 2 ] && exit 1

. /usr/share/jemaos_shell/shell_lib.sh

level=$1
msg=$2

emit_event $NOTIFY_TYPE_SYSTEM -1 $level "${msg}"
