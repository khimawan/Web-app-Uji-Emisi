import React, { useState, useEffect } from 'react';
import { parametersAPI, popupNotesAPI } from '../services/api';
import { Parameter, PopupNote } from '../types';
import toast from 'react-hot-toast';

const ParameterPage: React.FC = () => {
  const [parameters, setParameters] = useState<Parameter[]>([]);
  const [popupNotes, setPopupNotes] = useState<PopupNote[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showParamForm, setShowParamForm] = useState(false);
  const [showNoteForm, setShowNoteForm] = useState(false);
  const [editingParamId, setEditingParamId] = useState<string | null>(null);
  const [editingNoteId, setEditingNoteId] = useState<string | null>(null);

  const [paramForm, setParamForm] = useState({
    kategori: 'Bensin M',
    tahun_min: 0,
    tahun_max: null as number | null,
    tahun_operator: '<' as '<' | '>' | 'between',
    co_max: null as number | null,
    hc_max: null as number | null,
    opasitas_max: null as number | null,
  });

  const [noteForm, setNoteForm] = useState({
    jenis_kendaraan: 'Bensin M',
    tahun_operator: '>' as '>' | '<' | '>=' | '<=' | '=',
    tahun_value: 1990,
    parameter_uji: 'CO' as 'CO' | 'CO2' | 'HC' | 'O2' | 'Lambda' | 'Opasitas',
    nilai_operator: '>' as '>' | '<' | '>=' | '<=' | '=',
    nilai_value: 0,
    note: '',
  });

  useEffect(() => {
    fetchParameters();
    fetchPopupNotes();
  }, []);

  const fetchParameters = async () => {
    try {
      const response = await parametersAPI.getAll();
      if (response.success) {
        setParameters(response.data);
      }
    } catch (error) {
      toast.error('Failed to fetch parameters');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchPopupNotes = async () => {
    try {
      const response = await popupNotesAPI.getAll();
      if (response.success) {
        setPopupNotes(response.data);
      }
    } catch (error) {
      toast.error('Failed to fetch popup notes');
    }
  };

  const handleParamSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingParamId) {
        await parametersAPI.update(editingParamId, paramForm);
        toast.success('Parameter updated successfully');
      } else {
        await parametersAPI.create(paramForm);
        toast.success('Parameter added successfully');
      }
      setShowParamForm(false);
      setEditingParamId(null);
      resetParamForm();
      fetchParameters();
    } catch (error) {
      toast.error('Failed to save parameter');
    }
  };

  const handleNoteSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingNoteId) {
        await popupNotesAPI.update(editingNoteId, noteForm);
        toast.success('Popup note updated successfully');
      } else {
        await popupNotesAPI.create(noteForm);
        toast.success('Popup note added successfully');
      }
      setShowNoteForm(false);
      setEditingNoteId(null);
      resetNoteForm();
      fetchPopupNotes();
    } catch (error) {
      toast.error('Failed to save popup note');
    }
  };

  const handleEditParam = (param: Parameter) => {
    setParamForm({
      kategori: param.kategori,
      tahun_min: param.tahun_min,
      tahun_max: param.tahun_max,
      tahun_operator: param.tahun_operator,
      co_max: param.co_max,
      hc_max: param.hc_max,
      opasitas_max: param.opasitas_max,
    });
    setEditingParamId(param.id);
    setShowParamForm(true);
  };

  const handleEditNote = (note: PopupNote) => {
    setNoteForm({
      jenis_kendaraan: note.jenis_kendaraan,
      tahun_operator: note.tahun_operator,
      tahun_value: note.tahun_value,
      parameter_uji: note.parameter_uji,
      nilai_operator: note.nilai_operator,
      nilai_value: note.nilai_value,
      note: note.note,
    });
    setEditingNoteId(note.id);
    setShowNoteForm(true);
  };

  const handleDeleteParam = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this parameter?')) {
      try {
        await parametersAPI.delete(id);
        toast.success('Parameter deleted successfully');
        fetchParameters();
      } catch (error) {
        toast.error('Failed to delete parameter');
      }
    }
  };

  const handleDeleteNote = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this popup note?')) {
      try {
        await popupNotesAPI.delete(id);
        toast.success('Popup note deleted successfully');
        fetchPopupNotes();
      } catch (error) {
        toast.error('Failed to delete popup note');
      }
    }
  };

  const resetParamForm = () => {
    setParamForm({
      kategori: 'Bensin M',
      tahun_min: 0,
      tahun_max: null,
      tahun_operator: '<',
      co_max: null,
      hc_max: null,
      opasitas_max: null,
    });
  };

  const resetNoteForm = () => {
    setNoteForm({
      jenis_kendaraan: 'Bensin M',
      tahun_operator: '>',
      tahun_value: 1990,
      parameter_uji: 'CO',
      nilai_operator: '>',
      nilai_value: 0,
      note: '',
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
    <div className="space-y-8">
      <h1 className="text-2xl font-bold text-gray-900">Parameter & Konfigurasi</h1>

      {/* Parameters Section */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold">Parameter Standar Emisi</h2>
          <button
            onClick={() => { resetParamForm(); setEditingParamId(null); setShowParamForm(true); }}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            + Tambah Parameter
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Kategori</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tahun</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">CO(%)</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">HC(ppm)</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Opasitas(%)</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {parameters.map((param) => (
                <tr key={param.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900">
                    {param.kategori}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {param.tahun_operator === 'between'
                      ? `${param.tahun_min}-${param.tahun_max}`
                      : param.tahun_operator === '<'
                      ? `<${param.tahun_min}`
                      : `>${param.tahun_min}`}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {param.co_max ?? '-'}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {param.hc_max ?? '-'}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {param.opasitas_max ?? '-'}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEditParam(param)}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteParam(param.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Popup Notes Section */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold">Konfigurasi Pop-up Notifikasi</h2>
          <button
            onClick={() => { resetNoteForm(); setEditingNoteId(null); setShowNoteForm(true); }}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            + Tambah Pop-up Note
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Jenis Kendaraan</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tahun</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Parameter</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nilai</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Note</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {popupNotes.map((note) => (
                <tr key={note.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 text-sm text-gray-900">{note.jenis_kendaraan}</td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {note.tahun_operator}{note.tahun_value}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {note.parameter_uji}
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                    {note.nilai_operator}{note.nilai_value}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500 max-w-xs truncate">{note.note}</td>
                  <td className="px-4 py-3 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleEditNote(note)}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteNote(note.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Parameter Form Modal */}
      {showParamForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {editingParamId ? 'Edit Parameter' : 'Tambah Parameter'}
            </h2>
            <form onSubmit={handleParamSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Kategori</label>
                <select
                  value={paramForm.kategori}
                  onChange={(e) => setParamForm({ ...paramForm, kategori: e.target.value })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="Bensin M">Bensin M</option>
                  <option value="Bensin N&O">Bensin N&O</option>
                  <option value="Solar JBB">Solar JBB</option>
                  <option value="Solar GVW">Solar GVW</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Operator Tahun</label>
                <select
                  value={paramForm.tahun_operator}
                  onChange={(e) => setParamForm({ ...paramForm, tahun_operator: e.target.value as any })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="<">&lt; (Kurang dari)</option>
                  <option value=">">&gt; (Lebih dari)</option>
                  <option value="between">Between (Antara)</option>
                </select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Tahun Min</label>
                  <input
                    type="number"
                    value={paramForm.tahun_min}
                    onChange={(e) => setParamForm({ ...paramForm, tahun_min: parseInt(e.target.value) })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  />
                </div>
                {paramForm.tahun_operator === 'between' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Tahun Max</label>
                    <input
                      type="number"
                      value={paramForm.tahun_max ?? ''}
                      onChange={(e) => setParamForm({ ...paramForm, tahun_max: parseInt(e.target.value) || null })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                )}
              </div>

              {!paramForm.kategori.startsWith('Solar') && (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">CO Max (%)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={paramForm.co_max ?? ''}
                      onChange={(e) => setParamForm({ ...paramForm, co_max: parseFloat(e.target.value) || null })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">HC Max (ppm)</label>
                    <input
                      type="number"
                      step="0.01"
                      value={paramForm.hc_max ?? ''}
                      onChange={(e) => setParamForm({ ...paramForm, hc_max: parseFloat(e.target.value) || null })}
                      className="mt-1 block w-full border rounded-md px-3 py-2"
                    />
                  </div>
                </div>
              )}

              {paramForm.kategori.startsWith('Solar') && (
                <div>
                  <label className="block text-sm font-medium text-gray-700">Opasitas Max (%)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={paramForm.opasitas_max ?? ''}
                    onChange={(e) => setParamForm({ ...paramForm, opasitas_max: parseFloat(e.target.value) || null })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  />
                </div>
              )}

              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => { setShowParamForm(false); setEditingParamId(null); }}
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

      {/* Popup Note Form Modal */}
      {showNoteForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-4">
              {editingNoteId ? 'Edit Pop-up Note' : 'Tambah Pop-up Note'}
            </h2>
            <form onSubmit={handleNoteSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Jenis Kendaraan</label>
                <input
                  type="text"
                  value={noteForm.jenis_kendaraan}
                  onChange={(e) => setNoteForm({ ...noteForm, jenis_kendaraan: e.target.value })}
                  placeholder="Contoh: Bensin M, Bensin N&O"
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                />
                <p className="text-xs text-gray-500 mt-1">Pisahkan dengan koma untuk beberapa jenis</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Operator Tahun</label>
                  <select
                    value={noteForm.tahun_operator}
                    onChange={(e) => setNoteForm({ ...noteForm, tahun_operator: e.target.value as any })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  >
                    <option value=">">&gt;</option>
                    <option value="<">&lt;</option>
                    <option value=">=">&gt;=</option>
                    <option value="<=">&lt;=</option>
                    <option value="=">=</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Tahun</label>
                  <input
                    type="number"
                    value={noteForm.tahun_value}
                    onChange={(e) => setNoteForm({ ...noteForm, tahun_value: parseInt(e.target.value) })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Parameter Uji</label>
                <select
                  value={noteForm.parameter_uji}
                  onChange={(e) => setNoteForm({ ...noteForm, parameter_uji: e.target.value as any })}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                >
                  <option value="CO">CO</option>
                  <option value="CO2">CO2</option>
                  <option value="HC">HC</option>
                  <option value="O2">O2</option>
                  <option value="Lambda">Lambda</option>
                  <option value="Opasitas">Opasitas</option>
                </select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Operator Nilai</label>
                  <select
                    value={noteForm.nilai_operator}
                    onChange={(e) => setNoteForm({ ...noteForm, nilai_operator: e.target.value as any })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  >
                    <option value=">">&gt;</option>
                    <option value="<">&lt;</option>
                    <option value=">=">&gt;=</option>
                    <option value="<=">&lt;=</option>
                    <option value="=">=</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Nilai</label>
                  <input
                    type="number"
                    step="0.01"
                    value={noteForm.nilai_value}
                    onChange={(e) => setNoteForm({ ...noteForm, nilai_value: parseFloat(e.target.value) })}
                    className="mt-1 block w-full border rounded-md px-3 py-2"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">Note</label>
                <textarea
                  value={noteForm.note}
                  onChange={(e) => setNoteForm({ ...noteForm, note: e.target.value })}
                  rows={3}
                  className="mt-1 block w-full border rounded-md px-3 py-2"
                />
              </div>

              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => { setShowNoteForm(false); setEditingNoteId(null); }}
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

export default ParameterPage;
