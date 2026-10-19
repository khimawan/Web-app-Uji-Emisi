import React, { useCallback, useEffect, useRef, useState } from 'react';
import * as pdfjsLib from 'pdfjs-dist';
import type { PDFDocumentProxy, RenderTask } from 'pdfjs-dist';
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url';

pdfjsLib.GlobalWorkerOptions.workerSrc = workerUrl;

type PageStatus = 'idle' | 'rendering' | 'done' | 'error';

const ZOOM_STEPS = [1, 1.25, 1.5, 2];

/**
 * Satu halaman PDF. Dirender ke <canvas> hanya saat halaman ini mendekati
 * viewport, supaya dokumen panjang tetap ringan di ponsel.
 */
const PdfPage: React.FC<{
  doc: PDFDocumentProxy;
  pageNumber: number;
  scale: number;
  shouldRender: boolean;
  onEnterView: (pageNumber: number) => void;
}> = ({ doc, pageNumber, scale, shouldRender, onEnterView }) => {
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const renderTaskRef = useRef<RenderTask | null>(null);
  const [status, setStatus] = useState<PageStatus>('idle');

  useEffect(() => {
    const element = wrapperRef.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) onEnterView(pageNumber);
        });
      },
      { rootMargin: '600px 0px' }
    );

    observer.observe(element);
    return () => observer.disconnect();
  }, [onEnterView, pageNumber]);

  useEffect(() => {
    if (!shouldRender) return;

    let cancelled = false;

    const render = async () => {
      setStatus('rendering');
      const page = await doc.getPage(pageNumber);
      if (cancelled) return;

      const canvas = canvasRef.current;
      if (!canvas) return;

      const cssViewport = page.getViewport({ scale });
      const dpr = Math.min(window.devicePixelRatio || 1, 2);
      const bitmapViewport = page.getViewport({ scale: scale * dpr });

      canvas.width = Math.floor(bitmapViewport.width);
      canvas.height = Math.floor(bitmapViewport.height);
      canvas.style.width = `${Math.floor(cssViewport.width)}px`;
      canvas.style.height = `${Math.floor(cssViewport.height)}px`;

      const context = canvas.getContext('2d');
      if (context) {
        context.fillStyle = '#ffffff';
        context.fillRect(0, 0, canvas.width, canvas.height);
      }

      renderTaskRef.current?.cancel();
      const task = page.render({ canvas, viewport: bitmapViewport });
      renderTaskRef.current = task;

      try {
        await task.promise;
        if (!cancelled) setStatus('done');
      } catch (error) {
        const name = (error as { name?: string })?.name;
        if (!cancelled && name !== 'RenderingCancelledException') setStatus('error');
      }
    };

    render().catch(() => {
      if (!cancelled) setStatus('error');
    });

    return () => {
      cancelled = true;
      renderTaskRef.current?.cancel();
    };
  }, [doc, pageNumber, scale, shouldRender]);

  return (
    <div ref={wrapperRef} className="mb-4 flex justify-center">
      <div className="relative w-full" style={{ maxWidth: '100%' }}>
        <canvas
          ref={canvasRef}
          className={`mx-auto block max-w-full rounded border border-gray-300 shadow-sm ${status === 'done' ? '' : 'invisible'}`}
          aria-label={`Halaman ${pageNumber}`}
        />
        {status !== 'done' && (
          <div className="absolute inset-0 flex items-center justify-center rounded border border-gray-300 bg-gray-100">
            {status === 'error' ? (
              <span className="px-4 text-center text-xs text-red-600">
                Halaman {pageNumber} gagal ditampilkan.
              </span>
            ) : (
              <span className="flex items-center gap-2 text-xs text-gray-500">
                <span className="h-3 w-3 animate-spin rounded-full border-2 border-gray-400 border-t-transparent" />
                Memuat halaman {pageNumber}...
              </span>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

const PdfViewer: React.FC<{ url: string; title?: string }> = ({
  url,
  title = 'Dokumen PDF',
}) => {
  const [doc, setDoc] = useState<PDFDocumentProxy | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading');
  const [errorMessage, setErrorMessage] = useState('');
  const [baseWidth, setBaseWidth] = useState(0);
  const [containerWidth, setContainerWidth] = useState(0);
  const [zoomIndex, setZoomIndex] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [renderedPages, setRenderedPages] = useState<number[]>([1]);
  const measureRef = useRef<HTMLDivElement | null>(null);

  // Muat dokumen PDF.
  useEffect(() => {
    if (!url) {
      setStatus('error');
      setErrorMessage('Tidak ada file PDF.');
      return;
    }

    let cancelled = false;
    setStatus('loading');
    setErrorMessage('');
    setDoc(null);
    setCurrentPage(1);
    setZoomIndex(0);
    setRenderedPages([1]);

    const loadingTask = pdfjsLib.getDocument({ url });

    loadingTask.promise
      .then((loaded) => {
        if (cancelled) return;
        setDoc(loaded);
        setStatus('ready');
      })
      .catch((error: unknown) => {
        if (cancelled) return;
        setErrorMessage(
          error instanceof Error ? error.message : 'Gagal memuat dokumen PDF.'
        );
        setStatus('error');
      });

    return () => {
      cancelled = true;
      void loadingTask.destroy().catch(() => undefined);
    };
  }, [url]);

  // Ukur lebar area tampilan supaya halaman muat ke layar (fit to width).
  useEffect(() => {
    const element = measureRef.current;
    if (!element) return;

    const update = () => setContainerWidth(element.clientWidth);
    update();

    const observer = new ResizeObserver(update);
    observer.observe(element);
    return () => observer.disconnect();
  }, [status]);

  // Lebar halaman PDF dalam poin, dipakai untuk menghitung skala.
  useEffect(() => {
    if (!doc) return;

    let cancelled = false;
    doc
      .getPage(1)
      .then((page) => {
        if (!cancelled) setBaseWidth(page.getViewport({ scale: 1 }).width);
      })
      .catch(() => undefined);

    return () => {
      cancelled = true;
    };
  }, [doc]);

  const handleEnterView = useCallback((pageNumber: number) => {
    setCurrentPage(pageNumber);
    setRenderedPages((previous) =>
      previous.includes(pageNumber) ? previous : [...previous, pageNumber]
    );
  }, []);

  const fitScale =
    baseWidth > 0 && containerWidth > 0
      ? Math.min((containerWidth - 2) / baseWidth, 2)
      : 1;
  const scale = fitScale * ZOOM_STEPS[zoomIndex];

  const controlClass =
    'inline-flex h-10 min-w-[2.5rem] items-center justify-center rounded-md border border-gray-300 bg-white px-3 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-40';

  return (
    <div className="overflow-hidden rounded-lg border border-gray-300 bg-gray-50">
      {/* Kontrol */}
      <div className="flex flex-wrap items-center gap-2 border-b border-gray-300 bg-white px-3 py-2">
        <button
          type="button"
          className={controlClass}
          onClick={() => setZoomIndex((value) => Math.max(0, value - 1))}
          disabled={status !== 'ready' || zoomIndex === 0}
          aria-label="Perkecil"
        >
          −
        </button>
        <span className="min-w-[3.5rem] text-center text-sm font-medium text-gray-700">
          {Math.round(ZOOM_STEPS[zoomIndex] * 100)}%
        </span>
        <button
          type="button"
          className={controlClass}
          onClick={() =>
            setZoomIndex((value) => Math.min(ZOOM_STEPS.length - 1, value + 1))
          }
          disabled={status !== 'ready' || zoomIndex === ZOOM_STEPS.length - 1}
          aria-label="Perbesar"
        >
          +
        </button>

        <span className="ml-auto text-sm text-gray-600">
          {status === 'ready' ? `Hal. ${currentPage} / ${doc?.numPages ?? 0}` : ' '}
        </span>

        <a
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex h-10 items-center rounded-md bg-blue-600 px-3 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          Buka PDF
        </a>
      </div>

      {/* Area dokumen */}
      <div
        ref={measureRef}
        className="max-h-[70vh] overflow-auto overscroll-contain bg-gray-100 p-2 sm:max-h-[80vh] md:p-4"
      >
        {status === 'loading' && (
          <div className="flex items-center justify-center gap-2 py-16 text-sm text-gray-500">
            <span className="h-4 w-4 animate-spin rounded-full border-2 border-gray-400 border-t-transparent" />
            Memuat dokumen...
          </div>
        )}

        {status === 'error' && (
          <div className="px-4 py-10 text-center">
            <p className="mb-4 text-sm text-red-600">
              {errorMessage || 'Dokumen tidak dapat ditampilkan.'}
            </p>
            <a
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-10 items-center rounded-md bg-blue-600 px-4 text-sm font-medium text-white hover:bg-blue-700"
            >
              Buka di tab baru
            </a>
          </div>
        )}

        {status === 'ready' && doc && (
          <div>
            {Array.from({ length: doc.numPages }, (_, index) => index + 1).map(
              (pageNumber) => (
                <PdfPage
                  key={`${url}-${pageNumber}`}
                  doc={doc}
                  pageNumber={pageNumber}
                  scale={scale}
                  shouldRender={renderedPages.includes(pageNumber)}
                  onEnterView={handleEnterView}
                />
              )
            )}
          </div>
        )}
      </div>

      <p className="sr-only">{title}</p>
    </div>
  );
};

export default PdfViewer;