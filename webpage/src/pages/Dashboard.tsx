import { useEffect, useState } from 'react';
import { Play, Activity, Package, Server } from 'lucide-react';
import api, { EngineStats, Container } from '@/lib/api';
import { cn } from '@/lib/utils';

export default function Dashboard() {
  const [stats, setStats] = useState<EngineStats | null>(null);
  const [recentContainers, setRecentContainers] = useState<Container[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [statsData, containersData]: any = await Promise.all([
        api.get('/monitor/stats'),
        api.get('/containers'),
      ]);
      
      if (statsData.success) setStats(statsData.data);
      if (containersData.success) setRecentContainers(containersData.data?.slice(0, 5) || []);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return <div className="text-text-secondary">加载中...</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">仪表板</h1>
        <p className="text-text-secondary mt-1">系统概览与快速操作</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          icon={<Play className="w-6 h-6" />}
          label="运行中容器"
          value={stats?.containers_running || 0}
          color="emerald"
        />
        <StatCard
          icon={<Activity className="w-6 h-6" />}
          label="总容器数"
          value={stats?.containers_total || 0}
          color="cyan"
        />
        <StatCard
          icon={<Package className="w-6 h-6" />}
          label="镜像总数"
          value={stats?.images_total || 0}
          color="indigo"
        />
        <StatCard
          icon={<Server className="w-6 h-6" />}
          label="引擎版本"
          value={stats?.engine_version || 'N/A'}
          color="indigo"
          valueIsString
        />
      </div>

      {/* Recent Containers */}
      <div className="glass-card p-6">
        <h2 className="text-xl font-semibold text-text-primary mb-4">最近容器</h2>
        <div className="space-y-3">
          {recentContainers.length === 0 ? (
            <div className="text-center py-8 text-text-muted">
              暂无容器，前往模板市场创建第一个容器
            </div>
          ) : (
            recentContainers.map((container) => (
              <div
                key={container.name}
                className="flex items-center justify-between p-4 bg-surface-overlay rounded-lg border border-surface-border hover:border-accent-indigo/30 transition-colors"
              >
                <div className="flex items-center gap-4">
                  <div className={cn('status-dot', container.status === 'running' ? 'running' : 'stopped')} />
                  <div>
                    <div className="font-medium text-text-primary">{container.name}</div>
                    <div className="text-sm text-text-muted">{container.image}</div>
                  </div>
                </div>
                <div className="text-sm text-text-secondary max-w-[50%] text-right">
                  {(container.ports || []).join(', ')}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

interface StatCardProps {
  icon: React.ReactNode;
  label: string;
  value: number | string;
  color: 'emerald' | 'cyan' | 'indigo';
  valueIsString?: boolean;
}

function StatCard({ icon, label, value, color, valueIsString }: StatCardProps) {
  const colorMap = {
    emerald: 'from-emerald-500/20 to-emerald-500/5 border-emerald-500/30 text-emerald-500',
    cyan: 'from-cyan-500/20 to-cyan-500/5 border-cyan-500/30 text-cyan-500',
    indigo: 'from-indigo-500/20 to-indigo-500/5 border-indigo-500/30 text-indigo-500',
  };

  return (
    <div className={cn('glass-card p-6 border bg-gradient-to-br', colorMap[color])}>
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm text-text-secondary mb-2">{label}</div>
          <div className="text-3xl font-bold text-text-primary">
            {valueIsString ? value : typeof value === 'number' ? value.toLocaleString() : value}
          </div>
        </div>
        <div className="opacity-50">{icon}</div>
      </div>
    </div>
  );
}
