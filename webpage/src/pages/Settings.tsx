import { useEffect, useState } from 'react';
import { Server, Info, KeyRound, Users } from 'lucide-react';
import api from '@/lib/api';
import { useAuth } from '@/hooks/useAuth';
import UsersTab from '@/components/settings/UsersTab';
import PasswordTab from '@/components/settings/PasswordTab';

type TabKey = 'general' | 'engine' | 'password' | 'users' | 'about';

export default function Settings() {
  const [activeTab, setActiveTab] = useState<TabKey>('general');
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">系统设置</h1>
        <p className="text-text-secondary mt-1">配置管理平台参数</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Sidebar */}
        <div className="glass-card p-4 space-y-2 h-fit">
          <TabButton active={activeTab === 'general'} onClick={() => setActiveTab('general')} icon={<Server className="w-4 h-4" />}>
            通用配置
          </TabButton>
          <TabButton active={activeTab === 'engine'} onClick={() => setActiveTab('engine')} icon={<Server className="w-4 h-4" />}>
            引擎管理
          </TabButton>
          <TabButton active={activeTab === 'password'} onClick={() => setActiveTab('password')} icon={<KeyRound className="w-4 h-4" />}>
            修改密码
          </TabButton>
          {isAdmin && (
            <TabButton active={activeTab === 'users'} onClick={() => setActiveTab('users')} icon={<Users className="w-4 h-4" />}>
              用户管理
            </TabButton>
          )}
          <TabButton active={activeTab === 'about'} onClick={() => setActiveTab('about')} icon={<Info className="w-4 h-4" />}>
            关于
          </TabButton>
        </div>

        {/* Content */}
        <div className="lg:col-span-3">
          {activeTab === 'general' && <GeneralTab />}
          {activeTab === 'engine' && <EngineTab />}
          {activeTab === 'password' && <PasswordTab />}
          {activeTab === 'users' && <UsersTab />}
          {activeTab === 'about' && <AboutTab />}
        </div>
      </div>
    </div>
  );
}

function TabButton({ active, onClick, icon, children }: { active: boolean; onClick: () => void; icon: React.ReactNode; children: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all text-left ${
        active
          ? 'bg-accent-indigo/20 text-accent-indigo border border-accent-indigo/30'
          : 'text-text-secondary hover:bg-surface-overlay hover:text-text-primary'
      }`}
    >
      {icon}
      <span className="font-medium">{children}</span>
    </button>
  );
}

function GeneralTab() {
  const [config, setConfig] = useState<any>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const data: any = await api.get('/config');
      if (data.success) setConfig(data.data);
    } catch (error) {
      console.error('Failed to load config:', error);
    }
  };

  if (!config) return <div className="glass-card p-6 text-text-secondary">加载中...</div>;

  return (
    <div className="glass-card p-6 space-y-6">
      <h2 className="text-xl font-semibold text-text-primary">通用配置</h2>

      <div className="space-y-4">
        <ConfigItem label="引擎类型" value={config.engine || 'docker'} />
        <ConfigItem label="API 端口" value={config.port || 8080} />
        <ConfigItem label="监听地址" value={config.listen_addr || '127.0.0.1'} />
        <ConfigItem label="端口范围" value={config.port_range || 'N/A'} />
      </div>

      <div className="pt-4 border-t border-surface-border text-sm text-text-muted">
        配置文件位置: /opt/deskcli/conf/config.yaml
      </div>
    </div>
  );
}

function EngineTab() {
  const [engines, setEngines] = useState<any[]>([]);

  useEffect(() => {
    loadEngines();
  }, []);

  const loadEngines = async () => {
    try {
      const data: any = await api.get('/monitor/engines');
      if (data.success) setEngines(data.data || []);
    } catch (error) {
      console.error('Failed to load engines:', error);
    }
  };

  return (
    <div className="glass-card p-6 space-y-6">
      <h2 className="text-xl font-semibold text-text-primary">引擎管理</h2>

      <div className="space-y-4">
        {engines.map((engine) => (
          <div key={engine.name} className="flex items-center justify-between p-4 bg-surface-overlay rounded-lg border border-surface-border">
            <div className="flex items-center gap-3">
              <div className={`w-3 h-3 rounded-full ${engine.available ? 'bg-emerald-500' : 'bg-red-500'}`} />
              <div>
                <div className="font-medium text-text-primary capitalize">{engine.name}</div>
                <div className="text-sm text-text-muted">{engine.version || 'N/A'}</div>
              </div>
            </div>
            <span className={`px-3 py-1 rounded text-sm ${engine.available ? 'bg-emerald-500/20 text-emerald-500' : 'bg-red-500/20 text-red-500'}`}>
              {engine.available ? '可用' : '不可用'}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

function AboutTab() {
  return (
    <div className="glass-card p-6 space-y-6">
      <h2 className="text-xl font-semibold text-text-primary">关于 RDDocker</h2>

      <div className="space-y-4">
        <div>
          <div className="text-sm text-text-muted mb-1">版本</div>
          <div className="text-text-primary font-mono">v2.0.0</div>
        </div>

        <div>
          <div className="text-sm text-text-muted mb-1">项目介绍</div>
          <div className="text-text-secondary leading-relaxed">
            RDDocker 是一套开箱即用的桌面容器环境解决方案。通过 Docker/Podman/LXC 容器技术，
            以极低开销、秒级启动速度，创建出完整的 Linux 桌面环境。
          </div>
        </div>

        <div>
          <div className="text-sm text-text-muted mb-1">支持的操作系统</div>
          <div className="flex flex-wrap gap-2 mt-2">
            {['Debian', 'Ubuntu', 'Alpine', 'Fedora', 'Arch'].map((os) => (
              <span key={os} className="px-3 py-1 bg-accent-indigo/20 text-accent-indigo rounded-full text-sm border border-accent-indigo/30">
                {os}
              </span>
            ))}
          </div>
        </div>

        <div>
          <div className="text-sm text-text-muted mb-1">支持的桌面环境</div>
          <div className="flex flex-wrap gap-2 mt-2">
            {['Server', 'X11GUI', 'GNOME', 'Xfce', 'Plasma', 'Deepin', 'Lingmo', 'Hyprland', 'Niri'].map((de) => (
              <span key={de} className="px-3 py-1 bg-accent-cyan/20 text-accent-cyan rounded-full text-sm border border-accent-cyan/30">
                {de}
              </span>
            ))}
          </div>
        </div>

        <div className="pt-4 border-t border-surface-border">
          <div className="text-sm text-text-muted">
            © 2026 RDDocker Project | MIT License
          </div>
          <div className="mt-2">
            <a
              href="https://github.com/PIKACHUIM/RDDocker"
              target="_blank"
              rel="noopener noreferrer"
              className="text-accent-indigo hover:underline text-sm"
            >
              GitHub Repository →
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

function ConfigItem({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="flex items-center justify-between p-4 bg-surface-overlay rounded-lg border border-surface-border">
      <span className="text-text-muted">{label}</span>
      <span className="text-text-primary font-medium">{value}</span>
    </div>
  );
}
