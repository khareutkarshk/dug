export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body
        style={{
          margin: 0,
          minHeight: "100vh",
          fontFamily: "ui-sans-serif, system-ui, sans-serif",
          background: "linear-gradient(160deg, #0f172a 0%, #1e293b 50%, #0f766e 100%)",
          color: "#e2e8f0",
        }}
      >
        {children}
      </body>
    </html>
  );
}
