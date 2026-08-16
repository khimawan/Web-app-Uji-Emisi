import React, { useState, useEffect } from 'react';
import { hasilUjiAPI } from '../services/api';
import { Statistics, Pagination } from '../types';
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import toast from 'react-hot-toast';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const HasilUjiPage: React.FC = () => {
  const [statistics, setStatistics] = useState<Statistics | null>(null);
  const [data, setData] = useState<any[]>([]);
  const [pagination, setPagination] = useState<Pagination>({ total: 0, page: 1, limit: 10, totalPages: 0 });
  const [isLoading, setIsLoading] = useState(true);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterHasil, setFilterHasil] = useState('');
  const [filterKategori, setFilterKategori] = useState('');

  useEffect(() => {
    fetchData();
    fetchStatistics();
  }, [pagination.page, searchTerm, filterHasil, filterKategori]);

  const fetchData = async () => {
    try {
      const response = await hasilUjiAPI.getAll({
        page: pagination.page,
        limit: pagination.limit,
        search: searchTerm,
        hasil_uji: filterHasil,
        kategori: filterKategori,
      });
      if (response.success) {
        setData(response.data.items);
        setPagination(response.data.pagination);
      }
    } catch (error) {
      toast.error('Failed to fetch data');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchStatistics = async () => {
    try {
      const response = await hasilUjiAPI.getStatistics();
      if (response.success) {
        setStatistics(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch statistics');
    }
  };

  const handleExportCSV = async () => {
    try {
      const response = await hasilUjiAPI.exportCSV();
      const url = window.URL.createObjectURL(new Blob([response]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'hasil_uji_emisi.csv');
      document.body.appendChild(link);
      link.click();
      link.remove();
      toast.success('CSV downloaded');
    } catch (error) {
      toast.error('Failed to export CSV');
    }
  };

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setSelectedIds(data.map((d) => d.id));
    } else {
      setSelectedIds([]);
    }
  };

  const handleSelectOne = (id: string) => {
    if (selectedIds.includes(id)) {
      setSelectedIds(selectedIds.filter((i) => i !== id));
    } else {
      setSelectedIds([...selectedIds, id]);
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedIds.length === 0) {
      toast.error('No items selected');
      return;
    }
    if (window.confirm(`Are you sure you want to delete ${selectedIds.length} items?`)) {
      try {
        await Promise.all(selectedIds.map((id) => hasilUjiAPI.delete(id)));
        toast.success('Selected items deleted successfully');
        setSelectedIds([]);
        fetchData();
        fetchStatistics();
      } catch (error) {
        toast.error('Failed to delete selected items');
      }
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  // Prepare chart data
  const kategoriData = statistics
    ? [
        { name: 'Dinas', value: statistics.total_dinas },
        { name: 'Umum', value: statistics.total_umum },
      ]
    : [];

  const hasilData = statistics
    ? [
        { name: 'Lulus', value: statistics.total_lulus },
        { name: 'Tidak Lulus', value: statistics.total_tidak_lulus },
      ]
    : [];

  const jenisData = statistics
    ? [
        { name: 'Bensin', value: statistics.total_bensin },
        { name: 'Solar', value: statistics.total_solar },
      ]
    : [];

  const barData = statistics
    ? [
        { name: 'Dinas', jumlah: statistics.total_dinas },
        { name: 'Umum', jumlah: statistics.total_umum },
        { name: 'Bensin', jumlah: statistics.total_bensin },
        { name: 'Solar', jumlah: statistics.total_solar },
      ]
    : [];

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Data Hasil Uji</h1>

      {/* Statistics Cards */}
      {statistics && (
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-3xl font-bold text-blue-600">{statistics.total_kendaraan}</div>
            <div className="text-sm text-gray-500">Total Kendaraan</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-3xl font-bold text-indigo-600">{statistics.total_dinas}</div>
            <div className="text-sm text-gray-500">Kendaraan Dinas</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-3xl font-bold text-green-600">{statistics.total_umum}</div>
            <div className="text-sm text-gray-500">Kendaraan Umum</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-3xl font-bold text-emerald-600">{statistics.total_lulus}</div>
            <div className="text-sm text-gray-500">Lulus Uji</div>
          </div>
          <div className="bg-white rounded-lg shadow p-4 text-center">
            <div className="text-3xl font-bold text-red-600">{statistics.total_tidak_lulus}</div>
            <div className="text-sm text-gray-500">Tidak Lulus</div>
          </div>
        </div>
      )}

      {/* Charts */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-lg font-semibold mb-4 text-center">Dinas vs Umum</h3>
          <ResponsiveContainer width="100%" height={200}>
            <PieChart>
              <Pie
                data={kategoriData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {kategoriData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-lg font-semibold mb-4 text-center">Lulus vs Tidak Lulus</h3>
          <ResponsiveContainer width="100%" height={200}>
            <PieChart>
              <Pie
                data={hasilData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {hasilData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-lg font-semibold mb-4 text-center">Bensin vs Solar</h3>
          <ResponsiveContainer width="100%" height={200}>
            <PieChart>
              <Pie
                data={jenisData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {jenisData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Bar Chart */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-lg font-semibold mb-4">Diagram Batang Jumlah Kendaraan</h3>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={barData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Bar dataKey="jumlah" fill="#8884d8" />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <input
            type="text"
            placeholder="Search by plat nomor, merek..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="border rounded-md px-3 py-2"
          />
          <select
            value={filterHasil}
            onChange={(e) => setFilterHasil(e.target.value)}
            className="border rounded-md px-3 py-2"
          >
            <option value="">Semua Hasil</option>
            <option value="Lulus">Lulus</option>
            <option value="Tidak Lulus">Tidak Lulus</option>
          </select>
          <select
            value={filterKategori}
            onChange={(e) => setFilterKategori(e.target.value)}
            className="border rounded-md px-3 py-2"
          >
            <option value="">Semua Kategori</option>
            <option value="Dinas">Dinas</option>
            <option value="Umum">Umum</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {selectedIds.length > 0 && (
          <div className="bg-yellow-50 border-b p-3 flex items-center justify-between">
            <span className="text-sm text-yellow-700">{selectedIds.length} item(s) selected</span>
            <button
              onClick={handleDeleteSelected}
              className="px-3 py-1 bg-red-500 text-white rounded-md text-sm hover:bg-red-600"
            >
              Delete Selected
            </button>
          </div>
        )}

        <div className="overflow-x-auto max-h-96">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50 sticky top-0">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  <input
                    type="checkbox"
                    onChange={handleSelectAll}
                    checked={selectedIds.length === data.length && data.length > 0}
                  />
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">No</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Waktu</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Kategori</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Jenis</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Plat Nomor</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Merek</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipe</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Kapasitas</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tahun</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">CO</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">CO2</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">HC</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">O2</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Lambda</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Opasitas</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Hasil</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Valid</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Catatan</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {data.length === 0 ? (
                <tr>
                  <td colSpan={19} className="px-4 py-4 text-center text-gray-500">
                    No data available
                  </td>
                </tr>
              ) : (
                data.map((row, index) => (
                  <tr key={row.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 whitespace-nowrap">
                      <input
                        type="checkbox"
                        checked={selectedIds.includes(row.id)}
                        onChange={() => handleSelectOne(row.id)}
                      />
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {(pagination.page - 1) * pagination.limit + index + 1}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {new Date(row.tested_at).toLocaleString('id-ID')}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm">
                      <span
                        className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                          row.kategori === 'Dinas'
                            ? 'bg-blue-100 text-blue-800'
                            : 'bg-green-100 text-green-800'
                        }`}
                      >
                        {row.kategori}
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.jenis}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900">
                      {row.plat_nomor}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.merek}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.tipe}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.kapasitas_mesin}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.tahun_pembuatan}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.co ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.co2 ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.hc ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.o2 ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.lambda ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{row.opasitas ?? '-'}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm">
                      <span
                        className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                          row.hasil_uji === 'Lulus'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                        }`}
                      >
                        {row.hasil_uji}
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm">
                      <span
                        className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                          row.valid ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                        }`}
                      >
                        {row.valid ? 'Valid' : 'Tidak'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-500 max-w-xs truncate">{row.catatan ?? '-'}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200">
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

      {/* Export Buttons */}
      <div className="flex justify-end space-x-2">
        <button
          onClick={handleExportCSV}
          className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
        >
          Download CSV
        </button>
        <button
          onClick={() => toast.success('Excel export coming soon')}
          className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
        >
          Download Excel
        </button>
        <button
          onClick={() => toast.success('PDF export coming soon')}
          className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600"
        >
          Download PDF
        </button>
      </div>
    </div>
  );
};

export default HasilUjiPage;
