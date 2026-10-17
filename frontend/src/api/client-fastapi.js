/* API CLIENT
  ==========
  Every page imports { Api } from here instead
  of touching fetch() directly

  Requests go to relative paths like "/api/wings". In dev, vite.config.js
  proxies "/api/*" to http://localhost:8000 (FastAPI app), so there's
  no CORS to configure locally. In production you'd either serve both from
  the same origin, or set an env var here for an absolute backend URL.
*/


async function request(path, options = {}) {
  const res = await fetch(`/api${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(
      `${options.method || "GET"} ${path} failed (${res.status}): ${text}`
    );
  }

  if (res.status === 204) return null;
  return res.json();
}


export const Api = 
{
  getWings: () => request("/wings"),

  getCentreForWing: (wing) => request(`/wings/${encodeURIComponent(wing)}/centre`),

  getSpotAvailability: (wing) => request(`/spots/availability?wing=${encodeURIComponent(wing)}`),

 
  startDetection: async () => {
      const data = await request("/detection/start", {
        method: "POST",
      });

      return data.ok;
    },

    processDetectionFrame: async (imageBlob) => {
      const data = await request("/detection/frame", {
        method: "POST",
        headers: { "Content-Type": "image/jpeg" },
        body: imageBlob,
      });

      return data.ok;
    },

    getDetectionStatus: async () => {
      const data = await request("/detection/status");

      return data;
    },

    stopDetection: async () => {
      const data = await request("/detection/stop", {
        method: "POST",
      });

      return data.ok;
    },

  getVehicle: (licensePlate) => request(`/vehicles/${encodeURIComponent(licensePlate)}`),

  saveVehicle: (payload) =>
    request("/vehicles", { method: "POST", body: JSON.stringify(payload) }),

  getAvailableSpots: (wing, centreId, size) =>
    request(
      `/spots?wing=${encodeURIComponent(wing)}&centre_id=${encodeURIComponent(
        centreId
      )}&size=${encodeURIComponent(size)}`
    ),

  markEntry: (payload) => request("/entries", { method: "POST", body: JSON.stringify(payload) }),

  getSpot: (licensePlate) => request(`/vehicles/spot/${encodeURIComponent(licensePlate)}`),

  markExit: (licensePlate) =>
    request("/exits", { method: "POST", body: JSON.stringify({ license_plate: licensePlate }) }),

  getLatestBill: (licensePlate) => request(`/bills/${encodeURIComponent(licensePlate)}/latest`),
};
