import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URI,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Log requests for debugging purposes
axiosInstance.interceptors.request.use(
  (config) => {
    console.log('Request:', {
      url: config.url,
      method: config.method,
      data: config.data,
    });
    return config;
  },
  (error) => Promise.reject(error)
);

// Handle response errors with improved logging
axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ERR_NETWORK') {
      console.error('Network Error: Ensure the backend is running and CORS is configured correctly.');
      console.error('Backend URL:', import.meta.env.VITE_BACKEND_URI);
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
