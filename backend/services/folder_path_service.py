import os
from typing import Optional
import shutil
from datetime import datetime
from dotenv import load_dotenv
import supabase

from services import detection
load_dotenv(".//.env")

#Upload to blob directly
def create_folder_path(license_number:str) -> Optional[str]:

    folder_path_raw = detection.get_status().get("folder_path", "")
    if folder_path_raw =="" or license_number  =="":
        print("Error fetching status")
        return ""
    date = datetime.now().strftime("%d-%m-%Y")
    time = datetime.now().strftime("%H-%M-%S")
    vehicle_time =f'capture_log/{date}/{license_number}_{time}'
    bucket_name = os.getenv("S3_BUCKET")
    supabase_url = os.getenv("SUPABASE_URL")
    for filename in ["Captured_Image.jpg", "Car_Crop_Capture.jpg"]:
        file = vehicle_time +"/"+ filename
        with open(file, "rb") as f:
            supabase.storage.from_(bucket_name).upload(file=f,
            path=file, 
            file_options={"cache-control": "3600", "upsert": "false", "content-type":"image/jpeg"})
    #cleanup
    blobURL = supabase_url + "/storage/v1/object/public/" + bucket_name+"/"+vehicle_time
    return blobURL

# def create_folder_path(license_number:str) -> Optional[str]:

#     folder_path_raw = detection.get_status().get("folder_path", "")
#     if folder_path_raw =="" or license_number  =="":
#         print("Error fetching status")
#         return ""
#     date = datetime.now().strftime("%d-%m-%Y")
#     time = datetime.now().strftime("%H-%M-%S")
#     vehicle_time =f'capture_log/{date}/{license_number}_{time}'
#     os.makedirs(f'capture_log/{date}', exist_ok=True)
#     #os.makedirs(vehicle_time)
#     shutil.copytree(folder_path_raw,vehicle_time, dirs_exist_ok=True)
#     return vehicle_time
            