import { useEffect, useState } from 'react';
import { X, Box, Hash, Activity, Network } from 'lucide-react';
import api from '@/lib/api';
import { cn } from '@/lib/utils';

interface ContainerDetail {
  name: string;
  status: string;
  image: string;
  ports: string[];
  id: string;
}

interface ContainerDrawerProps {
  containerName: string;
  onClose: () => void;
}

export default function ContainerDrawer({ containerName, onClose }: ContainerDrawerProps) {
  const [detail, setDetail] = useState<ContainerDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDetail();
  }, [containerName]);

  const loadDetail = async () => {
    setIsLoading(true);
    setError('');
    try {
      const data: any = await api.get(`/containers/${containerName}`);
      if (data.success) {
        setDetail(data.data);
      } else {
        setError(data.message || '加载失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || '加载失败');
    } finally {
      setIsLoading(false);
    }
  };

  const isRunning = detail?.status === 'running';

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />

      {/* Drawer */}
      <div className="relative w-full max-w-md h-full bg-surface-raised border-l border-surface-border shadow-2xl overflow-y-auto">
        <div className="sticky top-0 bg-surface-raised/90 backdrop-blur-md border-b border-surface-border p-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-text-primary">容器详情</h2>
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-surface-overlay text-text-muted hover:text-text-primary transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 space-y-6">
          {isLoading ? (
            <div className="text-text-secondary text-sm">加载中...</div>
          ) : error ? (
            <div className="text-red-500 text-sm bg-red-500/10 px-4 py-2 rounded-lg border border-red-500/30">
              {error}
            </div>
          ) : detail ? (
            <>
              {/* Header */}
              <div className="space-y-2">
                <div className="flex items-center gap-3">
                  <div className={cn('status-dot', isRunning ? 'running' : 'stopped')} />
                  <span className="text-2xl font-bold text-text-primary">{detail.name}</span>
                </div>
                <span className={`inline-block px-3 py-1 rounded-full text-sm border ${
                  isRunning
                    ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/30'
                    : 'bg-red-500/10 text-red-500 border-red-500/30'
                }`}>
                  {isRunning ? '运行中' : '已停止'}
                </span>
              </div>

              {/* Details */}
              <div className="space-y-4">
                <DetailRow icon={<Box className="w-4 h-4" />} label="镜像" value={detail.image} />
                <DetailRow
                  icon={<Hash className="w-4 h-4" />}
                  label="容器 ID"
                  value={detail.id ? detail.id.slice(0, 12) : 'N/A'}
                  mono
                />
                <DetailRow
                  icon={<Activity className="w-4 h-4" />}
                  label="状态"
                  value={detail.status}
                />
              </div>

              {/* Ports */}
              <div className="space-y-2">
                <label className="text-sm text-text-secondary flex items-center gap-2">
                  <Network className="w-4 h-4" /> 端口映射
                </label>
                {detail.ports && detail.ports.length > 0 ? (
                  <div className="space-y-2">
                    {detail.ports.map((port, i) => (
                      <div
                        key={i}
                        className="px-3 py-2 bg-surface-overlay rounded-lg border border-surface-border text-sm font-mono text-text-secondary"
                      >
                        {port}
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-sm text-text-muted">无端口映射</div>
                )}
              </div>
            </>
          ) : null}
        </div>
      </div>
    </div>
  );
}

function DetailRow({ icon, label, value, mono }: { icon: React.ReactNode; label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between py-2 border-b border-surface-border">
      <span className="text-sm text-text-muted flex items-center gap-2">{icon}{label}</span>
      <span className={cn('text-sm text-text-primary', mono && 'font-mono')}>{value}</span>
    </div>
  );
}
