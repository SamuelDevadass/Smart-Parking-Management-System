1. VERCEL FRONTEND DEPLOYMENT

TASKS, COMMANDS, CONFIGURAIONS

DESCRIBE THE VERCEL.JSON ISSUE AND FIX

DEPLOYED DASHBOARD


ANALYTICS


2. RAILWAY PYTHON

Make the backend truly headless, how we removed / replaced GUI based libraries for headless versions

CONFIGURATIONS 

DESCRIBE THE NEED TO MOVE TO ONNX & MEMORY PROBLEM
GET DYNAMIC PORT & ASSIGN IT
ALLOW OUTBOUND IPV6

DEPLOYED DOCS

ANALYTICS

3. RAILWYAY GO 

CONFIGS

PORT PROBLEM

OUTBOUT IPV6

GO REROUTING FOR THE / AFTER .app


4. SUPABASE DB DEPLOYMENT AND USAGE 


5. ANALYTICS & OBSERVABILITY





Absolutely 😭 No celebrations until we have the deployment and debugging log written down. Let's make this a proper checklist you can keep for your M.Tech project documentation, future debugging, and deployment reproducibility.

I'll separate confirmed fixes, workarounds/configuration changes, and verification checks so we don't accidentally record a suspected cause as a proven fact.

# Smart Parking Management System — Deployment & Fix Log

Scope: React/Vite frontend · Go Huma backend · Python FastAPI detection service · Supabase PostgreSQL

## 1. Deployment setup

- Configure the Vercel project root directory.
  - Set the Vercel project's root directory to `frontend`, since the React/Vite frontend lives there.
  - This ensures Vercel builds the intended application rather than trying to build from the repository root.
- Deploy the Go Huma backend to Railway.
  - Confirm that the backend starts successfully in the production environment.
  - Confirm that the Huma API documentation endpoint (`/docs`) is accessible.
- Deploy the Python FastAPI detection/OCR service to Railway.
  - Confirm that the service starts and its `/docs` endpoint is accessible.
  - Keep this as a separate service from the Go backend.
- Configure the production API routing through Vercel.
  - The frontend uses relative `/api/...` requests.
  - Vercel forwards API requests to the Go backend hosted on Railway.
  - This avoids requiring each frontend component to hardcode the Railway backend URL.
- Make backend startup compatible with Railway's port configuration.
  - The Go server listens using the configured `PORT` environment variable.
  - The server binds to `":" + port` rather than assuming a fixed local port.
- Make the Go backend's `.env` file optional in production.
  - Local development can use `.env`.
  - Production configuration is supplied through Railway environment variables, so the application doesn't fail merely because a local `.env` file isn't present.

## 2. Python detection service fixes

- Resolve the OpenCV dependency issue.
  - The deployed Python service had an OpenCV dependency problem.
  - The dependency was changed to `opencv-python-headless==5.0.0.93`.
  - This allowed the service to run in the headless Railway environment.
- Fix the Python service URL's trailing-slash issue.
  - The Python service URL configured for the Go backend had a trailing slash.
  - Concatenating that URL with a route beginning with `/` could produce a double slash in the resulting URL.
  - The URL was normalized in Go using:
  ```
  strings.TrimRight(pythonURL, "/")
  ```
  - Alternatively, the trailing slash can be removed from the environment variable itself.
  - Record the chosen approach: confirm whether URL normalization in Go or changing the environment variable is the final configuration.
- Verify the detection/OCR integration.
  - The Python detection endpoint is `POST /api/detection/start`.
  - The Go backend can communicate with the Python service after the URL issue was addressed.
  - Detection and OCR were subsequently confirmed to work.

## 3. Supabase PostgreSQL connectivity

- Identify the database connectivity problem from Railway.
  - The initial direct Supabase database endpoint was not reachable from Railway because the connection attempted to use an IPv6 route that was unreachable in that environment.
- Identify the Supabase Session Pooler as the alternative to try.
  - The Session Pooler was recommended as a fallback for database connectivity.
  - The available log confirms that this was the proposed workaround, but does not independently establish the final connection configuration. Verify the actual Railway database URL before recording this as fully resolved.
