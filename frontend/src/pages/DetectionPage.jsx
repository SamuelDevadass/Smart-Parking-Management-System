//Request camera in frontend
import { useState, useEffect, useRef } from "react";
import { Api } from "../api/client.js";
import PlateChip from "../components/PlateChip.jsx";

export default function DetectionPage({ data, updateData, goTo }) {
  
  const [running, setRunning] = useState(false);
  const [error, setError] = useState(null);
  const videoRef = useRef(null);
  const [cameraError, setCameraError] = useState(null);
  const [cameraReady, setCameraReady] = useState(false);

  useEffect(() => {
    let stream;
    async function startCamera() {
      try {
        stream = await navigator.mediaDevices.getUserMedia({
          video: {
            facingMode: "environment",
            width: { ideal: 640 },
            height: { ideal: 480 },
          },
          audio: false,
        });
        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          setCameraReady(true);
        }
      } catch (err) {
        setCameraError(`Could not access camera: ${err.message}`);
      }
    }
    startCamera();
    return () => {
      stream?.getTracks().forEach((track) => track.stop());
    };
  }, []);

  // Clear the plate field each time this page is shown, same as
  useEffect(() => 
  {
    updateData({ licensePlate: "" });
  }, []);
  
  async function startDetection() 
  {
    setError(null);
    try 
    {
      await Api.startDetection();
      setRunning(true);
    } 
    catch (err) 
    {
      setError(err.message);
    }
  }

  async function stopDetection() 
  {
    try 
    {
      await Api.stopDetection();
      const interval = setInterval(async () => 
      {
        try 
        {
          const result = await Api.getDetectionStatus();
          if (result.status === "done") 
          {
            updateData({ licensePlate: result.license_plate || "" });
            setRunning(false);
            clearInterval(interval);
          } 
          else if (result.status === "error") 
          {
            setError("Detection failed on the backend.");
            setRunning(false);
            clearInterval(interval);
          }
        } 
        catch (err) 
        {
          setError(err.message);
          setRunning(false);
          clearInterval(interval);
        }
      }, 500);
    } 
    catch (err) 
    {
    setError(err.message);
    setRunning(false);
    }
  }

  return (
    <div className="panel">
      <h1 className="panel__title">Automatic License Plate Detection</h1>
      {error && <div className="error-banner">{error}</div>}
      <div className="status-row">
        <span className={`status-dot ${running ? "status-dot--running" : ""}`} />
        {running ? "Detection running…" : "Idle"}
    </div>
    <div className="field">
      <label>Live Camera Feed</label>
      <img src="/api/video"  alt={running ? "Live Detection Feed" : "Waiting for detection..."}
        style={{ width: "100%", maxWidth: "640px", border: "1px solid #ccc",
          borderRadius: "8px", minHeight: "360px", objectFit: "contain",}}/>
      {!running && (
        <p className="panel__hint">
          Press <strong>Start Detection</strong> to begin camera streaming.
        </p>
      )}
    </div>
    <div className="btn-row">
      <button className="btn btn--primary" onClick={startDetection} disabled={running}>
        Start Detection </button>
      <button className="btn btn--danger" onClick={stopDetection} disabled={!running}>
        Stop Detection</button>
    </div>
    <div className="field field--mono">
      <label htmlFor="plate-input">License Plate</label>
      <input id="plate-input" value={data.licensePlate}
          onChange={(e) => updateData({ licensePlate: e.target.value })}
          placeholder="Detected plate appears here — or type it manually"/>
    </div>
    <PlateChip value={data.licensePlate} />
      <div className="btn-row btn--spaced">
        <button className="btn btn--primary" disabled={!data.licensePlate} onClick={() => goTo("owner")}> 
          Continue </button>
        <button className="btn" disabled={!data.wing} onClick={() => goTo("spots")}> 
          Empty Spots</button>
      </div>
    </div>
  );
}
