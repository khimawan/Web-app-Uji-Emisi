import React, { useState, useEffect } from 'react';
import { kendaraanAPI } from '../services/api';
import { Kendaraan, Pagination } from '../types';
import toast from 'react-hot-toast';

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

  const [formData, setFormData] = useState({
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

  const handleDownloadTemplate = async () => {
    try {
      const response = await kendaraanAPI.downloadTemplate();
      const url = window.URL.createObjectURL(new Blob([response]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'template_kendaraan.csv');
      document.body.appendChild(link);
      link.click();
      link.remove();
      toast.success('Template downloaded');
    } catch (error) {
      toast.error('Failed to download template');
    }
  };

  const handleUploadCSV = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const response = await kendaraanAPI.uploadCSV(file);
      if (response.success) {
        toast.success(`Uploaded: ${response.data.success} success, ${response.data.failed} failed`);
        fetchKendaraan();
      }
    } catch (error) {
      toast.error('Failed to upload CSV');
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
        <div className="flex space-x-2">
          <button
            onClick={handleDownloadTemplate}
            className="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
          >
            Download Template
          </button>
          <label className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 cursor-pointer">
            Upload CSV
            <input type="file" accept=".csv" onChange={handleUploadCSV} className="hidden" />
          </label>
          <button
            onClick={() => { resetForm(); setEditingId(null); setShowForm(true); }}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            + Tambah Kendaraan
          </button>
        </div>
      </div>

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
            <option value="Dinas">Dinas</option>
            <option value="Umum">Umum</option>
          </select>
          <select
            value={filterJenis}
            onChange={(e) => setFilterJenis(e.target.value)}
            className="border rounded-md px-3 py-2"
          >
            <option value="">Semua Jenis</option>
            <option value="Bensin M">Bensin M</option>
            <option value="Bensin N&O">Bensin N&O</option>
            <option value="Solar JBB">Solar JBB</option>
            <option value="Solar GVW">Solar GVW</option>
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
                  onChange={(e) => setFormData({ ...formData, kategori: e.target.value })}
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
                  onChange={(e) => setFormData({ ...formData, jenis: e.target.value })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="Bensin M">Bensin M (Angkutan Penumpang)</option>
                  <option value="Bensin N&O">Bensin N&O (Angkutan Barang/Gandeng)</option>
                  <option value="Solar JBB">Solar JBB (≤3,5T)</option>
                  <option value="Solar GVW">Solar GVW (>3,5T)</option>
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
                  value={formData.kapasitas_mesin}
                  onChange={(e) => setFormData({ ...formData, kapasitas_mesin: parseInt(e.target.value) })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Tahun Pembuatan</label>
                <input
                  type="number"
                  value={formData.tahun_pembuatan}
                  onChange={(e) => setFormData({ ...formData, tahun_pembuatan: parseInt(e.target.value) })}
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
              {kendaraan.length === 0 ? (
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
