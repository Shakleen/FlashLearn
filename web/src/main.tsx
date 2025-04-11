import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import HomePage from "./pages/HomePage.tsx";
import "bootstrap/dist/css/bootstrap.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import DeckPage from "./pages/DeckPage.tsx";
import DeckFormPage from "./pages/DeckFormPage.tsx";
import DeleteConfirmationPage from "./pages/DeleteConfirmationPage.tsx";
import { Toaster } from "sonner";
import CardFormPage from "./pages/CardFormPage.tsx";

const router = createBrowserRouter([
  { path: "/", element: <HomePage /> },
  { path: "/deck/:id", element: <DeckPage /> },
  { path: "/deck/form/:id", element: <DeckFormPage /> },
  { path: "/deck/delete/:id", element: <DeleteConfirmationPage /> },
  { path: "/card/form/:id", element: <CardFormPage /> },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Toaster richColors />
    <RouterProvider router={router} />
  </StrictMode>
);
