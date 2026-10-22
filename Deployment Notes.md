# Smart Parking Management System — Deployment Documentation

# 1. Vercel Frontend Deployment

## 1.1 Tasks, Commands and Configurations

### Repository configuration

The frontend was deployed from the GitHub repository:

- Repository: `SamuelDevadass/Smart-Parking-Management-System`
- Branch: `FullstackWebApplication_Go`
- Root directory: `frontend`
- Framework: React with Vite

The root directory was configured explicitly because the frontend resides in a subdirectory of the repository.


The deployed production URL was: `https://parking-system-eta.vercel.app`

![alt text](<Screenshot 2026-10-22 201430.png>)


## 1.2 The Vercel Routing Issue and Fix

### Original problem

The original Vercel configuration forwarded API requests to the Go backend deployed on Railway.

```
{
  "rewrites": [
    {
      "source": "/api/:path*",
      "destination": "https://perceptive-fulfillment-production-18aa.up.railway.app/api/:path*"
    }
  ]
}
```

The general rewrite worked for requests such as:

```
/api/wings
```

However, the nested available-spots route returned a Vercel-generated HTML `NOT_FOUND` response containing a `bom1::...` identifier.


### Fix: Explicit rewrite for available spots

A specific rewrite was added before the general catch-all rewrite.

The deployed configuration discussed during debugging was:

```
{
  "rewrites": [
    {
      "source": "/api/:wing/spots/available/",
      "destination": "https://perceptive-fulfillment-production-18aa.up.railway.app/api/:wing/spots/available/"
    },
    {
      "source": "/api/:path*",
      "destination": "https://perceptive-fulfillment-production-18aa.up.railway.app/api/:path*"
    }
  ]
}
```

The explicit rule forwards the available-spots path to the corresponding Railway endpoint, while the general rule handles the remaining API paths.

# Railway Deployment

![alt text](<Screenshot 2026-10-20 180428.png>)

# 2. Railway Python Backend Deployment

## 2.1 Purpose and Architecture

The Python backend provides the license-plate detection and OCR functionality.

It was deployed as a separate FastAPI service on Railway. The Go backend communicates with it over HTTP, forwarding detection requests and receiving processed results.

Observed production service URL:

`https://smart-parking-management-system-production-feeb.up.railway.app`

![alt text](<Screenshot 2026-10-22 201208.png>)

## 2.2 Making the Backend Truly Headless

The original Python implementation included GUI-related dependencies and components, including:

- Tkinter.
- `ImageTk` from Pillow.
- GUI-oriented image display and camera interface code.
- OpenCV functionality that was initially installed through a GUI-capable package.

A cloud deployment environment such as Railway does not provide a normal desktop display environment. GUI dependencies and calls that expect a desktop can therefore cause import failures or runtime errors.

The backend was adapted to operate without requiring a desktop interface.

### Changes made

| Original concern                   | Deployment change                                               |
| ---------------------------------- | --------------------------------------------------------------- |
| GUI-based image display            | Avoid GUI display operations in the server-side processing path |
| Tkinter and `ImageTk` dependencies | Remove them from the runtime path where no GUI is needed        |
| GUI-capable OpenCV package         | Replace it with `opencv-python-headless`                        |
| Interactive desktop workflow       | Expose detection and processing through FastAPI endpoints       |
| Local camera/display assumptions   | Process incoming image data through HTTP requests               |

The intended deployment workflow became:

1. The frontend sends a detection request or image frame.
2. The Go backend forwards the request to the Python service.
3. Python processes the image using the detection and OCR pipeline.
4. The Python service returns the result to Go.
5. Go returns the response to the frontend.

This design avoids depending on a physical desktop display on the server.

### OpenCV dependency problem

The Python environment initially encountered an OpenCV dependency conflict involving packages that expected `opencv-python`.

The final Railway build command used during the successful deployment was:

```
pip install -r requirements.txt && pip uninstall -y opencv-python opencv-contrib-python opencv-python-headless && pip install --no-deps opencv-python-headless==5.0.0.93 && python -c "import cv2; print(cv2.__version__)"
```

