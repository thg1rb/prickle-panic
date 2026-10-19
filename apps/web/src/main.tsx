import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import { th } from "./i18n/th";
import "./styles.css";

document.title = th.pageTitle;
const description = document.querySelector<HTMLMetaElement>(
  'meta[name="description"]',
);
if (description) description.content = th.pageDescription;

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error("Root element was not found.");
}

createRoot(rootElement).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
