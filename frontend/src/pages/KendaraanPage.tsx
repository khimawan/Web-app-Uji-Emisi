import React, { useState, useEffect } from 'react';
import { kendaraanAPI } from '../services/api';
import { Kendaraan, Pagination } from '../types';
import toast from 'react-hot-toast';

// Nilai yang sama dengan backend (backend/controllers/kendaraan_csv.go).
const KATEGORI_VALID = ['Dinas', 'Umum'] as const;
const JENIS_VALID = [
  { nilai: 'Bensin M', label: 'Bensin M (Angkutan Penumpang)' },
  { nilai: 'Bensin N&O', label: 'Bensin N&O (Angkutan Barang/Gendang)' },
  { nilai: 'Solar JBB', label: 'Solar JBB (Berat Total ≤ 3,5 ton)' },
  { nilai: 'Solar GVW', label: 'Solar GVW (Berat Total > 3,5 ton)' },
] as const;

type FormKendaraan = {
  kategori: (typeof KATEGORI_VALID)[number];
  jenis: (typeof JENIS_VALID)[number]['nilai'];
  plat_nomor: string;
  merek: string;
  tipe: string;
  kapasitas_mesin: number;
  tahun_pembuatan: number;
};

type ErrorUpload = { row: number; message: string };

const FORMAT_FILE = ['.csv', '.txt', '.xlsx', '.xlsm'];

function unduhBlob(blob: Blob, namaFile: string) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = namaFile;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

