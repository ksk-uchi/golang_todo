import axios from "axios";

export const api = axios.create({
  baseURL: "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
    "X-CSRF-Token":
      typeof document !== "undefined"
        ? document.cookie
            .split(";")
            .map((c) => c.trim())
            .find((cookie) => cookie.startsWith("csrf_token="))
            ?.split("=")[1]
        : "",
  },
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);