This command:

1. Installs the declared Python requirements.
2. Removes the conflicting OpenCV distributions.
3. Installs the headless OpenCV package without resolving its dependencies again.
4. Verifies that `cv2` imports successfully.

## 2.3 ONNX Runtime and Memory Constraints

The Python service performs computer-vision operations involving YOLO-based vehicle/license-plate detection and OCR.

These workloads can consume considerably more memory and CPU than ordinary API requests. The Railway service had a resource allocation of:

- 1 GB memory.

This limited memory allocation made the inference pipeline a significant deployment consideration.

### Why move toward ONNX

ONNX (Open Neural Network Exchange) provides a model representation that can be executed using compatible inference runtimes.

The project moved toward compressed models and ONNX runtimes for YOLO and OCR to make the inference pipeline more suitable for the resource-constrained deployment environment.

The goals were to:

- Reduce model storage and loading overhead where possible.
- Avoid unnecessary heavyweight runtime dependencies.

The OCR pipeline used `rapidocr_onnxruntime`. The project also adopted ONNX-based inference for the YOLO component.

## 2.4 Railway Configurations

The Python service settings included:

| Setting                                | Value                                            |
| -------------------------------------- | ------------------------------------------------ |
| Source repository                      | `SamuelDevadass/Smart-Parking-Management-System` |
| Branch                                 | `FullstackWebApplication_Go`                     |
| Root directory                         | `/backend`                                       |
| Builder                                | Railpack                                         |
| Python version shown in build settings | Python 3.13.16                                   |
| Region                                 | US East (Virginia, USA)                          |
| Replicas                               | 1                                                |
| Memory limit                           | 1 GB                                             |
| CPU limit                              | 2 vCPU                                           |


### Dynamic port assignment

The final start command was:

```
uvicorn main:app --host 0.0.0.0 --port $PORT
```

The `$PORT` environment variable allows Railway to supply the appropriate port at runtime. The `0.0.0.0` host binding makes the server reachable through the container's network interface instead of limiting it to loopback.


The start command is important because it binds the application to all network interfaces and uses Railway's assigned port.

## 2.5 Outbound IPv6

The Railway service settings included an Outbound IPv6 option.

It was enabled to allow the service to make outbound connections to IPv6 destinations. This was relevant to the deployment's network-connectivity investigation, particularly connections involving Supabase.


## 2.6 Deployed Documentation

The Python FastAPI service's `/docs` endpoint was used to verify that the service had started and exposed its API documentation.

The service's documentation URL was:

`https://smart-parking-management-system-production-feeb.up.railway.app/docs`

![alt text](<Screenshot 2026-10-22 195516.png>)

# 3. Railway Go Backend Deployment

## 3.1 Purpose and Architecture

The Go backend uses Huma to expose the primary application API.

Its responsibilities include:

- Serving API endpoints used by the React frontend.
- Accessing Supabase PostgreSQL.
- Handling parking-related application logic.
- Communicating with the Python detection service over HTTP.
- Returning results to the frontend.

The Go backend is the primary API gateway for application requests, while the Python service provides specialized computer-vision processing.

## 3.2 Railway Configurations

The Go service was deployed from the same repository and branch as the frontend and Python service.

| Setting                            | Value                                                   |
| ---------------------------------- | ------------------------------------------------------- |
| Source repository                  | `SamuelDevadass/Smart-Parking-Management-System`        |
| Branch                             | `FullstackWebApplication_Go`                            |
| Root directory                     | `/backend_go`                                           |
| Builder                            | Railpack                                                |
| Go version shown in build settings | Go 1.27.1                                               |
| Region                             | US East (Virginia, USA)                                 |
| Replicas                           | 1                                                       |
| Memory limit                       | 1 GB                                                    |
| CPU limit                          | 2 vCPU                                                  |


The Go backend was exposed through a Railway-generated public domain:

`https://perceptive-fulfillment-production-18aa.up.railway.app`

