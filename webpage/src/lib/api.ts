import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

// Request interceptor - add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle auth errors
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

// Type definitions
export interface User {
  id: number;
  username: string;
  role: string;
}

export interface Container {
  name: string;
  image: string;
  status: string;
  ports: string[];
  id: string;
}

export interface ImageInfo {
  id: string;
  tags: string[];
  size: string;
  created: string;
}

export interface Template {
  id: number;
  os_type: string;
  os_version: string;
  de_name: string;
  display_name: string;
  description: string;
  image_full: string;
  icon: string;
  category: string;
  is_active: boolean;
  created_at: string;
}

export interface EngineStats {
  containers_running: number;
  containers_total: number;
  images_total: number;
  engine_version: string;
}

export interface SystemStatus {
  engine: string;
  version: string;
  containers_count: number;
}
