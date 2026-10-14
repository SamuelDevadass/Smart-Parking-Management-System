import os
from typing import Optional
from datetime import datetime
from dotenv import load_dotenv
from supabase import create_client, Client
from services import detection
load_dotenv(".//.env")

"""Function should take time and path as parameters from go if want to make it generic
even for entry take time from go
for exit no need to take time since we can modify get active session to return the supabase folder"""
def create_folder_path(license_number: str, action_subfolder: str, blob_folder) -> Optional[str]:
    folder_path_raw = detection.get_status().get("folder_path", "")
    if not folder_path_raw or not license_number:
        print("Error fetching detection status or license number")
        return ""
    if not action_subfolder:
        action_subfolder = "Entry"
    
    supabase_url = os.getenv("SUPABASE_URL")
    supabase_key = os.getenv("SUPABASE_SECRET_KEY")
    bucket_name = os.getenv("S3_BUCKET")
    if not supabase_url or not supabase_key or not bucket_name:
        raise ValueError("Missing Supabase URL, secret key, or bucket name")
    supabase: Client = create_client(supabase_url, supabase_key)

    if action_subfolder == "Entry":
        date = datetime.now().strftime("%d-%m-%Y")
        current_time = datetime.now().strftime("%H-%M-%S")
        vehicle_time = (f"capture_log/{date}/{license_number}_{current_time}")
    elif action_subfolder == "Exit" and blob_folder is not None:
        vehicle_time = blob_folder
    
    filenames = ["Captured_Image.jpg", "Car_Crop_Capture.jpg"]
    for filename in filenames:
        local_file = os.path.join(folder_path_raw, filename)
        if action_subfolder == "Entry":
            output_file = f"{vehicle_time}/{action_subfolder}/{filename}"
        elif action_subfolder == "Exit" and blob_folder is not None:
            output_file = f"{blob_folder}/{action_subfolder}/{filename}"
        with open(local_file, "rb") as f:
            supabase.storage.from_(bucket_name).upload(path=output_file,file=f,
            file_options={"cache-control": "3600", "upsert": "false", 
                                "content-type": "image/jpeg",},)

    print("Upload succeeded")
    direct_path = os.getenv("SUPABASE_WEB_VIEW_PREFIX")+vehicle_time
    return direct_path
    
        


# """Error fixed"""
# BRO 😂 that actually makes sense. The key clue is that restarting the Go terminal fixed the value.
# Your terminal was probably holding a stale environment variable.
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
            