![alt text](<Screenshot 2026-10-22 201309.png>)

### Build and startup

The configured build command was:

```
go build -o main.exe main.go
```

The start command was:

```
./main.exe
```
## 3.3 Port Configuration Problem

The Go server was configured to obtain its port from the `PORT` environment variable and bind to all interfaces.

The relevant server-address pattern was:

```
":" + port
```

The server then listens on the configured address rather than assuming that the local development port is also the production port.


## 3.4 Outbound IPv6 and Supabase Connectivity

The Railway service's Outbound IPv6 option was enabled during the database connectivity investigation.

The initial direct Supabase database connection was not reachable from Railway because the connection attempt encountered an unreachable IPv6 route.

## 3.5 Deployed Documentation

The Go Huma `/docs` endpoint was used to verify that the service started successfully and exposed its API routes.

![alt text](<Screenshot 2026-10-22 195826.png>)

# 4. Supabase Database Deployment and Usage

## 4.1 Role in the Architecture

Supabase provides PostgreSQL database services and object storage for the application.

The database stores application data required for the parking-management system, while object storage provides a place to store image-related data where configured.

![alt text](<Screenshot 2026-10-22 201614.png>)

## 4.2 Database

A strict BCNF Normalized database was deployed to SupaBase along with the existing data present on the local development of the Application & FullStack Web Application versions of the project. 

Supabase auto generated a Schema diagram of the relations / tables based on the SQL code provided.

![alt text](db/supabase/SUPABASE_SCHEMA.svg)

## 4.3 Object Storage

The project also used Supabase object storage with S3-compatible configuration.

The parking_log table has an image_folder_path attribute / column which stores the url of the image files uploaded to the Object Storage.

![alt text](<Screenshot 2026-10-22 201701.png>)

The storage uses an S3 bucket configuration, with a structured path for efficient indexing & retrieval of image files if required. 

It follows a hierarchical structure of Date -> License Number with Timestamp of Entry -> Entry & Exit 

![alt text](<Screenshot 2026-10-22 201742.png>)

![alt text](<Screenshot 2026-10-22 201810.png>)


# 5. Analytics and Observability

The following dashboard charts and Analytics were observed on the respective platforms after deployment

## 5.1 Vercel

![alt text](image-3.png)

![alt text](image-4.png)

- **CDN Requests**: Shows request volume and HTTP response status codes, indicating traffic patterns and request outcomes.

- **Fast Data Transfer**: Recorded approximately 9 GB outgoing and 4 MB incoming traffic, reflecting substantial outbound data transfer.


## 5.2 Railway

![alt text](image-5.png)

- **CPU**: *Python* has larger spikes, reaching roughly 0.27 vCPU;
           *Go's* CPU usage stays comparatively low.

- **Memory**: *Python* uses around 530–850 MB for much of the period
              *Go's* usage is minimal by comparison.

- **Network egress**: Both services show bursts of outbound traffic, with *Go* producing some of the largest spikes due to *Database* and *Object Storage* transfers.

## 5.3 Supabase

![alt text](image-6.png)

- **Max CPU usage**: 6.4% — Low CPU utilization, suggesting database processing demand was modest.

- **Disk I/O consumed**: 1% — Very little of the available disk I/O capacity was consumed.

- **Disk I/O remaining**: 99% — Substantial I/O capacity remained available.

- **Memory usage**: 51.4% — Around half of the available memory was in use, relatively stable across the period.

# TEARDOWN

The deployments and all services were stopped and taken down to prevent incurring extra costs since they all utlizied free tier trials on the respective platforms. 

The deployment experiment demonstrated the practical challenges of integrating multiple cloud services, managing deployment-specific configuration, debugging network and routing behavior, and balancing computer-vision workloads against constrained compute resources.

Refer to the video walkthrough for a live sample execution: 

[![Watch the Demo Video](https://markdown-videos-api.jorgenkh.no/youtube/gi900gtaSpA)](https://www.youtube.com/watch?v=gi900gtaSpA)