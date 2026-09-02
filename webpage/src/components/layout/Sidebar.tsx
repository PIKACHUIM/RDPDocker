import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Container, Package, Terminal, Settings } from 'lucide-react';
import { cn } from '@/lib/utils';

export default function Sidebar() {
  return (
    <aside className="w-64 bg-surface-raised border-r border-surface-border flex flex-col">
      <div className="p-6 border-b border-surface-border">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-accent-indigo to-accent-cyan bg-clip-text text-transparent">
          RDDocker
        </h1>
        <p className="text-xs text-text-muted mt-1">容器管理平台</p>
      </div>

      <nav className="flex-1 p-4 space-y-2">
        <NavItem to="/" icon={<LayoutDashboard className="w-5 h-5" />} label="仪表板" />
        <NavItem to="/containers" icon={<Container className="w-5 h-5" />} label="容器管理" />
        <NavItem to="/templates" icon={<Package className="w-5 h-5" />} label="模板市场" />
        <NavItem to="/terminal" icon={<Terminal className="w-5 h-5" />} label="终端" />
        <NavItem to="/settings" icon={<Settings className="w-5 h-5" />} label="系统设置" />
      </nav>

      <div className="p-4 border-t border-surface-border">
        <div className="text-xs text-text-muted">
          <div>版本: v2.0.0</div>
          <div className="mt-1">© 2026 RDDocker</div>
        </div>
      </div>
    </aside>
  );
}

interface NavItemProps {
  to: string;
  icon: React.ReactNode;
  label: string;
}

function NavItem({ to, icon, label }: NavItemProps) {
  return (
    <NavLink
      to={to}
      end={to === '/'}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-3 px-4 py-3 rounded-lg transition-all",
          isActive
            ? "bg-accent-indigo/20 text-accent-indigo border border-accent-indigo/30"
            : "text-text-secondary hover:bg-surface-overlay hover:text-text-primary"
        )
      }
    >
      {icon}
      <span className="font-medium">{label}</span>
    </NavLink>
  );
}
