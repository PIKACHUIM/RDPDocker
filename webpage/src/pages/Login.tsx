import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { Server, Container, Package } from 'lucide-react';
import { cn } from '@/lib/utils';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(username, password);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-surface via-surface-raised to-surface flex items-center justify-center p-4">
      {/* Animated grid background */}
      <div className="absolute inset-0 bg-[linear-gradient(rgba(99,102,241,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(99,102,241,0.03)_1px,transparent_1px)] bg-[size:50px_50px]" />
      
      <div className="relative w-full max-w-5xl grid md:grid-cols-2 gap-8 items-center">
        {/* Left - Brand section */}
        <div className="hidden md:block space-y-8">
          <div>
            <h1 className="text-5xl font-bold bg-gradient-to-r from-accent-indigo to-accent-cyan bg-clip-text text-transparent mb-4">
              RDDocker
            </h1>
            <p className="text-text-secondary text-lg">
              开箱即用的容器化桌面环境管理平台
            </p>
          </div>
          
          <div className="space-y-4">
            <FeatureItem icon={<Server className="w-5 h-5" />} text="支持 Docker / Podman / LXC" />
            <FeatureItem icon={<Container className="w-5 h-5" />} text="5 大操作系统 × 9 种桌面环境" />
            <FeatureItem icon={<Package className="w-5 h-5" />} text="秒级启动，极低开销" />
          </div>
        </div>

        {/* Right - Login form */}
        <div className="glass-card p-8 space-y-6 shadow-2xl glow-border">
          <div className="space-y-2">
            <h2 className="text-3xl font-bold text-text-primary">登录</h2>
            <p className="text-text-secondary">使用您的账号登录管理平台</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm text-text-secondary">用户名</label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-4 py-3 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
                placeholder="输入用户名"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm text-text-secondary">密码</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-4 py-3 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
                placeholder="输入密码"
                required
              />
            </div>

            {error && (
              <div className="text-red-500 text-sm bg-red-500/10 px-4 py-2 rounded-lg border border-red-500/30">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className={cn(
                "w-full py-3 rounded-lg font-medium transition-all",
                "bg-gradient-to-r from-accent-indigo to-accent-cyan text-white",
                "hover:shadow-lg hover:shadow-accent-indigo/30",
                "disabled:opacity-50 disabled:cursor-not-allowed"
              )}
            >
              {isLoading ? '登录中...' : '登录'}
            </button>
          </form>

          <div className="pt-4 border-t border-surface-border text-center text-sm text-text-muted">
            首次登录？默认账号：<span className="text-accent-cyan">admin</span> / <span className="text-accent-cyan">admin123</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function FeatureItem({ icon, text }: { icon: React.ReactNode; text: string }) {
  return (
    <div className="flex items-center gap-3 text-text-secondary">
      <div className="w-10 h-10 rounded-lg bg-accent-indigo/10 flex items-center justify-center text-accent-indigo">
        {icon}
      </div>
      <span>{text}</span>
    </div>
  );
}