- Document the final database connection settings.
  - Record whether the application uses the Supabase direct connection or Session Pooler.
  - Record the connection mode and relevant non-secret configuration.
  - Never include database passwords, API secrets, or credentials in the project log.

## 4. Frontend API client

- Use a central API client instead of scattered `fetch()` calls.
  - Frontend components call the shared `Api` client.
  - API requests use relative paths under `/api`, allowing Vercel to route them to Railway.
- Configure JSON request headers centrally.
  - The client sets `Content-Type: application/json` by default and permits request-specific headers to override or extend them.
- Add meaningful HTTP error handling.
  - The client checks `res.ok`.
  - For unsuccessful responses, it reads the response body and constructs an error containing the response status.
  - This helps distinguish an HTTP error from a successful response.
  The relevant pattern is:
  ```
  if (!res.ok) {
    const message = await res.text();
    const error = new Error(
      message || `Request failed: ${res.status}`
    );
    error.status = res.status;
    throw error;
  }

  return res.json();
  ```
- Keep API paths compatible with the Vite development proxy and production routing.
  - Use relative `/api/...` URLs in the frontend.
  - Configure the development proxy for local development and Vercel rewrites for production.
  - Avoid hardcoding a production host into individual page components.

## 5. Vercel rewrite issue — the important production fix

Confirmed working

- Verify the general API rewrite.
  - The original rewrite in `frontend/vercel.json` forwarded `/api/:path*` to the Go backend on Railway.
  - The `/api/wings` request worked through Vercel, confirming that the general routing rule handled at least that endpoint.
- Investigate the nested available-spots route.
  - A nested available-spots request returned Vercel-generated HTML containing `NOT_FOUND` and a `bom1::...` identifier.
  - The same endpoint worked when requested directly against Railway.
  - This distinction suggested that the failing request was being rejected or misrouted at the Vercel layer, rather than failing inside the Go handler.
  - This was a strong diagnostic clue, not definitive proof of the exact internal cause.
- Add a specific rewrite for the available-spots route.
  - An explicit rewrite was added before the general `/api/:path*` rule in `frontend/vercel.json`.
  - The specific rule forwards the wing-specific available-spots path to the corresponding Railway endpoint.
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
  - Important: preserve the actual working configuration from the deployed repository. The snippet above documents the configuration discussed during debugging; don't replace a working file blindly.
- Deploy the rewrite change and test again.
  - After deploying the explicit route, the available-spots request started working.
  - The exact reason Vercel failed to handle the nested path with the original rule remains unconfirmed.
  - The practical fix is confirmed by the successful end-to-end test.
- Keep the explicit rewrite before the catch-all rewrite.
  - This preserves the intended order of route matching in the configuration that resolved the issue.

## 6. Available-spots API and query parameters

- Verify the available-spots URL format.
  - The frontend builds a URL containing the encoded wing name and the query parameters `centre_id` and `size`.
  - Example structure:
  ```
  /api/{wing}/spots/available/?centre_id=3&size=Two%20Wheeler
  ```
- Verify that the direct Railway endpoint works.
  - The exact available-spots request worked directly against Railway before the Vercel-specific fix.
  - This helped isolate the production routing problem from the backend endpoint itself.
- Verify that the production flow works after the rewrite fix.
  - The user subsequently confirmed that the application was working end to end.
  - Record the successful production test as the final validation of the routing fix.

## 7. Vehicle entry, exit, and active-session handling

- Distinguish a backend application error from a proxy error.
  - The missing available-spots route produced Vercel HTML.
  - The vehicle spot lookup, by contrast, returned a JSON 404 from the Go backend.
  - This distinction was useful: the first pointed toward routing, while the second showed that the request reached the application.
- Interpret the active-session lookup carefully.
  - The request to `/api/vehicles/spot/MH-12-UY-0916` returned `No session details found for license plate ...`.
  - This can be expected when the vehicle has no active parking session.
  - A 404 is not automatically a bug; the response must be judged against the expected application state.
