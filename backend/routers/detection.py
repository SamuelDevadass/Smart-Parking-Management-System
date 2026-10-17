
from fastapi import APIRouter, HTTPException, Request

from LicensePlateDetector import Detector

router = APIRouter()
detection = Detector()


@router.post("/api/detection/start")
def start_detection():
    try:
        return detection.start_session()
    except RuntimeError as exc:
        raise HTTPException(status_code=409, detail=str(exc))
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@router.post("/api/detection/frame")
async def process_detection_frame(request: Request):
    image_bytes = await request.body()

    try:
        return detection.process_frame(image_bytes)
    except RuntimeError as exc:
        raise HTTPException(status_code=409, detail=str(exc))
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=str(exc))
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@router.post("/api/detection/stop")
def stop_detection():
    try:
        return detection.stop_session()
    except RuntimeError as exc:
        raise HTTPException(status_code=409, detail=str(exc))
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc))


@router.get("/api/detection/status")
def get_detection_status():
    return detection.get_status()
