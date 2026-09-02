import { useEffect, useState } from 'react';
import { Play, Square, RotateCw, Trash2, Terminal as TerminalIcon, Search } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import api, { Container } from '@/lib/api';
import { cn } from '@/lib/utils';
import ContainerDrawer from '@/components/containers/ContainerDrawer';

export default function Containers() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [filteredContainers, setFilteredContainers] = useState<Container[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState<'all' | 'running' | 'stopped'>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [detailName, setDetailName] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadContainers();
  }, []);

  useEffect(() => {
    let result = containers;
    
    // Filter by status
    if (filter === 'running') {
      result = result.filter(c => c.status === 'running');
    } else if (filter === 'stopped') {
      result = result.filter(c => c.status !== 'running');
    }
    
    // Filter by search
    if (searchQuery) {
      result = result.filter(c => 
        c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        c.image.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    setFilteredContainers(result);
  }, [containers, filter, searchQuery]);

  const loadContainers = async () => {
    try {
      const data: any = await api.get('/containers');
      if (data.success) {
        setContainers(data.data || []);
      }
    } catch (error) {
      console.error('Failed to load containers:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAction = async (name: string, action: 'start' | 'stop' | 'restart' | 'delete') => {
    try {
      if (action === 'delete') {
        if (!confirm(`确定要删除容器 ${name}？`)) return;
        await api.delete(`/containers/${name}`);
      } else {
        await api.post(`/containers/${name}/${action}`);
      }
      loadContainers();
    } catch (error) {
      console.error(`Failed to ${action} container:`, error);
    }
  };

  const handleOpenTerminal = (name: string) => {
    navigate(`/terminal?container=${name}`);
  };

  if (isLoading) {
    return <div className="text-text-secondary">加载中...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-text-primary">容器管理</h1>
          <p className="text-text-secondary mt-1">管理和监控所有容器</p>
        </div>
        <button
          onClick={() => navigate('/templates')}
          className="px-4 py-2 bg-gradient-to-r from-accent-indigo to-accent-cyan text-white rounded-lg hover:shadow-lg hover:shadow-accent-indigo/30 transition-all"
        >
          创建容器
        </button>
      </div>

      {/* Toolbar */}
      <div className="glass-card p-4 flex flex-wrap items-center gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
          <input
            type="text"
            placeholder="搜索容器名称或镜像..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
          />
        </div>
        
        <div className="flex gap-2">
          <FilterButton active={filter === 'all'} onClick={() => setFilter('all')}>
            全部 ({containers.length})
          </FilterButton>
          <FilterButton active={filter === 'running'} onClick={() => setFilter('running')}>
            运行中 ({containers.filter(c => c.status === 'running').length})
          </FilterButton>
          <FilterButton active={filter === 'stopped'} onClick={() => setFilter('stopped')}>
            已停止 ({containers.filter(c => c.status !== 'running').length})
          </FilterButton>
        </div>
      </div>

      {/* Container Grid */}
      {filteredContainers.length === 0 ? (
        <div className="glass-card p-12 text-center">
          <div className="text-text-muted mb-4">
            {searchQuery ? '未找到匹配的容器' : '暂无容器'}
          </div>
          {!searchQuery && (
            <button
              onClick={() => navigate('/templates')}
              className="px-4 py-2 bg-accent-indigo text-white rounded-lg hover:bg-accent-indigo/90 transition-colors"
            >
              前往模板市场
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          {filteredContainers.map((container) => (
            <ContainerCard
              key={container.name}
              container={container}
              onAction={handleAction}
              onOpenTerminal={handleOpenTerminal}
              onOpenDetail={() => setDetailName(container.name)}
            />
          ))}
        </div>
      )}

      {detailName && (
        <ContainerDrawer
          containerName={detailName}
          onClose={() => setDetailName(null)}
        />
      )}
    </div>
  );
}

interface ContainerCardProps {
  container: Container;
  onAction: (name: string, action: 'start' | 'stop' | 'restart' | 'delete') => void;
  onOpenTerminal: (name: string) => void;
  onOpenDetail: () => void;
}

function ContainerCard({ container, onAction, onOpenTerminal, onOpenDetail }: ContainerCardProps) {
  const isRunning = container.status === 'running';
  
  return (
    <div className="glass-card p-6 space-y-4 hover:border-accent-indigo/30 transition-colors">
      <div className="flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <div className={cn('status-dot', isRunning ? 'running' : 'stopped')} />
            <button
              onClick={onOpenDetail}
              className="font-semibold text-text-primary truncate text-left hover:text-accent-cyan transition-colors cursor-pointer"
              title="查看详情"
            >
              {container.name}
            </button>
          </div>
          <p className="text-sm text-text-muted truncate">{container.image}</p>
        </div>
      </div>

      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-text-muted">状态</span>
          <span className={cn('font-medium', isRunning ? 'text-emerald-500' : 'text-red-500')}>
            {isRunning ? '运行中' : '已停止'}
          </span>
        </div>
        {container.ports && container.ports.length > 0 && (
          <div className="flex justify-between">
            <span className="text-text-muted">端口</span>
            <span className="text-text-secondary text-right max-w-[60%]">
              {container.ports.join(', ')}
            </span>
          </div>
        )}
      </div>

      <div className="flex gap-2 pt-2 border-t border-surface-border">
        {isRunning ? (
          <>
            <ActionButton icon={<Square className="w-4 h-4" />} onClick={() => onAction(container.name, 'stop')} title="停止" />
            <ActionButton icon={<RotateCw className="w-4 h-4" />} onClick={() => onAction(container.name, 'restart')} title="重启" />
            <ActionButton icon={<TerminalIcon className="w-4 h-4" />} onClick={() => onOpenTerminal(container.name)} title="终端" className="text-accent-cyan" />
          </>
        ) : (
          <ActionButton icon={<Play className="w-4 h-4" />} onClick={() => onAction(container.name, 'start')} title="启动" className="text-emerald-500" />
        )}
        <ActionButton icon={<Trash2 className="w-4 h-4" />} onClick={() => onAction(container.name, 'delete')} title="删除" className="text-red-500 ml-auto" />
      </div>
    </div>
  );
}

function FilterButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'px-4 py-2 rounded-lg text-sm font-medium transition-all',
        active
          ? 'bg-accent-indigo text-white'
          : 'bg-surface-overlay text-text-secondary hover:text-text-primary'
      )}
    >
      {children}
    </button>
  );
}

function ActionButton({ icon, onClick, title, className }: { icon: React.ReactNode; onClick: () => void; title: string; className?: string }) {
  return (
    <button
      onClick={onClick}
      title={title}
      className={cn('p-2 rounded-lg hover:bg-surface-overlay transition-colors text-text-muted hover:text-text-primary', className)}
    >
      {icon}
    </button>
  );
}