- Verify vehicle type display.
  - The vehicle type was confirmed to display correctly after the backend response/model shape was adjusted.
- Verify entry/exit and detection/OCR workflows.
  - The latest overall status is that the user has confirmed the application is working.
  - For the project log, retain this as an end-to-end success confirmation. If you want a formal test record, capture separate evidence for successful entry, exit, and billing rather than treating one successful flow as proof of every individual path.

## 8. Final production smoke-test checklist

Use this checklist whenever you redeploy the app or change API routing.

## Production validation

0 / 17

Deployment

Vercel builds from the frontend directory

Go backend starts on the Railway-provided port

Go Huma /docs endpoint is accessible

Python FastAPI /docs endpoint is accessible

Backend can connect to Supabase PostgreSQL

Parking selection

Wings load through the production frontend

Selecting a wing returns the correct centre_id

Available spots load for the selected centre_id and vehicle size

No Vercel-generated NOT_FOUND response appears for nested API routes

Vehicle workflow

Vehicle details and vehicle type display correctly

Detection and OCR work through the deployed services

Vehicle entry succeeds and creates the expected active session

Exit retrieves the active session and occupied spot

Exit and billing complete as expected

Regression checks

Unknown or inactive vehicle sessions return an understandable response

API errors remain distinguishable from proxy/routing errors

A fresh deployment preserves the working Vercel rewrites

Reset checklistCopy checklist

## 9. Lessons learned — record these for future debugging

- Test the same endpoint at multiple layers. Compare the frontend, Vercel URL, and direct Railway URL to isolate which layer is failing.
- Inspect the response body, not just the status code. Vercel-generated HTML and a Go-generated JSON error can both be 404s but indicate very different problems.
- Check URL construction carefully. Trailing slashes and path concatenation can create malformed upstream URLs.
- Check query parameters and their types. Parameters such as `centre_id` must be valid integers, not `undefined` or an incorrectly encoded value.
- Keep frontend routing and backend routing separate in your diagnosis. An endpoint working directly on Railway doesn't guarantee that Vercel forwards the same path successfully.
- Preserve working configuration changes. Once the explicit Vercel rewrite resolved the issue, avoid removing it during unrelated cleanup without testing the nested endpoint again.
- Treat logs as evidence. Record the observed error, the change made, and the result. Mark a suspected root cause as a hypothesis unless it has been verified.

## 10. Deployment summary for your project report

### Production Deployment and Debugging Summary

The Smart Parking Management System was deployed using a React/Vite frontend on Vercel, a Go Huma backend and a Python FastAPI detection/OCR service on Railway, and PostgreSQL on Supabase.

During deployment, several configuration and integration issues were addressed. The Vercel project root was configured to point to the frontend directory, the Go server was configured to use the deployment platform's assigned port, and the backend was made compatible with environments where a local `.env` file is unavailable. The Python service's OpenCV dependency and the URL used for communication between the Go and Python services were also addressed.

A major production issue involved routing requests for available parking spots through Vercel. The endpoint worked directly against Railway, but the corresponding nested request through Vercel returned a platform-generated `NOT_FOUND` response. An explicit rewrite for the available-spots route was added before the general API rewrite. Following redeployment, the route worked successfully.

The frontend API client was kept centralized and configured to use relative `/api` paths, with HTTP error handling to help distinguish unsuccessful responses from successful ones. Additional checks covered wing selection, centre identification, spot availability, vehicle details, detection/OCR, and the entry/exit workflow.

The final production test confirmed that the application was working. The exact internal reason for the original Vercel routing failure was not conclusively established, so the explicit rewrite is documented as the verified practical fix rather than a proven explanation of Vercel's behavior.

One final documentation task: save the actual deployed `frontend/vercel.json`, the final dependency file for Python, and a sanitized list of Railway environment-variable names alongside this checklist. Those are the pieces you'll most want if you need to reproduce the deployment later. Never put secret values in the log.

Now we're documenting the victory instead of merely celebrating it.

---------------------------------------------------------------------------------------------------------------

Python headless error

GO build commands and http serve

