import React, { useState, useEffect } from 'react';
import { homeAPI } from '../services/api';
import { HomeContent } from '../types';
import toast from 'react-hot-toast';
import { useAuth } from '../context/AuthContext';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [descriptions, setDescriptions] = useState<HomeContent[]>([]);
  const [images, setImages] = useState<HomeContent[]>([]);
  const [workingInstruction, setWorkingInstruction] = useState<HomeContent | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchContent();
  }, []);

  const fetchContent = async () => {
    try {
      const response = await homeAPI.getContent();
      if (response.success) {
        setDescriptions(response.data.descriptions);
        setImages(response.data.images);
        setWorkingInstruction(response.data.working_instruction);
      }
    } catch (error) {
      toast.error('Failed to load content');
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Selamat Datang</h1>
        <p className="mt-2 text-gray-600">
          Aplikasi Uji Emisi Kendaraan
        </p>
      </div>

      {/* Descriptions */}
      {descriptions.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Tentang Aplikasi</h2>
          {descriptions.map((desc) => (
            <div key={desc.id} className="mb-4">
              {desc.title && <h3 className="text-lg font-medium text-gray-700 mb-2">{desc.title}</h3>}
              <p className="text-gray-600 whitespace-pre-wrap">{desc.description}</p>
            </div>
          ))}
        </div>
      )}

      {/* Images */}
      {images.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Gambar</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {images.map((img) => (
              <div key={img.id} className="relative">
                <img
                  src={img.file_path}
                  alt={img.title || 'Gambar'}
                  className="w-full h-48 object-cover rounded-lg"
                />
                {img.title && (
                  <div className="absolute bottom-0 left-0 right-0 bg-black bg-opacity-50 text-white p-2 rounded-b-lg">
                    {img.title}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Working Instruction */}
      {workingInstruction && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Working Instruction</h2>
          <div className="border rounded-lg overflow-hidden">
            <iframe
              src={workingInstruction.file_path}
              className="w-full h-96"
              title="Working Instruction"
            />
          </div>
          <div className="mt-4">
            <a
              href={workingInstruction.file_path}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Download PDF
            </a>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <a
          href="/kendaraan"
          className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-blue-100 text-blue-600">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-900">Input Kendaraan</h3>
              <p className="text-sm text-gray-500">Tambah data kendaraan baru</p>
            </div>
          </div>
        </a>

        <a
          href="/emisi"
          className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-green-100 text-green-600">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-900">Input Emisi</h3>
              <p className="text-sm text-gray-500">Input hasil uji emisi</p>
            </div>
          </div>
        </a>

        <a
          href="/hasil-uji"
          className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-purple-100 text-purple-600">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-900">Data Hasil Uji</h3>
              <p className="text-sm text-gray-500">Lihat laporan dan statistik</p>
            </div>
          </div>
        </a>
      </div>
    </div>
  );
};

export default Dashboard;
