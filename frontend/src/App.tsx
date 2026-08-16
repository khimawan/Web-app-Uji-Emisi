import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { AuthProvider, useAuth } from './context/AuthContext';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import KendaraanPage from './pages/KendaraanPage';
import EmisiPage from './pages/EmisiPage';
import HasilUjiPage from './pages/HasilUjiPage';
import ParameterPage from './pages/ParameterPage';
import AdminPage from './pages/AdminPage';
import Layout from './components/Layout';

const ProtectedRoute: React.FC<{ children: React.ReactNode; requiredRole?: string[] }> = ({ children, requiredRole }) => {
  const { isAuthenticated, user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && user && !requiredRole.includes(user.role)) {
    return <Navigate to="/" replace />;
  }

  return <Layout>{children}</Layout>;
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

function App() {
  return (
    <AuthProvider>
      <Router>
        <Toaster position="top-right" />
        <Routes>
          <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
          <Route path="/" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
          <Route path="/kendaraan" element={<ProtectedRoute><KendaraanPage /></ProtectedRoute>} />
          <Route path="/emisi" element={<ProtectedRoute><EmisiPage /></ProtectedRoute>} />
          <Route path="/hasil-uji" element={<ProtectedRoute><HasilUjiPage /></ProtectedRoute>} />
          <Route path="/parameter" element={<ProtectedRoute requiredRole={['admin', 'supervisor']}><ParameterPage /></ProtectedRoute>} />
          <Route path="/admin" element={<ProtectedRoute requiredRole={['admin', 'supervisor']}><AdminPage /></ProtectedRoute>} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
