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

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      const isLoginRequest = error.config.url.endsWith('/login');
      if (!isLoginRequest) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        if (!window.location.pathname.endsWith('/login')) {
          window.location.href = '/login?expired=true';
        }
      }
    }
    return Promise.reject(error);
  }
);

export const loginUser = (payload) => api.post('/login', payload);
export const registerUser = (payload) => api.post('/register', payload);
export const getAdminDashboardStats = () => api.get('/dashboard-stats');
export const createEvent = (payload) => api.post('/admin/events', payload);
export const getEvents = (params) => api.get('/events', { params });
export const getEventByID = (id) => api.get(`/events/${id}`);
export const buyTickets = (eventID, payload) => api.post(`/events/${eventID}/tickets`, payload);
export const getMyTickets = () => api.get('/my-tickets');
export const transferTicket = (ticketId, payload) => api.post(`/my-tickets/${ticketId}/transfer`, payload);
export const cancelTicket = (ticketId) => api.post(`/my-tickets/${ticketId}/cancel`);
export const updateEvent = (id, payload) => api.put(`/admin/events/${id}`, payload);
export const deleteEvent = (id) => api.delete(`/admin/events/${id}`);



