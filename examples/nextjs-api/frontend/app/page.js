"use client";

import { useState } from "react";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export default function Page() {
  const [payload, setPayload] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function callGateway() {
    setLoading(true);
    setError("");
    try {
      const res = await fetch(`${API_BASE}/api/catalog`);
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      const data = await res.json();
      setPayload({
        ...data,
        poweredBy: res.headers.get("x-powered-by"),
      });
    } catch (err) {
      setError(err.message || "request failed");
      setPayload(null);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main
      style={{
        maxWidth: 720,
        margin: "0 auto",
        padding: "4rem 1.5rem",
      }}
    >
      <p style={{ letterSpacing: "0.08em", textTransform: "uppercase", opacity: 0.7 }}>
        DUG example
      </p>
      <h1 style={{ fontSize: "2.5rem", margin: "0.4rem 0 1rem" }}>Next.js → API Gateway</h1>
      <p style={{ lineHeight: 1.6, opacity: 0.9 }}>
        This page calls <code>{API_BASE}/api/catalog</code>. The browser never talks to
        backends directly — DUG load-balances <code>/api/</code> to upstream services.
      </p>

      <button
        onClick={callGateway}
        disabled={loading}
        style={{
          marginTop: "1.5rem",
          padding: "0.85rem 1.4rem",
          border: 0,
          borderRadius: 8,
          background: "#14b8a6",
          color: "#042f2e",
          fontWeight: 700,
          cursor: "pointer",
        }}
      >
        {loading ? "Calling DUG…" : "Fetch via DUG"}
      </button>

      {error ? (
        <pre style={panelStyle}>{error}</pre>
      ) : null}

      {payload ? (
        <pre style={panelStyle}>{JSON.stringify(payload, null, 2)}</pre>
      ) : null}

      <pre style={{ ...panelStyle, opacity: 0.85 }}>
{`Browser (Next.js :3000)
        │
        │  GET /api/catalog
        ▼
     DUG :8080  ── /api/ ──▶  api1 / api2`}
      </pre>
    </main>
  );
}

const panelStyle = {
  marginTop: "1.5rem",
  padding: "1rem 1.2rem",
  background: "rgba(15, 23, 42, 0.75)",
  border: "1px solid rgba(148, 163, 184, 0.25)",
  borderRadius: 12,
  overflow: "auto",
};
