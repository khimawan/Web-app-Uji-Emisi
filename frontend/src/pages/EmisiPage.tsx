import React, { useState, useEffect } from 'react';
import { kendaraanAPI, hasilUjiAPI } from '../services/api';
import { Kendaraan } from '../types';
import toast from 'react-hot-toast';

const EmisiPage: React.FC = () => {
  const [kendaraanList, setKendaraanList] = useState<Kendaraan[]>([]);
  const [filteredKendaraan, setFilteredKendaraan] = useState<Kendaraan[]>([]);
  const [selectedKendaraan, setSelectedKendaraan] = useState<Kendaraan | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [showPopup, setShowPopup] = useState(false);
  const [popupData, setPopupData] = useState<{ hasil_uji: string; valid: boolean; popup_notes: string[] } | null>(null);

  const [formData, setFormData] = useState({
    co: '',
    co2: '',
    hc: '',
    o2: '',
    lambda: '',
    opasitas: '',
  });

  useEffect(() => {
    fetchKendaraan();
  }, []);

  useEffect(() => {
    if (searchTerm) {
      const filtered = kendaraanList.filter(
        (k) =>
          k.plat_nomor.toLowerCase().includes(searchTerm.toLowerCase()) ||
          k.merek.toLowerCase().includes(searchTerm.toLowerCase()) ||
          k.tipe.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setFilteredKendaraan(filtered);
    } else {
      setFilteredKendaraan(kendaraanList);
    }
  }, [searchTerm, kendaraanList]);

  const fetchKendaraan = async () => {
    try {
      const response = await kendaraanAPI.getAllList();
      if (response.success) {
        setKendaraanList(response.data);
        setFilteredKendaraan(response.data);
      }
    } catch (error) {
      toast.error('Failed to fetch kendaraan');
    } finally {
      setIsLoading(false);
    }
  };

  const isBensin = selectedKendaraan?.jenis.startsWith('Bensin');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedKendaraan) {
      toast.error('Please select a kendaraan');
      return;
    }

    try {
      const payload: any = {
        kendaraan_id: selectedKendaraan.id,
      };

      if (isBensin) {
        if (formData.co) payload.co = parseFloat(formData.co);
        if (formData.co2) payload.co2 = parseFloat(formData.co2);
        if (formData.hc) payload.hc = parseFloat(formData.hc);
        if (formData.o2) payload.o2 = parseFloat(formData.o2);
        if (formData.lambda) payload.lambda = parseFloat(formData.lambda);
      } else {
        if (formData.opasitas) payload.opasitas = parseFloat(formData.opasitas);
      }

      const response = await hasilUjiAPI.create(payload);
      if (response.success) {
        setPopupData({
          hasil_uji: response.data.hasil_uji.hasil_uji,
          valid: response.data.hasil_uji.valid,
          popup_notes: response.data.popup_notes || [],
        });
        setShowPopup(true);
        resetForm();
        toast.success('Hasil uji saved successfully');
      }
    } catch (error) {
      toast.error('Failed to save hasil uji');
    }
  };

  const resetForm = () => {
    setSelectedKendaraan(null);
    setSearchTerm('');
    setFormData({
      co: '',
      co2: '',
      hc: '',
      o2: '',
      lambda: '',
      opasitas: '',
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
      <h1 className="text-2xl font-bold text-gray-900">Input Data Emisi</h1>

      {/* Search and Select Kendaraan */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-4">Pilih Kendaraan</h2>
        <div className="relative">
          <input
            type="text"
            placeholder="Cari kendaraan berdasarkan plat nomor, merek, atau tipe..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full border rounded-md px-4 py-2 pl-10"
          />
          <svg
            className="absolute left-3 top-2.5 h-5 w-5 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>

        {searchTerm && filteredKendaraan.length > 0 && (
          <div className="mt-2 border rounded-md max-h-60 overflow-y-auto">
            {filteredKendaraan.map((k) => (
              <div
                key={k.id}
                onClick={() => {
                  setSelectedKendaraan(k);
                  setSearchTerm('');
                }}
                className="p-3 hover:bg-gray-50 cursor-pointer border-b last:border-b-0"
              >
                <div className="flex justify-between items-center">
                  <div>
                    <span className="font-medium">{k.plat_nomor}</span>
                    <span className="ml-2 text-gray-500">
                      {k.merek} {k.tipe}
                    </span>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      k.kategori === 'Dinas'
                        ? 'bg-blue-100 text-blue-800'
                        : 'bg-green-100 text-green-800'
                    }`}
                  >
                    {k.kategori}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Selected Kendaraan Info */}
      {selectedKendaraan && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="font-semibold text-blue-900">Kendaraan Dipilih</h3>
              <p className="text-blue-700">
                {selectedKendaraan.plat_nomor} - {selectedKendaraan.merek}{' '}
                {selectedKendaraan.tipe}
              </p>
              <p className="text-sm text-blue-600">
                {selectedKendaraan.kategori} | {selectedKendaraan.jenis} |{' '}
                {selectedKendaraan.tahun_pembuatan}
              </p>
            </div>
            <button
              onClick={() => setSelectedKendaraan(null)}
              className="text-red-500 hover:text-red-700"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* Form Input Emisi */}
      {selectedKendaraan && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">Input Hasil Uji Emisi</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            {isBensin ? (
              <>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">CO (%)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={formData.co}
                      onChange={(e) => setFormData({ ...formData, co: e.target.value })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">CO2 (%)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={formData.co2}
                      onChange={(e) => setFormData({ ...formData, co2: e.target.value })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">HC (ppm)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={formData.hc}
                      onChange={(e) => setFormData({ ...formData, hc: e.target.value })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">O2 (%)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={formData.o2}
                      onChange={(e) => setFormData({ ...formData, o2: e.target.value })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Lambda</label>
                    <input
                      type="number"
                      step="0.0001"
                      value={formData.lambda}
                      onChange={(e) => setFormData({ ...formData, lambda: e.target.value })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                </div>
                <p className="text-sm text-gray-500">
                  * Untuk kendaraan bensin, opasitas tidak perlu diinput
                </p>
              </>
            ) : (
              <>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Opasitas (%)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.opasitas}
                    onChange={(e) => setFormData({ ...formData, opasitas: e.target.value })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  />
                </div>
                <p className="text-sm text-gray-500">
                  * Untuk kendaraan solar, hanya perlu input opasitas
                </p>
              </>
            )}

            <div className="flex justify-end space-x-2">
              <button
                type="button"
                onClick={resetForm}
                className="px-4 py-2 border rounded-md text-gray-700 hover:bg-gray-50"
              >
                Batal
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
              >
                Simpan Hasil Uji
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Popup Hasil Uji */}
      {showPopup && popupData && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <div className="text-center">
              <div
                className={`mx-auto flex items-center justify-center h-16 w-16 rounded-full ${
                  popupData.hasil_uji === 'Lulus' ? 'bg-green-100' : 'bg-red-100'
                }`}
              >
                {popupData.hasil_uji === 'Lulus' ? (
                  <svg
                    className="h-10 w-10 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                ) : (
                  <svg
                    className="h-10 w-10 text-red-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                )}
              </div>

              <h3
                className={`mt-4 text-xl font-bold ${
                  popupData.hasil_uji === 'Lulus' ? 'text-green-600' : 'text-red-600'
                }`}
              >
                {popupData.hasil_uji === 'Lulus'
                  ? 'Lulus Uji Emisi'
                  : 'Tidak Lulus Uji Emisi'}
              </h3>

              <p className="mt-2 text-gray-600">
                Status:{' '}
                <span className={popupData.valid ? 'text-green-600' : 'text-red-600'}>
                  {popupData.valid ? 'Valid' : 'Tidak Valid'}
                </span>
              </p>

              {popupData.popup_notes && popupData.popup_notes.length > 0 && (
                <div className="mt-4 text-left">
                  <h4 className="font-semibold text-gray-700 mb-2">Catatan Perbaikan:</h4>
                  <ul className="list-disc list-inside text-sm text-gray-600 space-y-1">
                    {popupData.popup_notes.map((note, index) => (
                      <li key={index}>{note}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>

            <div className="mt-6">
              <button
                onClick={() => {
                  setShowPopup(false);
                  setPopupData(null);
                }}
                className="w-full px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EmisiPage;