use auto assigned port numbers never personal convenience port numbers




Smart-Parking-Management-System: Go-HumaChi
Deployments
Variables
Metrics
Console
Settings
Filter settings
Filter Settings...

/
Source
Source Repo
SamuelDevadass/Smart-Parking-Management-System



Disconnect
Root Directory
Configure where we should look for your code. Docs↗
Root directory
/backend_go
Branch connected to production
Changes made to this GitHub branch will be automatically pushed to this environment.
FullstackWebApplication_Go

Disconnect
Auto deploy unavailable
Could not load branches.Retry
Wait for CI
Trigger deployments after all GitHub actions have completed successfully.

Networking
Public networking
Access your application over HTTP with the following domains
perceptive-fulfillment-production-18aa.up.railway.app
Port 8080





Domain
perceptive-fulfillment-production-18aa
.up.railway.app

8080
Update your domain or target port


Cancel

Update

Generate Domain

Custom Domain

TCP Proxy
Private Networking
Communicate with this service from within the Railway network.
perceptive-fulfillment.railway.internal
IPv4 & IPv6


Ready to talk privately ·
You can also simply call me
perceptive-fulfillment

DNS
perceptive-fulfillment
.railway.internal

Endpoint name available!


Cancel

Update
Outbound IPv6
Enable your service to make outbound connections to IPv6 destinations.

Edge
Under Attack Mode
Being DDoSed? One click puts a browser check in front of your domains, shedding bot floods at the edge. Real visitors solve it once and browse normally.
Until turned off

Activate
Takes effect across the globe in ~20 seconds.

CDN Caching
Cache static assets and optionally HTML at the edge. Reduces latency and origin load for your service's domains. Changes are staged and take effect on deploy without restarting the service.

Edge Rules
Block, redirect, challenge, or override caching for matching requests before they reach your service. Rules run in order; drag to reorder them in the editor.
Edge rules are not available on your plan. Please upgrade to use them.

Scale
Regions & Replicas
Deploy replicas per region for horizontal scaling.
US East (Virginia, USA)

Replicas
1
Replica
Multi-region replicas are only available on the Pro plan.

Learn More↗
Replica Limits
Allocate a maximum vCPU and Memory for each replica.
CPU: 2 vCPU

Plan limit: 2 vCPU


Memory: 1 GB

Plan limit: 1 GB


Upgrade for higher limits
Build
Builder

Railpack

Default

1.27.1go@1.27.1
App builder developed by Railway. Docs↗

