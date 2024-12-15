// src/services/authService.js
import axios from "./axiosConfig";

// Set up the base URL for Axios
const API_URL = import.meta.env.VITE_BACKEND_URI + "/api/auth/v1";

export const register = async (Name, Username, Email, Password) => {
  try {
    const response = await axios.post(`${API_URL}/register`, {
      Name,
      Username,
      Email,
      Password,
    });
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

export const validateAuth = async () => {
  try {
    const response = await axios.get(`${API_URL}/validate`, {
      withCredentials: true,
    });
    return response.data;
  } catch (error) {
    // Check if the error is due to unauthorized status (401)
    if (error.response?.status === 401) {
      // If using React Router inside a component
      if (typeof window !== 'undefined') {
        window.location.href = '/auth';
      }
      return
    }
    throw error;
  }
};

export const login = async (Username, Password) => {
  try {
    console.log("Attempting login for:", Username);
    const response = await axios.post(`${API_URL}/login`, {
      Username,
      Password,
    });
    console.log("Login response:", response);
    return response.data;
  } catch (error) {
    console.error("Login error:", error);
    if (error.response) {
      throw error.response.data;
    } else if (error.request) {
      throw new Error(
        "No response received from server. Please check your connection."
      );
    } else {
      throw new Error(error.message);
    }
  }
};

export const logout = async () => {
  try {
    const response = await axios.post(`${API_URL}/logout`);
    // Clear cache (localStorage, sessionStorage, etc.)
    localStorage.clear(); // Clear all localStorage keys
    sessionStorage.clear(); // Clear all sessionStorage keys
    document.cookie = ""; // Clears document cookies (if applicable)
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};
