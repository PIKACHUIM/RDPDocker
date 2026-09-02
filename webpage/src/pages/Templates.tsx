import { useEffect, useState } from 'react';
import { Package, Cpu, HardDrive, Monitor } from 'lucide-react';
import api, { Template } from '@/lib/api';
import { cn } from '@/lib/utils';
import DeployDialog from '@/components/templates/DeployDialog';

export default function Templates() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [filteredTemplates, setFilteredTemplates] = useState<Template[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [category, setCategory] = useState<string>('all');
  const [deployDialogOpen, setDeployDialogOpen] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);

  useEffect(() => {
    loadTemplates();
  }, []);

  useEffect(() => {
    if (category === 'all') {
      setFilteredTemplates(templates);
    } else {
      setFilteredTemplates(templates.filter(t => t.de_name === category));
    }
  }, [templates, category]);

  const loadTemplates = async () => {
    try {
      const data: any = await api.get('/templates');
      if (data.success) {
        setTemplates(data.data || []);
      }
    } catch (error) {
      console.error('Failed to load templates:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeploy = (template: Template) => {
    setSelectedTemplate(template);
    setDeployDialogOpen(true);
  };

  const categories = [
    { id: 'all', name: '全部', icon: <Package /> },
    { id: 'server', name: 'Server', icon: <HardDrive /> },
    { id: 'gnome3', name: 'GNOME', icon: <Monitor /> },
    { id: 'xfce4l', name: 'Xfce', icon: <Monitor /> },
    { id: 'plasma', name: 'Plasma', icon: <Monitor /> },
    { id: 'deepin', name: 'Deepin', icon: <Monitor /> },
  ];

  if (isLoading) {
    return <div className="text-text-secondary">加载中...</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-text-primary">模板市场</h1>
        <p className="text-text-secondary mt-1">选择预配置的操作系统和桌面环境</p>
      </div>

      {/* Category Tabs */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => setCategory(cat.id)}
            className={cn(
              'flex items-center gap-2 px-4 py-2 rounded-lg whitespace-nowrap transition-all',
              category === cat.id
                ? 'bg-accent-indigo text-white'
                : 'bg-surface-raised text-text-secondary hover:text-text-primary hover:bg-surface-overlay'
            )}
          >
            <span className="w-4 h-4">{cat.icon}</span>
            <span className="font-medium">{cat.name}</span>
          </button>
        ))}
      </div>

      {/* Template Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredTemplates.map((template) => (
          <TemplateCard key={template.id} template={template} onDeploy={handleDeploy} />
        ))}
      </div>

      {deployDialogOpen && selectedTemplate && (
        <DeployDialog
          template={selectedTemplate}
          onClose={() => {
            setDeployDialogOpen(false);
            setSelectedTemplate(null);
          }}
        />
      )}
    </div>
  );
}

interface TemplateCardProps {
  template: Template;
  onDeploy: (template: Template) => void;
}

function TemplateCard({ template, onDeploy }: TemplateCardProps) {
  const getOSIcon = (osType: string) => {
    const icons: Record<string, string> = {
      debian: '🌀',
      ubuntu: '🟠',
      alpine: '🏔️',
      fedora: '🎩',
      arch: '🏛️',
    };
    return icons[osType] || '🐧';
  };

  return (
    <div className="glass-card p-6 space-y-4 hover:border-accent-indigo/30 transition-all group">
      <div className="flex items-start justify-between">
        <div className="text-4xl">{getOSIcon(template.os_type)}</div>
        <span className="px-2 py-1 rounded text-xs bg-accent-indigo/20 text-accent-indigo border border-accent-indigo/30">
          {template.category || 'desktop'}
        </span>
      </div>

      <div>
        <h3 className="font-semibold text-text-primary mb-1">{template.display_name}</h3>
        <p className="text-sm text-text-muted line-clamp-2">{template.description}</p>
      </div>

      <div className="flex items-center gap-2 text-xs text-text-secondary">
        <Cpu className="w-3 h-3" />
        <span>{template.os_type} {template.os_version}</span>
        <span className="text-text-muted">•</span>
        <span>{template.de_name}</span>
      </div>

      <button
        onClick={() => onDeploy(template)}
        className="w-full py-2 bg-gradient-to-r from-accent-indigo to-accent-cyan text-white rounded-lg opacity-0 group-hover:opacity-100 hover:shadow-lg hover:shadow-accent-indigo/30 transition-all"
      >
        一键部署
      </button>
    </div>
  );
}


