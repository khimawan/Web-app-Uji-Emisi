import React, { useEffect, useRef, useState } from 'react';

export type TemplateFormat = 'csv' | 'excel';

/**
 * Satu tombol "Template" dengan menu dropdown CSV / Excel, menggantikan dua
 * tombol terpisah. Dipakai bersama oleh halaman Input Kendaraan dan Data Hasil Uji
 * supaya formatnya konsisten.
 */
const TemplateDownloadMenu: React.FC<{
  onDownload: (format: TemplateFormat) => void;
  disabled?: boolean;
}> = ({ onDownload, disabled = false }) => {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!isOpen) return;

    const onPointerDown = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setIsOpen(false);
    };

    document.addEventListener('mousedown', onPointerDown);
    document.addEventListener('keydown', onKeyDown);

    return () => {
      document.removeEventListener('mousedown', onPointerDown);
      document.removeEventListener('keydown', onKeyDown);
    };
  }, [isOpen]);

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setIsOpen((value) => !value)}
        disabled={disabled}
        aria-haspopup="menu"
        aria-expanded={isOpen}
        className="inline-flex items-center gap-1.5 px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 disabled:cursor-not-allowed disabled:opacity-50"
      >
        Template
        <svg
          className={`h-4 w-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {isOpen && (
        <div
          role="menu"
          // Di mobile tombol berada di sisi kiri, jadi dropdown ditambat ke
          // kiri (left-1) supaya melebar ke kanan dan tidak terpotong teksnya.
          // Di layar lebar kembali rata kanan, sedikit menjorok ke kanan.
          className="absolute left-1 z-30 mt-1 w-44 max-w-[calc(100vw-2.5rem)] overflow-hidden rounded-md border border-gray-200 bg-white py-1 shadow-lg sm:left-auto sm:-right-2"
        >
          {(['csv', 'excel'] as const).map((format) => (
            <button
              key={format}
              type="button"
              role="menuitem"
              onClick={() => {
                setIsOpen(false);
                onDownload(format);
              }}
              className="block w-full whitespace-nowrap px-4 py-2 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100"
            >
              Template {format === 'csv' ? 'CSV' : 'Excel'}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default TemplateDownloadMenu;