from fastapi import APIRouter, Request, HTTPException
from services import detection
from services import folder_path_service
router = APIRouter()

@router.post("/api/detection/start")
async def start_detection(request: Request):
    image_bytes = await request.body()

    if not image_bytes:
        raise HTTPException(status_code=400, detail="No image received")

    detection.process_frame(image_bytes)
    return {"ok": True}


@router.get("/api/detection/status")
def detection_status():
    return detection.get_status()


@router.post("/api/detection/stop")
def stop_detection():
    detection.stop_detection()
    return {"ok": True}

# for blob path
@router.get("/api/detection/folder/{license_plate}")
def blob_path(license_plate: str, action_subfolder: str, blob_folder: str):
    return folder_path_service.create_folder_path(license_plate, action_subfolder, blob_folder)

