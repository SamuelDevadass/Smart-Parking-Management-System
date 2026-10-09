# Smart Parking Management System with License Plate Detection
### Project Evolution Roadmap

| Raspberry PI Version | Local Monolith | Modular Local Monolith | Application | Fullstack Web App | Fullstack Web App with Go | Dockerized Web App | 
| :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| [![Completed](https://img.shields.io/badge/Completed-a5d6a7?style=for-the-badge)](#) | [![Completed](https://img.shields.io/badge/Completed-a5d6a7?style=for-the-badge)](#) | [![Completed](https://img.shields.io/badge/Completed-a5d6a7?style=for-the-badge)](#) | [![Completed](https://img.shields.io/badge/Completed-a5d6a7?style=for-the-badge)](#) | [![Completed](https://img.shields.io/badge/Completed-a5d6a7?style=for-the-badge)](#) | <br><br>[![Current](https://img.shields.io/badge/Current-2e7d32?style=for-the-badge)](#) <br>Distributed Deployment | [![Upcoming](https://img.shields.io/badge/Upcoming-757575?style=for-the-badge)](#) | 

# FullStack Web Application with Go

<li>FullStack Web Application
<li>Database Normalized upto BCNF  
<li>Automatic License Plate Detection
<li>Revised API Gateway  
<li>Reverse proxy for AI workloads

### ARCHITECTURE DIAGRAM

<img width="982" height="457" alt="image" src="https://github.com/user-attachments/assets/abb04c1d-6a95-4990-8ac3-594ddbcd18b9" />


| Component | Tech Stack | Deployed At |
| :--- | :--- | :--- |
| **Frontend** | React | Vercel |
| **API Layer** | Go Huma | Railway |
| **Backend** | Reverse Proxy (Go to FastAPI) | Railway |
| **Database** | PostgreSQL + PGX | Supabase |
| **Capture** | Blob Storage | Supabase Blob with S3 configuration |       

## Video Walkthrough

[![Watch the Demo Video](https://markdown-videos-api.jorgenkh.no/youtube/7mYoQlRec2c)](https://youtu.be/7mYoQlRec2c)

## FEATURED UPDATES

API Layer : FastApi + Python --> Golang Huma + Reverse Proxy to FastAPI

Database  : Python + psycopg --> Golang + pgx

-----------------------------------------------------------------------------------------------------

## IMPROVEMENTS OVER FULLSTACK WEB APPLICATION

i. Create structured folder path ordered by date of capture files and log them in Database

ii. Added stric type based checking in compliance with standard Go norms

iii. Auto generated Huma OpenAPI documentation for efficient testing 

-----------------------------------------------------------------------------------------------------

### DATABASE

Normalized upto BCNF

Deployed to SupaBase, existed data exported to SupaBase DB

Implemented using PostGreSQL and PGX with Go

Designed to support organizations with multiple centres, wings, floors and spots
    
    Eg: Organization Seasons Mall
            Centre: Seasons Mall Parking Centre
                Wing:   Seasons Mall Parking Centre Wing 1
                Floors: B1,B2,B3
                Spots: Multiple spots for Two Wheelers and Four Wheelers Available per floor
Supports comprehensive logging for robust performance

### FRONTEND

Implemented using React 

Designed to support multiple pages dedicated to specific user activities
    
    Eg: Choose Wing
        Automatic License Number Plate Detection
        Input User Details
        Mark Entry/Exit
        Generate Bill 
Supports Interactive design for enhanced efficiency

### BACKEND 

Implemented using Go Huma with Reverse Proxy to FastAPI

Designed to support secure REST APIs for backend endpoints 

    Eg: /api/wings
        /api/spots/
        /api/detection/ --> Reverse Proxy to FastAPI
        /api/vehicles/
        /api/entries/
        /api/exits/
        /api/bills/
        /api/video/     --> Reverse Proxy to FastAPI  
Supports Fixed Input and Response models defined in Go for better security over HTTP messages

### AUTOMATIC LICENSE PLATE DETECTION

Integrated YOLO_Monolith for Automatic License Number Detection

Implemented using YOLO, OpenCV, EasyOCR and Pillow

Includes automatic fallbacks with custom image enhancement pipeline against confidence score checks


-----------------------------------------------------------------------------------------------------

### QUICK DEV PATCHES

#### Tools and Technologies

PostgreSQL, Pgx, Go Huma, YOLO, EasyOCR, OpenCV, Pillow, RapidOCR

#### CONFIGURATIONS

i. Folder Structure
    
    root/
    -   .env
    -   .gitignore
    -   README.md
    -   LICENSE
    -   requirements.txt
    -   requirements.bak
    +---.vscode/
           settings.json
    +---frontend/
        -   eslintconfig.js
        -   index.html
        -   package-lock.json
        -   package.json
        -   vite.config.js
        +---node_modules/
        +---public/
            -   favicon.svg
            -   icons.svg 
        +---src/
            -   App.jsx
            -   index.css
            -   main.jsx  
            +---api/
                -   client.js
            +---components/
                -   PlateChip.jsx
                    StepBar.jsx       
            +---pages/
                -   BillingPage.jsx
                -   DetectionPage.jsx
                -   EmptySpotsPage.jsx
                -   EntryExitPage.jsx
                -   OwnerDetailsPage.jsx
                -   SelectWingPage.jsx
    +---backend/
        -   LicensePlateDetector.py
        -   schemas.py
        -   yolo26n.pt
        +---routers/
            -   bills.py
            -   detection.py
            -   spots.py
            -   vehicles_entry_exit.py
            -   video.py
            -   wings.py
        +---scans/
            -   <timestampted folders of captures>
        +---services/
            -   billing.py
            -   db.py
            -   detection.py
    +---backend_go/
        +---catpure_log/
        +---handlers/   
            +---wings/
                -   wings.go
        -   bills.go
        -   detection.go
        -   spots.go
        -   vehicle_entry_exit.go
        -   video.go
        +---models/
        -   models.go
        -   wings.go
        +---repositories/
        -   db.go
        +---services/
        -   billing.go
        -   detection_status.go
    +---db/
    -   .sql
    -   ER DIAGRAM.jpg
    -   FINAL SCHEMA (BCNF).png       
    +---venv/

ii. Environment Variables
    
Configure following variables to configure the environment:
        
        DB_NAME=
        DB_USER=
        DB_PW=
        DB_HOST=
        FRONTEND_URL=
        VITE_BACKEND_URL=
        BACKEND_URL=
        DB_URL=

iii. Environment Setup
    
-> Clone the repository/branch using command: 

    https://github.com/SamuelDevadass/Smart-Parking-Management-System.git
    
-> Create venv using command:
    
    python -m create venv venv
    
-> Activate the venv
    
-> Navigate to backend and install requirements using command:

    pip install -r requirements.txt

  - (Incase required, install other requirements from requirements.bak file)

-> Navigate to frontend and install requirements using command:

    npm install

-> Implement the DB using postgresql with the SQL from ./db/.sql
    
   - Refer to diagrams in ./db and issues for detailed Schema and ER diagram
    
-> Run the application:

   - Navigate to the backend and run:

    uvicorn main:app --reload --port 8001

   - Navigate to the backend_go and run:

    go mod tidy
    go build -o main.exe main.go
    go run main.go (or ./main.exe)

   - Navigate to the frontend and run:

    npm run dev

#### PRO TIPS

i. Ensure the .gitignore contains the lines:

        frame*.jpg
        202*/
        202*/*.jpg
      
This ensures the folders created during YOLO capture cycles are not tracked by Git

ii. venv startup:
    
        -> Set-ExecutionPolicy -Scope Process -ExecutionPolicy Remote
    
        -> ./venv/scripts/Activate

iii. Pipreqs library 
    
To add only those libraries actually imported in the files run the command:

    python -m pipreqs.pipreqs . --force

# CITATION

If you use this project, please credit https://github.com/SamuelDevadass/Smart-Parking-Management-System

Citation: [Smart Parking Management System / FullstackWebApplication_Go], Samuel Devadass (2026). 

Available at: [https://github.com/SamuelDevadass/Smart-Parking-Management-System/tree/FullstackWebApplication_Go]
