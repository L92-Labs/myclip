package instance

// The running MyClip takes one command from later launches, over a socket
// (a named pipe on Windows) in the run dir: "toggle" shows or hides the
// popup. It reads nothing and changes nothing else, so a desktop shortcut
// can run `myclip --toggle` (Wayland compositors keep global keys for
// themselves).

// Toggle is the only command.
const Toggle = "toggle"
