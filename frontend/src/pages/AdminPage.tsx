import React, { useState, useEffect } from 'react';
import { usersAPI } from '../services/api';
import { User, Pagination } from '../types';
import { useAuth } from '../context/AuthContext';
import toast from 'react-hot-toast';

const AdminPage: React.FC = () => {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [pagination, setPagination] = useState<Pagination>({ total: 0, page: 1, limit: 10, totalPages: 0 });
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    nama: '',
    username: '',
    password: '',
    role: 'anggota' as 'admin' | 'supervisor' | 'anggota',
  });

  // Username admin disembunyikan; username supervisor & anggota tetap terlihat.
  const sembunyikanUsername = (role: User['role']) => role === 'admin';

  // Akun admin hanya boleh mengubah username & password-nya sendiri. Akun lain
  // (admin lain atau supervisor) tidak boleh menyentuh kredensial admin.
  const targetAdalahDiriSendiri = Boolean(editingId) && editingId === currentUser?.id;
  const akunAdminMilikOrangLain =
    Boolean(editingId) && formData.role === 'admin' && !targetAdalahDiriSendiri;
  const bolehUbahUsername = !akunAdminMilikOrangLain;
  const bolehUbahPassword = !akunAdminMilikOrangLain;

  useEffect(() => {
    fetchUsers();
  }, [pagination.page]);

  const fetchUsers = async () => {
    setIsLoading(true);
    try {
      const response = await usersAPI.getAll({ page: pagination.page, limit: pagination.limit });
      if (response.success) {
        setUsers(response.data.items);
        setPagination(response.data.pagination);
      }
    } catch (error) {
      toast.error('Failed to fetch users');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingId) {
        // Username & password hanya dikirim bila diisi DAN pemohon berhak
        // mengubahnya, sehingga akun admin milik orang lain tetap utuh.
        const updateData: { nama: string; role: string; username?: string; password?: string } = {
          nama: formData.nama,
          role: formData.role,
        };
        if (bolehUbahUsername && formData.username) {
          updateData.username = formData.username;
        }
        if (bolehUbahPassword && formData.password) {
          updateData.password = formData.password;
        }
        await usersAPI.update(editingId, updateData);
        toast.success('User updated successfully');
      } else {
        await usersAPI.create(formData);
        toast.success('User added successfully');
      }
      setShowForm(false);
      setEditingId(null);
      resetForm();
      fetchUsers();
    } catch (error) {
      toast.error('Failed to save user');
    }
  };

  const handleEdit = (user: User) => {
    setFormData({
      nama: user.nama,
      // Username admin tidak dimuat ke form sehingga tidak pernah tampil di layar.
      username: user.role === 'admin' ? '' : user.username,
      password: '',
      role: user.role,
    });
    setEditingId(user.id);
    setShowForm(true);
  };

  const handleDelete = async (id: string) => {
    const userToDelete = users.find(u => u.id === id);
    if (userToDelete?.role === 'admin') {
      toast.error('Cannot delete admin user');
      return;
    }
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await usersAPI.delete(id);
        toast.success('User deleted successfully');
        fetchUsers();
      } catch (error) {
        toast.error('Failed to delete user');
      }
    }
  };

  const resetForm = () => {
    setFormData({
      nama: '',
      username: '',
      password: '',
      role: 'anggota',
    });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Manajemen Users</h1>
        <button
          onClick={() => { resetForm(); setEditingId(null); setShowForm(true); }}
          className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
        >
          + Tambah User
        </button>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">No</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">              {users.map((user, index) => (
                <tr key={user.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {(pagination.page - 1) * pagination.limit + index + 1}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {user.nama}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {sembunyikanUsername(user.role) ? (
                      <span className="select-none tracking-widest" aria-label="Disembunyikan">••••••••</span>
                    ) : (
                      user.username
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        user.role === 'admin'
                          ? 'bg-red-100 text-red-800'
                          : user.role === 'supervisor'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-green-100 text-green-800'
                      }`}
                    >
                      {user.role}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEdit(user)}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      Edit
                      {user.role === 'admin' && user.id !== currentUser?.id && (
                        <svg
                          className="w-3.5 h-3.5 ml-1 inline-block"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <title>Kredensial hanya dapat diubah oleh admin ini sendiri</title>
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                          />
                        </svg>
                      )}
                    </button>
                    {user.role !== 'admin' && (
                      <button
                        onClick={() => handleDelete(user.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Delete
                      </button>
                    )}
                  </td>
                </tr>
              ))}
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

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {editingId ? 'Edit User' : 'Tambah User'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Nama</label>
                <input
                  type="text"
                  value={formData.nama}
                  onChange={(e) => setFormData({ ...formData, nama: e.target.value })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                  required
                />
              </div>

              <div>
                <label htmlFor="user-username" className="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  id="user-username"
                  name="username"
                  type={formData.role === 'admin' && editingId ? 'password' : 'text'}
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  disabled={!bolehUbahUsername}
                  className="mt-1 block w-full border rounded-md px-3 py-2 disabled:bg-gray-100 disabled:text-gray-500 disabled:cursor-not-allowed"
                  autoComplete="off"
                  placeholder={!bolehUbahUsername ? '••••••••' : undefined}
                  required={!editingId || bolehUbahUsername}
                />
                {editingId && formData.role === 'admin' && (
                  <p className="mt-1 text-xs text-gray-500">
                    {bolehUbahUsername
                      ? 'Ini akun Anda sendiri. Isi dengan username baru bila ingin menggantinya.'
                      : 'Username akun admin hanya dapat diubah oleh admin itu sendiri.'}
                  </p>
                )}
              </div>

              <div>
                <label htmlFor="user-password" className="block text-sm font-medium text-gray-700">
                  Password {editingId && '(Kosongkan jika tidak ingin mengubah)'}
                </label>
                <input
                  id="user-password"
                  name="password"
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  disabled={!bolehUbahPassword}
                  className="mt-1 block w-full border rounded-md px-3 py-2 disabled:bg-gray-100 disabled:text-gray-500 disabled:cursor-not-allowed"
                  autoComplete="new-password"
                  required={!editingId}
                />
                {editingId && formData.role === 'admin' && (
                  <p className="mt-1 text-xs text-gray-500">
                    {bolehUbahPassword
                      ? 'Isi bila ingin mengganti password akun Anda sendiri.'
                      : 'Password akun admin hanya dapat diubah oleh admin itu sendiri.'}
                  </p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Role</label>
                <select
                  value={formData.role}
                  onChange={(e) => setFormData({ ...formData, role: e.target.value as any })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="anggota">Anggota</option>
                  <option value="supervisor">Supervisor</option>
                  <option value="admin">Admin</option>
                </select>
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
    </div>
  );
};

export default AdminPage;
