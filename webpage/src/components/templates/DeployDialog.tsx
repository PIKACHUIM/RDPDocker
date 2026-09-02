import { useEffect, useState } from 'react';
import { Server, Settings2, PackageOpen } from 'lucide-react';
import api, { Template } from '@/lib/api';

interface DeployDialogProps {
  template: Template;
  onClose: () => void;
}

interface EngineStatus {
  name: string;
  available: boolean;
  version: string;
}

const SOFTWARE_OPTIONS = [
  { id: 'firefox', name: 'Firefox 浏览器' },
  { id: 'chrome', name: 'Chrome 浏览器' },
  { id: 'vscode', name: 'VS Code' },
  { id: 'qq', name: 'QQ' },
];

const PORT_LABELS = [
  { key: 'ssh', label: 'SSH', internal: 22 },
  { key: 'rdp', label: 'RDP', internal: 3389 },
  { key: 'nx', label: 'NX', internal: 4000 },
  { key: 'vnc', label: 'VNC', internal: 5900 },
];

export default function DeployDialog({ template, onClose }: DeployDialogProps) {
  const [containerName, setContainerName] = useState('');
  const [engines, setEngines] = useState<EngineStatus[]>([]);
  const [engine, setEngine] = useState('docker');
  const [portMode, setPortMode] = useState<'auto' | 'manual'>('auto');
  const [manualPorts, setManualPorts] = useState<Record<string, string>>({
    ssh: '', rdp: '', nx: '', vnc: '',
  });
  const [softwares, setSoftwares] = useState<string[]>([]);
  const [isDeploying, setIsDeploying] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadEngines();
  }, []);

  const loadEngines = async () => {
    try {
      const data: any = await api.get('/monitor/engines');
      if (data.success) {
        const available = (data.data || []).filter((e: EngineStatus) => e.available);
        setEngines(available);
        if (available.length > 0) setEngine(available[0].name);
      }
    } catch (err) {
      console.error('Failed to load engines:', err);
    }
  };

  const toggleSoftware = (id: string) => {
    setSoftwares((prev) =>
      prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]
    );
  };

  const buildPorts = (): string[] | undefined => {
    if (portMode === 'auto') return undefined;
    const ports: string[] = [];
    for (const p of PORT_LABELS) {
      const ext = manualPorts[p.key].trim();
      if (ext) ports.push(`${ext}:${p.internal}`);
    }
    return ports.length > 0 ? ports : undefined;
  };

  const handleDeploy = async () => {
    setError('');
    if (!containerName.trim()) {
      setError('请输入容器名称');
      return;
    }
    const ports = buildPorts();

    setIsDeploying(true);
    try {
      const result: any = await api.post(`/templates/${template.id}/deploy`, {
        name: containerName.trim(),
        engine,
        ports,
        softwares: softwares.length > 0 ? softwares : undefined,
      });
      if (result.success) {
        onClose();
        window.location.href = '/containers';
      }
    } catch (err: any) {
      setError(err.response?.data?.message || '部署失败');
    } finally {
      setIsDeploying(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div className="glass-card p-6 max-w-lg w-full space-y-6 my-8">
        <div>
          <h2 className="text-2xl font-bold text-text-primary">部署容器</h2>
          <p className="text-text-secondary mt-1">{template.display_name}</p>
        </div>

        <div className="space-y-4">
          {/* Container name */}
          <div className="space-y-2">
            <label className="text-sm text-text-secondary">容器名称</label>
            <input
              type="text"
              value={containerName}
              onChange={(e) => setContainerName(e.target.value)}
              className="w-full px-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
              placeholder={`例如: my-${template.os_type}-${template.de_name}`}
            />
          </div>

          {/* Engine selection */}
          <div className="space-y-2">
            <label className="text-sm text-text-secondary flex items-center gap-2">
              <Server className="w-4 h-4" /> 容器引擎
            </label>
            <select
              value={engine}
              onChange={(e) => setEngine(e.target.value)}
              className="w-full px-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo outline-none transition-all"
            >
              {engines.length === 0 && <option value="docker">docker</option>}
              {engines.map((e) => (
                <option key={e.name} value={e.name}>
                  {e.name} ({e.version || 'N/A'})
                </option>
              ))}
            </select>
            <p className="text-xs text-text-muted">仅列出当前检测到可用的引擎</p>
          </div>

          {/* Port configuration */}
          <div className="space-y-2">
            <label className="text-sm text-text-secondary flex items-center gap-2">
              <Settings2 className="w-4 h-4" /> 端口配置
            </label>
            <div className="flex gap-2">
              <button
                onClick={() => setPortMode('auto')}
                className={`flex-1 px-3 py-2 rounded-lg text-sm border transition-colors ${
                  portMode === 'auto'
                    ? 'bg-accent-indigo/20 text-accent-indigo border-accent-indigo/30'
                    : 'bg-surface-overlay text-text-secondary border-surface-border'
                }`}
              >
                自动分配
              </button>
              <button
                onClick={() => setPortMode('manual')}
                className={`flex-1 px-3 py-2 rounded-lg text-sm border transition-colors ${
                  portMode === 'manual'
                    ? 'bg-accent-indigo/20 text-accent-indigo border-accent-indigo/30'
                    : 'bg-surface-overlay text-text-secondary border-surface-border'
                }`}
              >
                手动指定
              </button>
            </div>

            {portMode === 'manual' && (
              <div className="grid grid-cols-2 gap-2 mt-2">
                {PORT_LABELS.map((p) => (
                  <div key={p.key} className="space-y-1">
                    <label className="text-xs text-text-muted">
                      {p.label} (容器 :{p.internal})
                    </label>
                    <input
                      type="number"
                      value={manualPorts[p.key]}
                      onChange={(e) =>
                        setManualPorts((prev) => ({ ...prev, [p.key]: e.target.value }))
                      }
                      className="w-full px-3 py-1.5 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo outline-none text-sm"
                      placeholder="宿主端口"
                    />
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Software selection */}
          <div className="space-y-2">
            <label className="text-sm text-text-secondary flex items-center gap-2">
              <PackageOpen className="w-4 h-4" /> 预装软件
            </label>
            <div className="flex flex-wrap gap-2">
              {SOFTWARE_OPTIONS.map((sw) => (
                <button
                  key={sw.id}
                  onClick={() => toggleSoftware(sw.id)}
                  className={`px-3 py-1.5 rounded-lg text-sm border transition-colors ${
                    softwares.includes(sw.id)
                      ? 'bg-accent-cyan/20 text-accent-cyan border-accent-cyan/30'
                      : 'bg-surface-overlay text-text-secondary border-surface-border hover:text-text-primary'
                  }`}
                >
                  {sw.name}
                </button>
              ))}
            </div>
          </div>

          {/* Summary */}
          <div className="p-4 bg-surface-overlay rounded-lg border border-surface-border text-sm space-y-2">
            <div className="flex justify-between">
              <span className="text-text-muted">镜像</span>
              <span className="text-text-secondary">{template.image_full}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-text-muted">引擎</span>
              <span className="text-text-secondary">{engine}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-text-muted">端口</span>
              <span className="text-text-secondary">
                {portMode === 'auto' ? '自动分配' : '手动指定'}
              </span>
            </div>
          </div>

          {error && (
            <div className="text-red-500 text-sm bg-red-500/10 px-4 py-2 rounded-lg border border-red-500/30">
              {error}
            </div>
          )}
        </div>

        <div className="flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 py-2 bg-surface-overlay text-text-secondary rounded-lg hover:bg-surface-border transition-colors"
          >
            取消
          </button>
          <button
            onClick={handleDeploy}
            disabled={isDeploying}
            className="flex-1 py-2 bg-gradient-to-r from-accent-indigo to-accent-cyan text-white rounded-lg hover:shadow-lg hover:shadow-accent-indigo/30 transition-all disabled:opacity-50"
          >
            {isDeploying ? '部署中...' : '确认部署'}
          </button>
        </div>
      </div>
    </div>
  );
}