const KendaraanPage: React.FC = () => {
  const [kendaraan, setKendaraan] = useState<Kendaraan[]>([]);
  const [pagination, setPagination] = useState<Pagination>({ total: 0, page: 1, limit: 10, totalPages: 0 });
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterKategori, setFilterKategori] = useState('');
  const [filterJenis, setFilterJenis] = useState('');
  const [uploadErrors, setUploadErrors] = useState<ErrorUpload[]>([]);

  const [formData, setFormData] = useState<FormKendaraan>({
    kategori: 'Umum',
    jenis: 'Bensin M',
    plat_nomor: '',
    merek: '',
    tipe: '',
    kapasitas_mesin: 0,
    tahun_pembuatan: 2020,
  });

  useEffect(() => {
    fetchKendaraan();
  }, [pagination.page, searchTerm, filterKategori, filterJenis]);

  const fetchKendaraan = async () => {
    setIsLoading(true);
    try {
      const response = await kendaraanAPI.getAll({
        page: pagination.page,
        limit: pagination.limit,
        search: searchTerm,
        kategori: filterKategori,
        jenis: filterJenis,
      });
      if (response.success) {
        setKendaraan(response.data.items);
        setPagination(response.data.pagination);
      }
    } catch (error) {
      toast.error('Failed to fetch kendaraan');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingId) {
        await kendaraanAPI.update(editingId, formData);
        toast.success('Kendaraan updated successfully');
      } else {
        await kendaraanAPI.create(formData);
        toast.success('Kendaraan added successfully');
      }
      setShowForm(false);
      setEditingId(null);
      resetForm();
      fetchKendaraan();
    } catch (error) {
      toast.error('Failed to save kendaraan');
    }
  };

  const handleEdit = (k: Kendaraan) => {
    setFormData({
      kategori: k.kategori,
      jenis: k.jenis,
      plat_nomor: k.plat_nomor,
      merek: k.merek,
      tipe: k.tipe,
      kapasitas_mesin: k.kapasitas_mesin,
      tahun_pembuatan: k.tahun_pembuatan,
    });
    setEditingId(k.id);
    setShowForm(true);
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this kendaraan?')) {
      try {
        await kendaraanAPI.delete(id);
        toast.success('Kendaraan deleted successfully');
        fetchKendaraan();
      } catch (error) {
        toast.error('Failed to delete kendaraan');
      }
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedIds.length === 0) {
      toast.error('No items selected');
      return;
    }
    if (window.confirm(`Are you sure you want to delete ${selectedIds.length} items?`)) {
      try {
        await Promise.all(selectedIds.map(id => kendaraanAPI.delete(id)));
        toast.success('Selected items deleted successfully');
        setSelectedIds([]);
        fetchKendaraan();
      } catch (error) {
        toast.error('Failed to delete selected items');
      }
    }
  };

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setSelectedIds(kendaraan.map(k => k.id));
    } else {
      setSelectedIds([]);
    }
  };

  const handleSelectOne = (id: string) => {
    if (selectedIds.includes(id)) {
      setSelectedIds(selectedIds.filter(i => i !== id));
    } else {
      setSelectedIds([...selectedIds, id]);
    }
  };

  const handleDownloadTemplate = async (format: 'csv' | 'excel') => {
    try {
      const response = await kendaraanAPI.downloadTemplate(format);
      unduhBlob(new Blob([response]), format === 'excel' ? 'template_kendaraan.xlsx' : 'template_kendaraan.csv');
      toast.success(format === 'excel' ? 'Template Excel berhasil diunduh' : 'Template CSV berhasil diunduh');
    } catch (error) {
      toast.error('Gagal mengunduh template');
    }
  };

  const handleUploadFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;

    const ext = `.${file.name.split('.').pop()?.toLowerCase() ?? ''}`;
    if (!FORMAT_FILE.includes(ext)) {
      setUploadErrors([]);
      toast.error(`Format ${ext || 'file'} belum didukung. Gunakan ${FORMAT_FILE.join(', ')}.`);
      return;
    }

    try {
      const response = await kendaraanAPI.uploadFile(file);
      if (!response.success) {
        setUploadErrors([]);
        toast.error(response.message || 'Upload gagal');
        return;
      }

      const { success = 0, failed = 0, errors = [] } = response.data || {};
      setUploadErrors(errors);

      if (failed > 0) {
        toast.error(`Upload selesai: ${success} tersimpan, ${failed} baris ditolak`);
      } else {
        toast.success(`${success} data kendaraan berhasil diunggah`);
      }
      fetchKendaraan();
    } catch (error) {
      setUploadErrors([]);
      toast.error('Gagal mengunggah file');
    }
  };

  const resetForm = () => {
    setFormData({
      kategori: 'Umum',
      jenis: 'Bensin M',
      plat_nomor: '',
      merek: '',
      tipe: '',
      kapasitas_mesin: 0,
      tahun_pembuatan: 2020,
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Input Data Kendaraan</h1>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => handleDownloadTemplate('csv')}
            className="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
          >
            Template CSV
          </button>
          <button
            onClick={() => handleDownloadTemplate('excel')}
            className="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
          >
            Template Excel
          </button>
          <label className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 cursor-pointer">
            Upload File
            <input type="file" accept={FORMAT_FILE.join(',')} onChange={handleUploadFile} className="hidden" />
          </label>
          <button
            onClick={() => { resetForm(); setEditingId(null); setShowForm(true); }}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            + Tambah Kendaraan
          </button>
        </div>
      </div>

      {uploadErrors.length > 0 && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex justify-between items-start">
            <h2 className="text-sm font-semibold text-red-800">
              {uploadErrors.length} baris ditolak saat upload
            </h2>
            <button onClick={() => setUploadErrors([])} className="text-red-600 hover:text-red-800 text-sm">
              Tutup
            </button>
          </div>
          <ul className="mt-2 space-y-1 text-sm text-red-700 max-h-60 overflow-y-auto">
            {uploadErrors.map((err, i) => (
              <li key={`${err.row}-${i}`}>
                <span className="font-medium">Baris {err.row}:</span> {err.message}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Panduan format upload */}
      <details className="bg-white rounded-lg shadow p-4">
        <summary className="cursor-pointer text-sm font-semibold text-gray-700">
          Panduan format upload CSV &amp; Excel (semua jenis input yang tersedia)
        </summary>
        <div className="mt-3 text-sm text-gray-600">
          <ol className="list-decimal list-inside space-y-1 mb-4">
            <li>Klik <span className="font-medium">Template CSV</span> atau <span className="font-medium">Template Excel</span> — keduanya memuat panduan, contoh baris, dan (untuk Excel) daftar nilai yang bisa dipilih lewat dropdown.</li>
            <li>Buka template, ganti/hapus baris contoh, lalu isi data Anda.</li>
            <li>Simpan file (tetap <code className="text-xs bg-gray-100 px-1 rounded">.csv</code> atau <code className="text-xs bg-gray-100 px-1 rounded">.xlsx</code>) lalu unggah lewat tombol <span className="font-medium">Upload File</span>.</li>
            <li>Baris yang formatnya salah ditolak per baris lengkap dengan alasannya; baris lain tetap tersimpan.</li>
          </ol>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="font-medium text-gray-700 mb-1">Kolom file (urutan bebas):</p>
              <ul className="list-disc list-inside space-y-0.5">
                {['kategori', 'jenis', 'plat_nomor', 'merek', 'tipe', 'kapasitas_mesin', 'tahun_pembuatan'].map((k) => (
                  <li key={k}><code className="text-xs bg-gray-100 px-1 rounded">{k}</code></li>
                ))}
              </ul>
            </div>
            <div>
              <p className="font-medium text-gray-700 mb-1">Nilai yang tersedia:</p>
              <ul className="list-disc list-inside space-y-0.5">
                <li>Kategori: {KATEGORI_VALID.join(' / ')}</li>
                {JENIS_VALID.map((j) => (
                  <li key={j.nilai}>Jenis: {j.label}</li>
                ))}
              </ul>
            </div>
            <div className="md:col-span-2">
              <p className="font-medium text-gray-700 mb-1">Aturan penulisan:</p>
              <ul className="list-disc list-inside space-y-0.5">
                <li>Format file: {FORMAT_FILE.join(', ')} — <code className="text-xs bg-gray-100 px-1 rounded">.xls</code> (Excel 97-2003) tidak didukung, simpan ulang sebagai <code className="text-xs bg-gray-100 px-1 rounded">.xlsx</code>.</li>
                <li>kapasitas_mesin: angka bulat cc (contoh <code className="text-xs bg-gray-100 px-1 rounded">1496</code>).</li>
                <li>tahun_pembuatan: angka bulat 4 digit (contoh <code className="text-xs bg-gray-100 px-1 rounded">2020</code>).</li>
                <li>plat_nomor: huruf kapital, maks. 20 karakter (contoh <code className="text-xs bg-gray-100 px-1 rounded">B 1234 CDE</code>).</li>
                <li>Pemisah kolom <code className="text-xs bg-gray-100 px-1 rounded">,</code> / <code className="text-xs bg-gray-100 px-1 rounded">;</code> / tab / <code className="text-xs bg-gray-100 px-1 rounded">|</code> otomatis dideteksi.</li>
                <li>Baris kosong dan baris berawalan <code className="text-xs bg-gray-100 px-1 rounded">#</code> otomatis dilewati.</li>
              </ul>
            </div>
          </div>
        </div>
      </details>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <input
            type="text"
            placeholder="Search by plat nomor, merek, tipe..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="border rounded-md px-3 py-2"
          />
          <select
            value={filterKategori}
            onChange={(e) => setFilterKategori(e.target.value)}
            className="border rounded-md px-3 py-2"
          >
            <option value="">Semua Kategori</option>
            {KATEGORI_VALID.map((k) => (
              <option key={k} value={k}>{k}</option>
            ))}
          </select>
          <select
            value={filterJenis}
            onChange={(e) => setFilterJenis(e.target.value)}
            className="border rounded-md px-3 py-2"
          >
            <option value="">Semua Jenis</option>
            {JENIS_VALID.map((j) => (
              <option key={j.nilai} value={j.nilai}>{j.nilai}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {editingId ? 'Edit Kendaraan' : 'Tambah Kendaraan'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Kategori</label>
                <select
                  value={formData.kategori}
                  onChange={(e) => setFormData({ ...formData, kategori: e.target.value as FormKendaraan['kategori'] })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="Dinas">Dinas</option>
                  <option value="Umum">Umum</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Jenis Kendaraan</label>
                <select
                  value={formData.jenis}
                  onChange={(e) => setFormData({ ...formData, jenis: e.target.value as FormKendaraan['jenis'] })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  {JENIS_VALID.map((j) => (
                    <option key={j.nilai} value={j.nilai}>{j.label}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Plat Nomor</label>
                <input
                  type="text"
                  value={formData.plat_nomor}
                  onChange={(e) => setFormData({ ...formData, plat_nomor: e.target.value })}
                  placeholder="Contoh: B 1234 CDE"
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Merek Mobil</label>
                <input
                  type="text"
                  value={formData.merek}
                  onChange={(e) => setFormData({ ...formData, merek: e.target.value })}
                  placeholder="Contoh: Toyota"
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Tipe Mobil</label>
                <input
                  type="text"
                  value={formData.tipe}
                  onChange={(e) => setFormData({ ...formData, tipe: e.target.value })}
                  placeholder="Contoh: Avanza"
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Kapasitas Mesin (cc)</label>
                <input
                  type="number"
                  min={1}
                  step={1}
                  value={formData.kapasitas_mesin}
                  onChange={(e) => setFormData({ ...formData, kapasitas_mesin: Number(e.target.value) || 0 })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Tahun Pembuatan</label>
                <input
                  type="number"
                  min={1900}
                  max={new Date().getFullYear() + 2}
                  step={1}
                  value={formData.tahun_pembuatan}
                  onChange={(e) => setFormData({ ...formData, tahun_pembuatan: Number(e.target.value) || 0 })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => { setShowForm(false); setEditingId(null); }}
                  className="px-4 py-2 border rounded-md text-gray-700 hover:bg-gray-50"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
                >
                  Simpan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {selectedIds.length > 0 && (
          <div className="bg-yellow-50 border-b p-3 flex items-center justify-between">
            <span className="text-sm text-yellow-700">
              {selectedIds.length} item(s) selected
            </span>
            <button
              onClick={handleDeleteSelected}
              className="px-3 py-1 bg-red-500 text-white rounded-md text-sm hover:bg-red-600"
            >
              Delete Selected
            </button>
          </div>
        )}

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <input
                    type="checkbox"
                    onChange={handleSelectAll}
                    checked={selectedIds.length === kendaraan.length && kendaraan.length > 0}
                  />
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">No</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Plat Nomor</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kategori</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jenis</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Merek</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tipe</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kapasitas</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tahun</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {isLoading ? (
                <tr>
                  <td colSpan={10} className="px-6 py-4 text-center text-gray-500">
                    Memuat data...
                  </td>
                </tr>
              ) : kendaraan.length === 0 ? (
                <tr>
                  <td colSpan={10} className="px-6 py-4 text-center text-gray-500">
                    No data available
                  </td>
                </tr>
              ) : (
                kendaraan.map((k, index) => (
                  <tr key={k.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <input
                        type="checkbox"
                        checked={selectedIds.includes(k.id)}
                        onChange={() => handleSelectOne(k.id)}
                      />
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {(pagination.page - 1) * pagination.limit + index + 1}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {k.plat_nomor}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        k.kategori === 'Dinas' ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'
                      }`}>
                        {k.kategori}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {k.jenis}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {k.merek}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {k.tipe}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {k.kapasitas_mesin} cc
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {k.tahun_pembuatan}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        onClick={() => handleEdit(k)}
                        className="text-blue-600 hover:text-blue-900 mr-3"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(k.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
          <div className="flex-1 flex justify-between sm:hidden">
            <button
              onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
              disabled={pagination.page === 1}
              className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              Previous
            </button>
            <button
              onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
              disabled={pagination.page === pagination.totalPages}
              className="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              Next
            </button>
          </div>
          <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p className="text-sm text-gray-700">
                Showing <span className="font-medium">{(pagination.page - 1) * pagination.limit + 1}</span> to{' '}
                <span className="font-medium">
                  {Math.min(pagination.page * pagination.limit, pagination.total)}
                </span>{' '}
                of <span className="font-medium">{pagination.total}</span> results
              </p>
            </div>
            <div>
              <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                {Array.from({ length: pagination.totalPages }, (_, i) => i + 1).map((page) => (
                  <button
                    key={page}
                    onClick={() => setPagination({ ...pagination, page })}
                    className={`relative inline-flex items-center px-4 py-2 border text-sm font-medium ${
                      pagination.page === page
                        ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                        : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                    }`}
                  >
                    {page}
                  </button>
                ))}
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KendaraanPage;
