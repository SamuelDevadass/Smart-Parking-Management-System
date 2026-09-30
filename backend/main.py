"""Run with:  uvicorn main:app --reload --port 8000 (can change port number)"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from routers import wings, spots, detection, vehicles_entry_exit, bills, video
from dotenv import load_dotenv
import os, uvicorn
load_dotenv("./.env")

app = FastAPI(title="Smart Parking Management System API")

# The Vite dev server runs on :5173 and is a different origin from :8000,
# so the browser blocks the frontend's fetch() calls unless we explicitly
# allow it here. In production, tighten allow_origins to your real domain.
app.add_middleware(CORSMiddleware,
                    allow_origins=[f"{os.getenv("FRONTEND_URL")}"],
                    allow_methods=["*"],
                    allow_headers=["*"],)

# ---------------------------------------------------------------------------
# Wings
# ---------------------------------------------------------------------------
app.include_router(wings.router)

# ---------------------------------------------------------------------------
# Spot availability
# ---------------------------------------------------------------------------
app.include_router(spots.router)

# ---------------------------------------------------------------------------
# Detection
# ---------------------------------------------------------------------------
app.include_router(detection.router)

# ---------------------------------------------------------------------------
# Vehicles / owners & Entry / exit
# ---------------------------------------------------------------------------
app.include_router(vehicles_entry_exit.router)

# ---------------------------------------------------------------------------
# Billing
# ---------------------------------------------------------------------------
app.include_router(bills.router)

# ---------------------------------------------------------------------------
# Video Feed
# ---------------------------------------------------------------------------
app.include_router(video.router)

"""To run FastAPI on a custom port"""
# if __name__ == "__main__":
#     # Railway/Render pass a PORT environment variable dynamically
#     port = int(os.getenv("PORT", 8001))
#     uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)