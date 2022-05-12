# GoKey
Simple proof of concept Windows Go keylogger via conventional Windows APIs (SetWindowsHookEx low level keyboard hook)

## Build
By default the built Go exe would launch with a console. You can make this a little stealthier by building with:
```
go build -ldflags -H=windowsgui
```

## Function
- Hooks keys pressed via `SetWindowsHookEx`
- Monitors and logs window changes (switching of active window) via `SetWinEventHook`
- Converts low level keyboard events to typed text via `ToAscii` + `GetKeyboardState` - this handles things like caps lock and the shift key for us
- Window change catching leverages `AttachThreadInput` to new focused windows to allow for acquisition of keyboard state from within the keylogger thread

`AttachThreadInput` is a method of allowing the keylogging thread to get an accurate keyboard state. Per documentation:

https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardstate
```
An application can call this function to retrieve the current status of all the virtual keys. The status changes as a thread removes keyboard messages from its message queue. The status does not change as keyboard messages are posted to the thread's message queue, nor does it change as keyboard messages are posted to or retrieved from message queues of other threads. (Exception: Threads that are connected through AttachThreadInput share the same keyboard state.)
```

Without this adaptation the keyboard state would always be equal to whatever it was upon launching the keylogger as it is should never be in focus and is just running in the background.

The current implementation only begins attaching after the first window switch post-launch of the keylogger. This means that capitalization etc... could be innacurate until the first window switch.

"Caching" of threads that have already had `AttachThreadInput` executed on them is done via a simple map. This means that a long running logging campaign could potentially have a thread id collision where a new thread gets the same id as an older thread - current implementation would not handle that very gracefully. Base keylogging would still work but capitalization etc... may be innacurate as well as the window title.

Window titles are obtained via `GetWindowText` - sometimes these return as empty and will be marked as `[Unknown]`.

## Exfiltration
A very simple example of exfiltration is included within the `callhome` package. It sends an initial `beacon` on startup with hostname and username. An example of using a timer to exfiltrate non-empty logs that also include hostname and username is also demonstrated.