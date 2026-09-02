import { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from '@/hooks/useAuth';
import LoginPage from '@/pages/Login';
import MainLayout from '@/components/layout/MainLayout';

// Split per route so the login screen does not pull in recharts (Dashboard)
// or xterm.js (Terminal), which dominate the bundle.
const Dashboard = lazy(() => import('@/pages/Dashboard'));
const Containers = lazy(() => import('@/pages/Containers'));
const Templates = lazy(() => import('@/pages/Templates'));
const Terminal = lazy(() => import('@/pages/Terminal'));
const Settings = lazy(() => import('@/pages/Settings'));

function PageFallback() {
  return (
    <div className="flex flex-1 items-center justify-center py-24">
      <div className="flex items-center gap-3 text-text-secondary">
        <span className="h-2 w-2 animate-pulse rounded-full bg-accent-indigo" />
        <span className="text-sm">加载中…</span>
      </div>
    </div>
  );
}

function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-surface">
        <div className="text-text-secondary">Loading...</div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Suspense fallback={<PageFallback />}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/"
              element={
                <AuthGuard>
                  <MainLayout />
                </AuthGuard>
              }
            >
              <Route index element={<Dashboard />} />
              <Route path="containers" element={<Containers />} />
              <Route path="templates" element={<Templates />} />
              <Route path="terminal" element={<Terminal />} />
              <Route path="settings" element={<Settings />} />
            </Route>
          </Routes>
        </Suspense>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
