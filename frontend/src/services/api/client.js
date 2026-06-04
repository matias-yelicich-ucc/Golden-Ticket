import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const loginUser = (payload) => api.post('/login', payload);
export const registerUser = (payload) => api.post('/register', payload);
export const createEvent = (payload) => api.post('/admin/events', payload);
export const getEvents = () => api.get('/events');
