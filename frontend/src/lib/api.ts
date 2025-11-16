import axios from "axios";

const getBaseURL = () => {
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }

  if (typeof window !== "undefined") {
    return "/api";
  }

  const appUrl = process.env.APP_URL || "http://localhost:3000";
  return `${appUrl}/api`;
};

export const api = axios.create({
  baseURL: getBaseURL(),
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const apiError = getApiError(error);
    return Promise.reject(apiError);
  },
);

export const getApiError = (error: unknown) => {
  if (axios.isAxiosError(error)) {
    const { response } = error;
    if (response) {
      const apiError = new Error(
        response.data.message || response.statusText || "An API error occurred",
      ) as Error & { status?: number; data?: any };

      apiError.status = response.status;
      apiError.data = response.data;

      return apiError;
    }

    return new Error(error.message || "Network error occurred");
  }

  if (error instanceof Error) {
    return error;
  }

  // Handle string errors
  if (typeof error === "string") {
    return new Error(error);
  }

  return new Error("An unexpected error occurred");
};