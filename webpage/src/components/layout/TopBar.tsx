import { useAuth } from '@/hooks/useAuth';
import { LogOut, User } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function TopBar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="h-16 bg-surface-raised border-b border-surface-border flex items-center justify-between px-6">
      <div className="flex items-center gap-2 text-sm text-text-secondary">
        <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
        <span>引擎已连接</span>
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2 text-sm">
          <User className="w-4 h-4 text-text-muted" />
          <span className="text-text-secondary">{user?.username}</span>
          <span className="px-2 py-0.5 rounded text-xs bg-accent-indigo/20 text-accent-indigo border border-accent-indigo/30">
            {user?.role}
          </span>
        </div>
        
        <button
          onClick={handleLogout}
          className="p-2 rounded-lg hover:bg-surface-overlay text-text-muted hover:text-text-primary transition-colors"
          title="退出登录"
        >
          <LogOut className="w-4 h-4" />
        </button>
      </div>
    </header>
  );
}
