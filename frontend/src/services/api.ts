import axios from 'axios';
import { User, Kendaraan, HasilUji, Parameter, PopupNote, HomeContent, Statistics, Pagination } from '../types';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (username: string, password: string) => {
    const response = await api.post<{ success: boolean; data: { token: string; user: User } }>('/auth/login', { username, password });
    return response.data;
  },
  logout: async () => {
    const response = await api.post('/auth/logout');
    return response.data;
  },
  getCurrentUser: async () => {
    const response = await api.get<{ success: boolean; data: User }>('/auth/me');
    return response.data;
  },
};

// Kendaraan API
export const kendaraanAPI = {
  getAll: async (params?: { page?: number; limit?: number; search?: string; kategori?: string; jenis?: string }) => {
    const response = await api.get<{ success: boolean; data: { items: Kendaraan[]; pagination: Pagination } }>('/kendaraan', { params });
    return response.data;
  },
  getAllList: async () => {
    const response = await api.get<{ success: boolean; data: Kendaraan[] }>('/kendaraan/all');
    return response.data;
  },
  search: async (query: string) => {
    const response = await api.get<{ success: boolean; data: Kendaraan[] }>('/kendaraan/search', { params: { q: query } });
    return response.data;
  },
  getById: async (id: string) => {
    const response = await api.get<{ success: boolean; data: Kendaraan }>(`/kendaraan/${id}`);
    return response.data;
  },
  create: async (data: Omit<Kendaraan, 'id' | 'created_by' | 'created_at' | 'updated_at'>) => {
    const response = await api.post<{ success: boolean; data: Kendaraan }>('/kendaraan', data);
    return response.data;
  },
  update: async (id: string, data: Omit<Kendaraan, 'id' | 'created_by' | 'created_at' | 'updated_at'>) => {
    const response = await api.put<{ success: boolean; data: Kendaraan }>(`/kendaraan/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/kendaraan/${id}`);
    return response.data;
  },
  uploadCSV: async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post('/kendaraan/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
  downloadTemplate: async () => {
    const response = await api.get('/kendaraan/template/csv', { responseType: 'blob' });
    return response.data;
  },
};

// Hasil Uji API
export const hasilUjiAPI = {
  getAll: async (params?: { page?: number; limit?: number; search?: string; hasil_uji?: string; kategori?: string; tanggal_mulai?: string; tanggal_akhir?: string }) => {
    const response = await api.get<{ success: boolean; data: { items: any[]; pagination: Pagination } }>('/hasil-uji', { params });
    return response.data;
  },
  getStatistics: async () => {
    const response = await api.get<{ success: boolean; data: Statistics }>('/hasil-uji/statistics');
    return response.data;
  },
  create: async (data: { kendaraan_id: string; co?: number; co2?: number; hc?: number; o2?: number; lambda?: number; opasitas?: number }) => {
    const response = await api.post<{ success: boolean; data: { hasil_uji: HasilUji; popup_notes: string[] } }>('/hasil-uji', data);
    return response.data;
  },
  update: async (id: string, data: { kendaraan_id: string; co?: number; co2?: number; hc?: number; o2?: number; lambda?: number; opasitas?: number }) => {
    const response = await api.put<{ success: boolean; data: HasilUji }>(`/hasil-uji/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/hasil-uji/${id}`);
    return response.data;
  },
  exportCSV: async () => {
    const response = await api.get('/hasil-uji/export/csv', { responseType: 'blob' });
    return response.data;
  },
};

// Parameters API
export const parametersAPI = {
  getAll: async (kategori?: string) => {
    const response = await api.get<{ success: boolean; data: Parameter[] }>('/parameters', { params: { kategori } });
    return response.data;
  },
  create: async (data: Omit<Parameter, 'id' | 'is_active' | 'created_at' | 'updated_at'>) => {
    const response = await api.post<{ success: boolean; data: Parameter }>('/parameters', data);
    return response.data;
  },
  update: async (id: string, data: Partial<Parameter>) => {
    const response = await api.put<{ success: boolean; data: Parameter }>(`/parameters/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/parameters/${id}`);
    return response.data;
  },
};

// Popup Notes API
export const popupNotesAPI = {
  getAll: async () => {
    const response = await api.get<{ success: boolean; data: PopupNote[] }>('/popup-notes');
    return response.data;
  },
  create: async (data: Omit<PopupNote, 'id' | 'is_active' | 'created_at' | 'updated_at'>) => {
    const response = await api.post<{ success: boolean; data: PopupNote }>('/popup-notes', data);
    return response.data;
  },
  update: async (id: string, data: Partial<PopupNote>) => {
    const response = await api.put<{ success: boolean; data: PopupNote }>(`/popup-notes/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/popup-notes/${id}`);
    return response.data;
  },
};

// Home Content API
export const homeAPI = {
  getContent: async () => {
    const response = await api.get<{ success: boolean; data: { descriptions: HomeContent[]; images: HomeContent[]; working_instruction: HomeContent } }>('/home');
    return response.data;
  },
  updateDescription: async (id: string, data: { title: string; description: string }) => {
    const response = await api.put(`/home/description/${id}`, data);
    return response.data;
  },
  uploadImage: async (file: File, title: string) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('title', title);
    const response = await api.post('/home/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
  deleteImage: async (id: string) => {
    const response = await api.delete(`/home/image/${id}`);
    return response.data;
  },
  uploadWorkingInstruction: async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post('/home/working-instruction', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
  deleteWorkingInstruction: async (id: string) => {
    const response = await api.delete(`/home/working-instruction/${id}`);
    return response.data;
  },
};

// Users API
export const usersAPI = {
  getAll: async (params?: { page?: number; limit?: number; role?: string }) => {
    const response = await api.get<{ success: boolean; data: { items: User[]; pagination: Pagination } }>('/users', { params });
    return response.data;
  },
  create: async (data: { nama: string; username: string; password: string; role: string }) => {
    const response = await api.post<{ success: boolean; data: User }>('/users', data);
    return response.data;
  },
  update: async (id: string, data: { nama?: string; username?: string; password?: string; role?: string }) => {
    const response = await api.put<{ success: boolean; data: User }>(`/users/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    const response = await api.delete(`/users/${id}`);
    return response.data;
  },
};

export default api;
