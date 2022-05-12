package wintypes

var (
	//https://docs.microsoft.com/en-us/windows/win32/winauto/event-constants
	EVENT_OBJECT_FOCUS DWORD = 0x8005
	//https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-keyup
	WM_KEYUP DWORD = 0x101

	WINEVENT_OUTOFCONTEXT   DWORD = 4
	WINEVENT_SKIPOWNPROCESS DWORD = 2

	//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexa
	WH_KEYBOARD_LL = 13
)

// https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
type (
	//typedef int BOOL;
	BOOL int32
	//typedef unsigned char BYTE;
	BYTE byte
	//typedef unsigned long DWORD;
	DWORD uint32
	//typedef PVOID HANDLE;
	HANDLE uintptr
	//typedef HANDLE HHOOK;
	HHOOK HANDLE
	//typedef HANDLE HINSTANCE;
	HINSTANCE HANDLE
	//typedef HINSTANCE HMODULE;
	HMODULE HANDLE
	//typedef HANDLE HWND;
	HWND HANDLE
	//typedef long LONG;
	LONG int32
	/*
		#if defined(_WIN64)
		typedef __int64 LONG_PTR;
		#else
		typedef long LONG_PTR;
		#endif
	*/
	LONG_PTR uintptr
	//typedef LONG_PTR LPARAM;
	LPARAM LONG_PTR
	//typedef LONG_PTR LRESULT;
	LRESULT LONG_PTR
	//typedef UINT_PTR WPARAM;
	WPARAM uintptr
	//https://docs.microsoft.com/en-us/windows/win32/winauto/hwineventhook
	//typedef HANDLE HWINEVENTHOOK;
	HWINEVENTHOOK HANDLE
	//typedef BYTE *PBYTE;
	PBYTE []BYTE
	/*
		https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-hookproc
		LRESULT Hookproc(
			int code,
			[in] WPARAM wParam,
			[in] LPARAM lParam
		)
	*/
	HOOKPROC func(int, WPARAM, LPARAM) LRESULT
	/*
		https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-wineventproc
		void Wineventproc(
			HWINEVENTHOOK hWinEventHook,
			DWORD event,
			HWND hwnd,
			LONG idObject,
			LONG idChild,
			DWORD idEventThread,
			DWORD dwmsEventTime
		)
	*/
	WINEVENTPROC func(HWINEVENTHOOK, DWORD, HWND, LONG, LONG, DWORD, DWORD) uintptr
)

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
typedef struct tagMSG {
  HWND   hwnd;
  UINT   message;
  WPARAM wParam;
  LPARAM lParam;
  DWORD  time;
  POINT  pt;
  DWORD  lPrivate;
} MSG, *PMSG, *NPMSG, *LPMSG;
*/
type MSG struct {
	Hwnd     HWND
	Message  uint32
	WParam   WPARAM
	LParam   LPARAM
	Time     DWORD
	Pt       POINT
	LPrivate DWORD
}

/*
https://docs.microsoft.com/en-us/previous-versions/dd162805(v=vs.85)
typedef struct tagPOINT {
  LONG x;
  LONG y;
} POINT, *PPOINT;
*/
type POINT struct {
	X, Y LONG
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
typedef struct tagKBDLLHOOKSTRUCT {
  DWORD     vkCode;
  DWORD     scanCode;
  DWORD     flags;
  DWORD     time;
  ULONG_PTR dwExtraInfo;
} KBDLLHOOKSTRUCT, *LPKBDLLHOOKSTRUCT, *PKBDLLHOOKSTRUCT;
*/
type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}
