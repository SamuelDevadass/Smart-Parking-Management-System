import os
from typing import Optional
import shutil
from datetime import datetime
from dotenv import load_dotenv

from services import detection
load_dotenv(".//.env")

def create_folder_path(license_number:str) -> Optional[str]:

    folder_path_raw = detection.get_status().get("folder_path", "")
    if folder_path_raw =="" or license_number  =="":
        print("Error fetching status")
        return ""
    date = datetime.now().strftime("%d-%m-%Y")
    time = datetime.now().strftime("%H-%M-%S")
    vehicle_time =f'capture_log/{date}/{license_number}_{time}'
    os.makedirs(f'capture_log/{date}', exist_ok=True)
    #os.makedirs(vehicle_time)
    shutil.copytree(folder_path_raw,vehicle_time, dirs_exist_ok=True)
    return vehicle_time



            