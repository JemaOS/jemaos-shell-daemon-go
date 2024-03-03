#!/bin/bash
[ $# -lt 3 ] && exit 1

. /usr/share/jemaos_shell/shell_lib.sh

data=$1
exdata=$2
msg=$3

emit_event $NOTIFY_TYPE_CUSTOM $data $exdata "${msg}"
