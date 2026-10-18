/**
 * Membangun URL file statis yang bisa diakses browser.
 *
 * Path dari backend bisa berupa "uploads/home/x.pdf" (relatif, data lama) atau
 * "/uploads/home/x.pdf" (absolut, data baru). Path relatif akan salah resolve
 * ke origin frontend, sehingga harus dinormalisasi menjadi path absolut.
 */
export function fileUrl(filePath?: string | null): string {
  if (!filePath) return '';

  const path = filePath.trim();
  if (/^(https?:)?\/\//i.test(path) || /^(data|blob):/i.test(path)) return path;

  return '/' + path.replace(/^\/+/, '');
}

/**
 * Nama file untuk atribut `download` agar pengguna menyimpan file dengan nama
 * yang sama seperti aslinya, bukan nama UUID.
 */
export function fileNameFrom(filePath?: string | null, fallback = 'file'): string {
  if (!filePath) return fallback;
  const segments = filePath.split(/[\\/]/).filter(Boolean);
  return segments[segments.length - 1] || fallback;
}
