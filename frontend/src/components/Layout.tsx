import React, { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navItems = [
    { path: '/', label: 'Home', roles: ['admin', 'supervisor', 'anggota'] },
    { path: '/kendaraan', label: 'Input Kendaraan', roles: ['admin', 'supervisor', 'anggota'] },
    { path: '/emisi', label: 'Input Emisi', roles: ['admin', 'supervisor', 'anggota'] },
    { path: '/hasil-uji', label: 'Data Hasil Uji', roles: ['admin', 'supervisor', 'anggota'] },
    { path: '/parameter', label: 'Parameter', roles: ['admin', 'supervisor'] },
    { path: '/admin', label: 'Users', roles: ['admin', 'supervisor'] },
  ];

  const filteredNavItems = navItems.filter(item =>
    user && item.roles.includes(user.role)
  );

  const isActive = (path: string) => location.pathname === path;

  // Tutup sidebar setiap kali pindah halaman.
  useEffect(() => {
    setIsSidebarOpen(false);
  }, [location.pathname]);

  // Kunci scroll body dan tutup dengan tombol Escape selama sidebar terbuka.
  useEffect(() => {
    if (!isSidebarOpen) return;

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setIsSidebarOpen(false);
    };

    document.addEventListener('keydown', onKeyDown);
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';

    return () => {
      document.removeEventListener('keydown', onKeyDown);
      document.body.style.overflow = previousOverflow;
    };
  }, [isSidebarOpen]);

  const linkClasses = (active: boolean) =>
    `block truncate rounded-md px-3 py-2 text-sm font-medium transition-colors ${
      active
        ? 'bg-blue-800 text-white'
        : 'text-blue-100 hover:bg-blue-600 hover:text-white'
    }`;

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="sticky top-0 z-30 bg-blue-600 text-white shadow-lg">
        <div className="mx-auto max-w-7xl px-3 sm:px-4">
          <div className="flex h-14 items-center justify-between gap-2 md:h-16">
            {/* Kiri: tombol menu (hanya mobile) + logo */}
            <div className="flex min-w-0 items-center gap-1 sm:gap-2">
              <button
                type="button"
                onClick={() => setIsSidebarOpen(true)}
                aria-label="Buka menu navigasi"
                aria-expanded={isSidebarOpen}
                aria-controls="sidebar-nav"
                className="-ml-1 inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-md text-blue-50 hover:bg-blue-700 focus:outline-none focus-visible:ring-2 focus-visible:ring-white md:hidden"
              >
                <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>

              <Link to="/" className="flex min-w-0 items-center gap-2">
                <svg className="h-7 w-7 shrink-0 md:h-8 md:w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="truncate text-base font-bold md:text-xl">Uji Emisi</span>
              </Link>
            </div>

            {/* Tengah: navigasi horizontal (desktop) */}
            <nav className="hidden items-center gap-1 md:flex" aria-label="Navigasi utama">
              {filteredNavItems.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  aria-current={isActive(item.path) ? 'page' : undefined}
                  className={`truncate rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                    isActive(item.path)
                      ? 'bg-blue-700 text-white'
                      : 'text-blue-100 hover:bg-blue-500 hover:text-white'
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </nav>

            {/* Kanan: nama pengguna + logout */}
            <div className="flex shrink-0 items-center gap-2">
              <span className="hidden max-w-[12rem] truncate text-sm text-blue-100 sm:inline">
                {user?.nama} ({user?.role})
              </span>
              <button
                type="button"
                onClick={handleLogout}
                className="inline-flex items-center gap-1.5 rounded-md bg-red-500 px-2.5 py-1.5 text-xs font-medium transition-colors hover:bg-red-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-white sm:px-3 sm:py-2 sm:text-sm"
              >
                <svg className="h-4 w-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                <span className="hidden sm:inline">Logout</span>
                <span className="sr-only sm:hidden">Logout</span>
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Sidebar mobile */}
      <div
        className={`fixed inset-0 z-40 bg-black/50 transition-opacity duration-200 md:hidden ${
          isSidebarOpen ? 'opacity-100' : 'pointer-events-none opacity-0'
        }`}
        onClick={() => setIsSidebarOpen(false)}
        aria-hidden="true"
      />

      <aside
        id="sidebar-nav"
        className={`fixed left-0 top-0 z-50 flex h-full w-64 max-w-[82vw] flex-col bg-blue-700 text-white shadow-xl transition-transform duration-200 md:hidden ${
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
        aria-hidden={!isSidebarOpen}
      >
        <div className="flex h-14 shrink-0 items-center justify-between border-b border-blue-600 px-3">
          <span className="truncate font-bold">Menu</span>
          <button
            type="button"
            onClick={() => setIsSidebarOpen(false)}
            aria-label="Tutup menu navigasi"
            className="-mr-1 inline-flex h-10 w-10 items-center justify-center rounded-md text-blue-100 hover:bg-blue-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-white"
          >
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav className="flex-1 space-y-1 overflow-y-auto overscroll-contain p-2" aria-label="Navigasi utama">
          {filteredNavItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              tabIndex={isSidebarOpen ? undefined : -1}
              aria-current={isActive(item.path) ? 'page' : undefined}
              className={linkClasses(isActive(item.path))}
            >
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="shrink-0 border-t border-blue-600 p-3">
          <p className="truncate text-sm text-blue-100">{user?.nama}</p>
          <p className="truncate text-xs capitalize text-blue-200">{user?.role}</p>
        </div>
      </aside>

      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
};

export default Layout;