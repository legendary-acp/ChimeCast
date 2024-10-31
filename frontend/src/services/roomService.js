// src/services/authService.js
import axios from 'axios';

// Set up the base URL for Axios
const API_URL = import.meta.env.VITE_BACKEND_URI + '/api/room/v1';

export const fetchRooms = async () => {
  try {
    const response = await axios.get(`${API_URL}/`);
    if (response.data == null){
        return []
    }
    console.log(response.data)
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

export const createRoom = async (Name) => {
    try {
      const response = await axios.post(`${API_URL}/`,{
        Name
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  };