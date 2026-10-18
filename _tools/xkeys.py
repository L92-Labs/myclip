"""Press key chords on an X display through XTest: xkeys.py super+shift+c down return"""
import ctypes, sys, time
x = ctypes.CDLL("libX11.so.6"); t = ctypes.CDLL("libXtst.so.6")
x.XOpenDisplay.restype = ctypes.c_void_p
x.XStringToKeysym.restype = ctypes.c_ulong
x.XKeysymToKeycode.argtypes = [ctypes.c_void_p, ctypes.c_ulong]
t.XTestFakeKeyEvent.argtypes = [ctypes.c_void_p, ctypes.c_uint, ctypes.c_int, ctypes.c_ulong]
x.XFlush.argtypes = [ctypes.c_void_p]
d = x.XOpenDisplay(None)
names = {"super": "Super_L", "shift": "Shift_L", "ctrl": "Control_L", "alt": "Alt_L", "down": "Down", "up": "Up",
         "return": "Return", "escape": "Escape", "space": "space"}
def code(n): return x.XKeysymToKeycode(d, x.XStringToKeysym(names.get(n, n).encode()))
for chord in sys.argv[1:]:
    keys = [code(k) for k in chord.split("+")]
    for k in keys: t.XTestFakeKeyEvent(d, k, 1, 0)
    for k in reversed(keys): t.XTestFakeKeyEvent(d, k, 0, 0)
    x.XFlush(d); time.sleep(0.25)
