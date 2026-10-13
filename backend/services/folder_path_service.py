import os
from typing import Optional
from datetime import datetime
from dotenv import load_dotenv
from supabase import create_client, Client
from services import detection
load_dotenv(".//.env")

def create_folder_path(license_number: str) -> Optional[str]:
    folder_path_raw = detection.get_status().get("folder_path", "")
    if not folder_path_raw or not license_number:
        print("Error fetching detection status or license number")
        return ""
    supabase_url = os.getenv("SUPABASE_URL")
    supabase_key = os.getenv("SUPABASE_SECRET_KEY")
    bucket_name = os.getenv("S3_BUCKET")
    if not supabase_url or not supabase_key or not bucket_name:
        raise ValueError("Missing Supabase URL, secret key, or bucket name")
    supabase: Client = create_client(supabase_url, supabase_key)
    date = datetime.now().strftime("%d-%m-%Y")
    current_time = datetime.now().strftime("%H-%M-%S")
    vehicle_time = (f"capture_log/{date}/{license_number}_{current_time}")
    filenames = ["Captured_Image.jpg", "Car_Crop_Capture.jpg"]
    for filename in filenames:
        local_file = os.path.join(folder_path_raw, filename)
        output_file = f"{vehicle_time}/{filename}"
        with open(local_file, "rb") as f:
            supabase.storage.from_(bucket_name).upload(path=output_file,file=f,
                file_options={"cache-control": "3600", "upsert": "false", 
                              "content-type": "image/jpeg",},)
    # Return the folder prefix URL.
    # This is not a direct URL to either image.
    direct_path ="https://supabase.com/dashboard/project/smxmlxxaifxqckxnaqce/storage/files/buckets/parking-logs?path="+vehicle_time
    return direct_path
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
            