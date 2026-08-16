export interface User {
  id: string;
  nama: string;
  username: string;
  role: 'admin' | 'supervisor' | 'anggota';
  created_at: string;
  updated_at: string;
}

export interface Kendaraan {
  id: string;
  kategori: 'Dinas' | 'Umum';
  jenis: 'Bensin M' | 'Bensin N&O' | 'Solar JBB' | 'Solar GVW';
  plat_nomor: string;
  merek: string;
  tipe: string;
  kapasitas_mesin: number;
  tahun_pembuatan: number;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface HasilUji {
  id: string;
  kendaraan_id: string;
  co?: number;
  co2?: number;
  hc?: number;
  o2?: number;
  lambda?: number;
  opasitas?: number;
  hasil_uji: 'Lulus' | 'Tidak Lulus';
  valid: boolean;
  catatan?: string;
  tested_by: string;
  tested_at: string;
  created_at: string;
  updated_at: string;
}

export interface Parameter {
  id: string;
  kategori: string;
  tahun_min: number;
  tahun_max?: number;
  tahun_operator: '<' | '>' | 'between';
  co_max?: number;
  hc_max?: number;
  opasitas_max?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface PopupNote {
  id: string;
  jenis_kendaraan: string;
  tahun_operator: '>' | '<' | '>=' | '<=' | '=';
  tahun_value: number;
  parameter_uji: 'CO' | 'CO2' | 'HC' | 'O2' | 'Lambda' | 'Opasitas';
  nilai_operator: '>' | '<' | '>=' | '<=' | '=';
  nilai_value: number;
  note: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface HomeContent {
  id: string;
  content_type: 'description' | 'image' | 'working_instruction';
  title?: string;
  description?: string;
  file_path?: string;
  file_type?: string;
  sort_order: number;
  is_active: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface Pagination {
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data: T;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface Statistics {
  total_kendaraan: number;
  total_dinas: number;
  total_umum: number;
  total_lulus: number;
  total_tidak_lulus: number;
  total_bensin: number;
  total_solar: number;
}
