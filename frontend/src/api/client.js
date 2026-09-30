/*
  API CLIENT
  ==========
  Every page imports { Api } from here instead
  of touching fetch() directly.

  Requests go to relative paths like "/api/wings".
  Vite proxies "/api/*" to the Go Huma backend in development.
*/

async function request(path, options = {}) {
  const res = await fetch(`/api${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
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

export const Api = {
  // ============================================================
  // WINGS
  // ============================================================

  getWings: async () => {
    const data = await request("/wings");
    return data.wings;
  },

  getCentreForWing: async (wing) => {
    const data = await request(
      `/wings/${encodeURIComponent(wing)}/centre`
    );

    console.log("FULL RESPONSE:", data);
    console.log("CENTRE_ID:", data?.centre_id);


    return data.centre_id;
  },

  // ============================================================
  // SPOTS
  // ============================================================

getSpotAvailability: (wing) =>
  request(`/spots/availability?wing=${encodeURIComponent(wing)}`),

  getAvailableSpots: async (wing, centreId, size) => {
    const data = await request(
      `/${encodeURIComponent(wing)}/spots/available/?centre_id=${encodeURIComponent(
        centreId
      )}&size=${encodeURIComponent(size)}`
    );

    return data.available_spots_list;
  },

  // ============================================================
  // DETECTION
  // ============================================================

  startDetection: async () => {
    const data = await request("/detection/start", {
      method: "POST",
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

  // ============================================================
  // VEHICLES
  // ============================================================

  getVehicle: async (licensePlate) => {
      const data = await request(
          `/vehicles/${encodeURIComponent(licensePlate)}`
      );

      return data.vehicle_details;
  },

  saveVehicle: async (payload) => {
    const data = await request("/vehicles", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    return data;
  },

  // ============================================================
  // ENTRY
  // ============================================================

  markEntry: async (payload) => {
    const data = await request("/entries", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    return data;
  },

  // ============================================================
  // VEHICLE SPOT
  // ============================================================

  getSpot: async (licensePlate) => {
    const data = await request(
      `/vehicles/spot/${encodeURIComponent(licensePlate)}`
    );

    return data;
  },

  // ============================================================
  // EXIT
  // ============================================================

  markExit: (licensePlate) =>
  request("/exits", {
    method: "PUT",
    body: JSON.stringify({
      license_plate: licensePlate
    })
  }),

  // ============================================================
  // BILL
  // ============================================================

  getLatestBill: async (licensePlate) => {
    const data = await request(
      `/bills/${encodeURIComponent(licensePlate)}/latest`
    );

    return data;
  },
};
