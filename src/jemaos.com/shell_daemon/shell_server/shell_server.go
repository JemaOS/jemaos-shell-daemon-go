package shell_server

import (
    "fmt"
    "os"
    "runtime"
    "strings"

    te "jemaos.com/shell_daemon/shell_server/task_executor"
    "github.com/godbus/dbus/v5"
)

// Define the D-Bus path and interface constants
var DbusPath = dbus.ObjectPath("/io/jemaos/ShellDaemon")
const DbusIface = "io.jemaos.ShellInterface"
const ShellCommand = 1

// Debug related begin
// Debug flag to enable or disable debug prints
const debug = false

// Function to trace the caller function's name
func trace() string {
    pc, _, _, ok := runtime.Caller(1)
    if !ok {
        return "?"
    }

    fn := runtime.FuncForPC(pc)
    return fn.Name()
}

// Function to print debug messages if debug is enabled
func dPrintln(a ...interface{}) {
    if debug {
        fmt.Println(a...)
    }
}

// Debug related end

// DbusServer represents the D-Bus server
type DbusServer struct {
    dbus_ch chan *te.AsyncResult // Channel for asynchronous task results
    excutor *te.TaskList         // Task executor
    conn    *dbus.Conn           // D-Bus connection
}

// NewServer creates a new DbusServer instance
func NewServer(conn *dbus.Conn) *DbusServer {
    server := &DbusServer{
        make(chan *te.AsyncResult),
        te.NewTaskList(),
        conn,
    }
    go server.ListenAsyncCh()
    return server
}

// ListenAsyncCh listens for asynchronous task results and handles them
func (server *DbusServer) ListenAsyncCh() {
    for {
        select {
        case aResult := <-server.dbus_ch:
            dPrintln(trace(), aResult)
            server.ShellNotifying(ShellCommand,
                aResult.Key,
                aResult.Code,
                aResult.Msg)
        }
    }
}

// Define common errors and empty results
var ErrCommandNotFound = dbus.NewError("no command script found", nil)
var EmptyResult = &te.TaskResult{0, ""}

// SyncExec executes a script synchronously
func (server *DbusServer) SyncExec(script string) (*te.TaskResult, *dbus.Error) {
    dPrintln(trace(), script)
    args := strings.Fields(script)
    if len(args) < 1 {
        return EmptyResult, ErrCommandNotFound
    }
    ch := make(chan *te.TaskResult)
    go server.excutor.SyncExec(args, ch)
    result := <-ch
    dPrintln(trace(), result)
    return result, nil
}

// AsyncExec executes a script asynchronously
func (server *DbusServer) AsyncExec(script string) (*te.TaskResult, *dbus.Error) {
    args := strings.Fields(script)
    if len(args) < 1 {
        return EmptyResult, ErrCommandNotFound
    }
    ch := make(chan *te.TaskResult)
    go server.excutor.AsyncExec(args, ch, server.dbus_ch)
    result := <-ch
    return result, nil
}

// AsyncExec2 is a backward-compatible version of AsyncExec
func (server *DbusServer) AsyncExec2(script string) (*te.TaskResult, *dbus.Error) {
    return server.AsyncExec(script)
}

// GetTaskState retrieves the state of a task by its ID
func (server *DbusServer) GetTaskState(key int) (*te.TaskResult, *dbus.Error) {
    _, err := server.excutor.GetTask(key)
    if err != nil {
        return EmptyResult, nil
    }
    return &te.TaskResult{key, server.excutor.GetState(key)}, nil
}

// GetAsyncTaskOutput retrieves the output of an asynchronous task
func (server *DbusServer) GetAsyncTaskOutput(key int, lines int) (*te.TaskResult, *dbus.Error) {
    ch := make(chan *te.TaskResult)
    go server.excutor.GetAsyncTaskOutput(key, lines, ch)
    result := <-ch
    dPrintln(trace(), result)
    return result, nil
}

// GetDaemonState retrieves the state of the daemon
func (server *DbusServer) GetDaemonState() (*te.TaskResult, *dbus.Error) {
    return &te.TaskResult{server.excutor.GetCounter(), server.excutor.GetAllStates()}, nil
}

// ForceCloseTask forcibly closes a task by its ID
func (server *DbusServer) ForceCloseTask(key int) (*te.TaskResult, *dbus.Error) {
    server.excutor.RemoveTask(key)
    return server.GetTaskState(key)
}

// Exit shuts down the server and removes all tasks
func (server *DbusServer) Exit() {
    server.excutor.RemoveAllTasks()
    os.Exit(1)
}

// ShellNotifying sends a notification via D-Bus
func (server *DbusServer) ShellNotifying(s_type int, handler int, state int, msg string) error {
    dPrintln(trace(), s_type, handler, state, msg)
    return server.conn.Emit(DbusPath, DbusIface+".ShellNotifying", s_type, handler, state, msg)
}

// EmitNotification emits a notification and returns the result
func (server *DbusServer) EmitNotification(s_type int, handler int, state int, msg string) (int, *dbus.Error) {
    err := server.ShellNotifying(s_type, handler, state, msg)
    if err != nil {
        return -1, nil
    }
    return 0, nil
}