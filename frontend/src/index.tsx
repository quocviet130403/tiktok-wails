import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App"
import "./app/globals.css" // Import global CSS for Tailwind

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement)
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
