
import { useState, useEffect, useRef } from "react";
import { Api } from "../api/client.js";
import PlateChip from "../components/PlateChip.jsx";

export default function DetectionPage({ data, updateData, goTo }) {
  const [running, setRunning] = useState(false);
  const [error, setError] = useState(null);
  const [cameraError, setCameraError] = useState(null);
  const [cameraReady, setCameraReady] = useState(false);

  const videoRef = useRef(null);
  const canvasRef = useRef(null);
  const runningRef = useRef(false);
  const loopPromiseRef = useRef(null);
  const stoppingRef = useRef(false);

  // Request access to the frontend camera.
  useEffect(() => {
    let stream;
    let cancelled = false;

    async function startCamera() {
      try {
        if (!navigator.mediaDevices?.getUserMedia) {
          throw new Error("Camera access is unavailable in this browser or context.");
        }

        stream = await navigator.mediaDevices.getUserMedia({
          video: {
            facingMode: "environment",
            width: { ideal: 640 },
            height: { ideal: 480 },
          },
          audio: false,
        });

        if (cancelled) {
          stream.getTracks().forEach((track) => track.stop());
          return;
        }

        const video = videoRef.current;
        if (!video) {
          throw new Error("Video element not found.");
        }

        video.srcObject = stream;

        await new Promise((resolve, reject) => {
          video.onloadedmetadata = resolve;
          video.onerror = () => reject(new Error("Video failed to load."));
        });

        if (cancelled) return;

        await video.play();

        if (!cancelled) {
          setCameraReady(true);
          setCameraError(null);
        }
      } catch (err) {
        if (!cancelled) {
          setCameraError(`Could not access camera: ${err.message}`);
        }
      }
    }

    startCamera();

    return () => {
      cancelled = true;
      runningRef.current = false;

      stream?.getTracks().forEach((track) => track.stop());

      const video = videoRef.current;
      if (video) {
        video.srcObject = null;
      }
    };
  }, []);

  // Clear the plate field when this page mounts.
  useEffect(() => {
    updateData({ licensePlate: "" });
  }, []);

  // Encode the current camera frame as a JPEG Blob.
  async function encodeImage() {
    const video = videoRef.current;
    const canvas = canvasRef.current;

    if (!video || !canvas || video.readyState < 2) {
      throw new Error("Camera frame is not ready.");
    }

    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;

    const context = canvas.getContext("2d");

    if (!context) {
      throw new Error("Could not access the canvas context.");
    }

    context.drawImage(video, 0, 0, canvas.width, canvas.height);

    return new Promise((resolve, reject) => {
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error("Failed to encode camera frame."));
          }
        },
        "image/jpeg",
        0.85
      );
    });
  }

  function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  // Start one backend session, then send frames to the frame endpoint.
  async function startDetection() {
    if (
      runningRef.current ||
      loopPromiseRef.current ||
      stoppingRef.current
    ) {
      return;
    }

    setError(null);
    setCameraError(null);

    try {
      const started = await Api.startDetection();

      if (!started) {
        throw new Error("Could not start the detection session.");
      }

      runningRef.current = true;
      setRunning(true);

      const loop = (async () => {
        while (runningRef.current) {
          const imageBlob = await encodeImage();

          // Send this frame to the backend for YOLO processing.
          const ok = await Api.processDetectionFrame(imageBlob);

          if (!ok) {
            throw new Error("Backend did not accept the camera frame.");
          }

          await sleep(350);
        }
      })();

      loopPromiseRef.current = loop;

      await loop;
    } catch (err) {
      setError(err.message || "Failed to run detection.");
    } finally {
      runningRef.current = false;
      loopPromiseRef.current = null;
      setRunning(false);
    }
  }

  // Stop frame uploads, request final OCR, and retrieve the result.
  async function stopDetection() {
    if (stoppingRef.current) return;

    if (!runningRef.current && !loopPromiseRef.current) return;

    stoppingRef.current = true;
    runningRef.current = false;
    setError(null);

    try {
      // Wait for any frame request already in progress to finish.
      if (loopPromiseRef.current) {
        try {
          await loopPromiseRef.current;
        } catch {
          // A frame error should not prevent attempting to stop the session.
        }
      }

      // The backend performs final capture and OCR here.
      await Api.stopDetection();

      // Poll until the backend finishes OCR.
      while (true) {
        const result = await Api.getDetectionStatus();

        if (result.status === "done") {
          updateData({
            licensePlate: result.license_plate || "",
          });
          break;
        }

        if (result.status === "error") {
          throw new Error("Detection failed on the backend.");
        }

        if (result.status === "idle") {
          throw new Error("The detection session is no longer active.");
        }

        await sleep(500);
      }
    } catch (err) {
      setError(err.message || "Failed to stop detection.");
    } finally {
      runningRef.current = false;
      loopPromiseRef.current = null;
      stoppingRef.current = false;
      setRunning(false);
    }
  }

  return (
    <div className="panel">
      <h1 className="panel__title">
        Automatic License Plate Detection
      </h1>

      {error && <div className="error-banner">{error}</div>}

      {cameraError && (
        <div className="error-banner">{cameraError}</div>
      )}

      <div className="status-row">
        <span
          className={`status-dot ${
            running ? "status-dot--running" : ""
          }`}
        />
        {running ? "Detection running…" : "Idle"}
      </div>

      <div className="field">
        <video
          ref={videoRef}
          autoPlay
          muted
          playsInline
          style={{ display: "none" }}
        />

        <canvas ref={canvasRef} hidden />

        <label>Live Camera Feed</label>

        <img
          src="/api/video"
          alt={
            running
              ? "Live Detection Feed"
              : "Waiting for detection..."
          }
          style={{
            width: "100%",
            maxWidth: "640px",
            border: "1px solid #ccc",
            borderRadius: "8px",
            minHeight: "360px",
            objectFit: "contain",
          }}
        />

        {!running && (
          <p className="panel__hint">
            Press <strong>Start Detection</strong> to begin camera
            streaming.
          </p>
        )}
      </div>

      <div className="btn-row">
        <button
          className="btn btn--primary"
          onClick={startDetection}
          disabled={!cameraReady || running || stoppingRef.current}
        >
          Start Detection
        </button>

        <button
          className="btn btn--danger"
          onClick={stopDetection}
          disabled={!running || stoppingRef.current}
        >
          Stop Detection
        </button>
      </div>

      <div className="field field--mono">
        <label htmlFor="plate-input">License Plate</label>

        <input
          id="plate-input"
          value={data.licensePlate}
          onChange={(e) =>
            updateData({ licensePlate: e.target.value })
          }
          placeholder="Detected plate appears here — or type it manually"
        />
      </div>

      <PlateChip value={data.licensePlate} />

      <div className="btn-row btn--spaced">
        <button
          className="btn btn--primary"
          disabled={!data.licensePlate}
          onClick={() => goTo("owner")}
        >
          Continue
        </button>

        <button
          className="btn"
          disabled={!data.wing}
          onClick={() => goTo("spots")}
        >
          Empty Spots
        </button>
      </div>
    </div>
  );
}
