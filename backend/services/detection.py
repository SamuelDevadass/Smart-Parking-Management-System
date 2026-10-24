from LicensePlateDetector import Detector

_detector = None


def get_detector() -> Detector:
    global _detector

    if _detector is None:
        _detector = Detector()

    return _detector


def start_detection() -> dict:
    return get_detector().start_session()


def process_detection_frame(image_bytes: bytes) -> dict:
    return get_detector().process_frame(image_bytes)


def stop_detection() -> dict:
    return get_detector().stop_session()


def get_status() -> dict:
    return get_detector().get_status()


def get_live_frame() -> bytes | None:
    return get_detector().get_frame_bytes()