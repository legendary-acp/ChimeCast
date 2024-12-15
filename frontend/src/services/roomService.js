import axios from './axiosConfig';

// Set up the base URL for Axios
const API_URL = import.meta.env.VITE_BACKEND_URI + '/api/room/v1';

// Get all rooms
export const fetchRooms = async () => {
  try {
    const response = await axios.get(`${API_URL}/`);
    if (response.data == null) {
      return []
    }
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Create a new room
export const createRoom = async (Name) => {
  try {
    const response = await axios.post(`${API_URL}/`, {
      Name
    });
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Get room status
export const getRoomStatus = async (roomId) => {
  try {
    const response = await axios.get(`${API_URL}/${roomId}/status`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Join a room
export const joinRoom = async (roomId) => {
  try {
    const response = await axios.post(`${API_URL}/${roomId}/join`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Get room participants
export const getRoomParticipants = async (roomId) => {
  try {
    const response = await axios.get(`${API_URL}/${roomId}/participants`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Admit user to room (host only)
export const admitUser = async (roomId, userId) => {
  try {
    const response = await axios.post(`${API_URL}/${roomId}/admit/${userId}`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Deny user from room (host only)
export const denyUser = async (roomId, userId) => {
  try {
    const response = await axios.post(`${API_URL}/${roomId}/deny/${userId}`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Leave room
export const leaveRoom = async (roomId) => {
  try {
    const response = await axios.post(`${API_URL}/${roomId}/leave`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// Set up WebSocket connection
export const getWebSocketUrl = (roomId) => {
  return `${import.meta.env.VITE_WS_URI}/api/room/v1/${roomId}/ws`;
};