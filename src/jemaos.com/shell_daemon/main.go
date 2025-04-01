package main

import (
    "fmt"
    "os"

    "jemaos.com/shell_daemon/shell_server"
    "github.com/godbus/dbus/v5"
    "github.com/godbus/dbus/v5/introspect"
)

// Define the introspection XML for the D-Bus interface
const intro = `
<node>
    <interface name="io.jemaos.ShellInterface">
        <method name="SyncExec">
            <arg direction="in" type="s"/> <!-- Input: script as a string -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="AsyncExec">
            <arg direction="in" type="s"/> <!-- Input: script as a string -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="AsyncExec2">
            <arg direction="in" type="s"/> <!-- Input: script as a string -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="GetTaskState">
            <arg direction="in" type="i"/> <!-- Input: task ID as an integer -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="GetAsyncTaskOutput">
            <arg direction="in" type="ii"/> <!-- Input: task ID and number of lines -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="GetDaemonState">
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="ForceCloseTask">
            <arg direction="in" type="i"/> <!-- Input: task ID as an integer -->
            <arg direction="out" type="is"/> <!-- Output: integer and string -->
        </method>
        <method name="EmitNotification">
            <arg direction="in" type="iiis"/> <!-- Input: integers and a string -->
            <arg direction="out" type="i"/> <!-- Output: integer -->
        </method>
        <method name="Exit">
            <!-- No arguments for this method -->
        </method>
        <signal name="ShellNotifying">
            <arg type="iiis"/> <!-- Signal with integers and a string -->
        </signal>
    </interface>` + introspect.IntrospectDataString + `</node>`

func main() {
    // Connect to the system D-Bus
    conn, err := dbus.SystemBus()
    if err != nil {
        panic(err) // Exit if the connection fails
    }
    defer conn.Close() // Ensure the connection is closed on exit

    // Create a new D-Bus server
    server := shell_server.NewServer(conn)

    // Export the server and introspection data to the D-Bus
    conn.Export(server, shell_server.DbusPath, shell_server.DbusIface)
    conn.Export(introspect.Introspectable(intro), shell_server.DbusPath,
        "org.freedesktop.DBus.Introspectable")

    // Request a unique name on the D-Bus
    reply, err := conn.RequestName("io.jemaos.ShellDaemon", dbus.NameFlagDoNotQueue)
    if err != nil {
        panic(err) // Exit if the name request fails
    }
    if reply != dbus.RequestNameReplyPrimaryOwner {
        fmt.Fprintln(os.Stderr, "name already taken") // Print error if name is taken
        os.Exit(1) // Exit with an error code
    }

    // Print a message indicating the server is listening
    fmt.Println("Listening on", shell_server.DbusIface, shell_server.DbusPath, "...")

    // Block the main thread to keep the server running
    select {}
}