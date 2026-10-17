import threading
from LicensePlateDetector import Detector

_detector = None            # lazy-init — loading YOLO/OCR models is slow
_stop_event = threading.Event()
_thread: threading.Thread | None = None

# Shared state read by GET /detection/status. Only ever written by
# run_detection() (in the background thread) and start_detection()
# (in the request thread) — simple enough that a lock isn't needed since
# each field is only ever written by one side at a time.
_state = {"status": "idle", "license_plate": "", "folder_path": ""}


def _run_detection():
    global _detector
    try:
        if _detector is None:
            _detector = Detector()
        text, folder_path = _detector.start(stop_event=_stop_event)
        _state["status"] = "done"
        _state["license_plate"] = (text or "").strip()
        _state["folder_path"] = folder_path
    except Exception as exc:  # noqa: BLE001 - surface any detector failure to the frontend
        print(f"Detection error: {exc}")
        _state["status"] = "error"


def start_detection():
    global _thread
    if _thread is not None and _thread.is_alive():
        return  # already running, ignore duplicate start
    _stop_event.clear()
    _state["status"] = "running"
    _state["license_plate"] = ""
    _thread = threading.Thread(target=_run_detection, daemon=True)
    _thread.start()


def stop_detection():
    _stop_event.set()


def get_status() -> dict:
    return dict(_state)

def get_live_frame() -> bytes | None:
    global _detector
    if _detector is None :
        return None
    return _detector.get_frame_bytes()
