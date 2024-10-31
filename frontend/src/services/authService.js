// src/services/authService.js
import axios from 'axios';

// Set up the base URL for Axios
const API_URL = import.meta.env.VITE_BACKEND_URI + '/api/auth/v1';

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

export const login = async (Username, Password) => {
  try {
    const response = await axios.post(`${API_URL}/login`, { Username, Password });
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

export const logout = async () => {
  try {
    const response = await axios.post(`${API_URL}/logout`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};