Custom Build Command
Override the default build command that is run when building your app. Docs↗
Build command
go build -o main.exe main.go
Watch Paths
Gitignore-style rules to trigger a new deployment based on what file paths have changed. Docs↗
Add pattern
Add pattern e.g. /backend_go/src/**

Deploy
Custom Start Command
Command that will be run to start new deployments. Docs↗
Start command
./main.exe
Add pre-deploy step (Docs↗)
Teardown
Configure old deployment termination when a new one is started. Docs↗

Cron Schedule
Run the service according to the specified cron schedule.

Add Schedule
Healthcheck Path
Endpoint to be called before a deploy completes to ensure the new deployment is live. Docs↗

Healthcheck Path
Serverless
Containers will scale down to zero and then scale up based on traffic. Requests while the container is sleeping will be queued and served when the container wakes up. Docs↗

Restart Policy
Configure what to do when the process exits. Docs↗
On Failure

Restart the container if it exits with a non-zero exit code.


Your plan only supports up to 10 retries

Upgrade
Number of times to try and restart the service if it stopped due to an error.
Max restart retries
10
Config-as-code
Railway Config File
Config as Code is deprecated. Prefer Infrastructure as Code↗. Existing config files keep working until 2026-12-01. Starting 2026-08-28, services that have never used Config as Code cannot opt in. Legacy docs↗

Add File Path
Feature-flags



Delete Service
Deleting this service will permanently delete all its deployments and remove it from this environment. This cannot be undone.


-------------------------------------------------------------------------------------------

Smart-Parking-Management-System: Python-FastAPI
Deployments
Variables
Metrics
Console
Settings
Filter settings
Filter Settings...

/
Source
Source Repo
SamuelDevadass/Smart-Parking-Management-System



Disconnect
Root Directory
Configure where we should look for your code. Docs↗
Root directory
/backend
Branch connected to production
Changes made to this GitHub branch will be automatically pushed to this environment.
FullstackWebApplication_Go

Disconnect
Auto deploy unavailable
Could not load branches.Retry
Wait for CI
Trigger deployments after all GitHub actions have completed successfully.

Networking
Public networking
Access your application over HTTP with the following domains
smart-parking-management-system-production-feeb.up.railway.app



Domain
smart-parking-management-system-production-feeb
.up.railway.app

Edit Port
Update your domain or choose a target port


Cancel

Update

Generate Domain

Custom Domain

TCP Proxy
Private Networking
Communicate with this service from within the Railway network.
smart-parking-management-system.railway.internal
IPv4 & IPv6


Ready to talk privately ·
You can also simply call me
smart-parking-management-system

DNS
smart-parking-management-system
.railway.internal

Endpoint name available!


Cancel

Update
Outbound IPv6
Enable your service to make outbound connections to IPv6 destinations.

Edge
Under Attack Mode
Being DDoSed? One click puts a browser check in front of your domains, shedding bot floods at the edge. Real visitors solve it once and browse normally.
Until turned off

Activate
Takes effect across the globe in ~20 seconds.

CDN Caching
Cache static assets and optionally HTML at the edge. Reduces latency and origin load for your service's domains. Changes are staged and take effect on deploy without restarting the service.

Edge Rules
Block, redirect, challenge, or override caching for matching requests before they reach your service. Rules run in order; drag to reorder them in the editor.
Edge rules are not available on your plan. Please upgrade to use them.

Scale
Regions & Replicas
Deploy replicas per region for horizontal scaling.
US East (Virginia, USA)

Replicas
1
Replica
Multi-region replicas are only available on the Pro plan.

Learn More↗
Replica Limits
Allocate a maximum vCPU and Memory for each replica.
CPU: 2 vCPU

Plan limit: 2 vCPU


Memory: 1 GB

Plan limit: 1 GB


Upgrade for higher limits
Build
Builder

Railpack

Default

3.13.16python@3.13.16
App builder developed by Railway. Docs↗

Custom Build Command
Override the default build command that is run when building your app. Docs↗
Build command
pip install -r requirements.txt && pip uninstall -y opencv-python opencv-contrib-python opencv-python-headless && pip install --no-deps opencv-python-headless==5.0.0.93 && python -c "import cv2; print(cv2.__version__)"
Watch Paths
Gitignore-style rules to trigger a new deployment based on what file paths have changed. Docs↗
Add pattern
Add pattern e.g. /backend/src/**

Deploy
Custom Start Command
Command that will be run to start new deployments. Docs↗
Start command
uvicorn main:app --host 0.0.0.0 --port $PORT
Add pre-deploy step (Docs↗)
Teardown
Configure old deployment termination when a new one is started. Docs↗

Cron Schedule
Run the service according to the specified cron schedule.

Add Schedule
Healthcheck Path
Endpoint to be called before a deploy completes to ensure the new deployment is live. Docs↗

Healthcheck Path
Serverless
Containers will scale down to zero and then scale up based on traffic. Requests while the container is sleeping will be queued and served when the container wakes up. Docs↗

Restart Policy
Configure what to do when the process exits. Docs↗
On Failure

Restart the container if it exits with a non-zero exit code.


Your plan only supports up to 10 retries

Upgrade
Number of times to try and restart the service if it stopped due to an error.
Max restart retries
10
Config-as-code
Railway Config File
Config as Code is deprecated. Prefer Infrastructure as Code↗. Existing config files keep working until 2026-12-01. Starting 2026-08-28, services that have never used Config as Code cannot opt in. Legacy docs↗

Add File Path
Feature-flags



Delete Service
Deleting this service will permanently delete all its deployments and remove it from this environment. This cannot be undone.

Delete service


