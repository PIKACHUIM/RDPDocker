import { useEffect, useState } from 'react';
import { UserPlus, Shield, User as UserIcon } from 'lucide-react';
import api from '@/lib/api';

interface AdminUser {
  id: number;
  username: string;
  role: string;
  enabled: boolean;
  created_at: string;
}

export default function UsersTab() {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newRole, setNewRole] = useState<'admin' | 'user'>('user');
  const [error, setError] = useState('');

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const data: any = await api.get('/users');
      if (data.success) setUsers(data.data || []);
    } catch (err) {
      console.error('Failed to load users:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreate = async () => {
    setError('');
    if (newUsername.trim().length < 3) {
      setError('用户名至少 3 个字符');
      return;
    }
    if (newPassword.length < 8) {
      setError('密码至少 8 位');
      return;
    }
    try {
      await api.post('/users', {
        username: newUsername.trim(),
        password: newPassword,
        role: newRole,
      });
      setNewUsername('');
      setNewPassword('');
      setNewRole('user');
      setShowForm(false);
      loadUsers();
    } catch (err: any) {
      setError(err.response?.data?.message || '创建失败');
    }
  };

  const handleToggle = async (user: AdminUser) => {
    try {
      await api.put(`/users/${user.id}/enabled`, { enabled: !user.enabled });
      loadUsers();
    } catch (err) {
      console.error('Failed to toggle user:', err);
    }
  };

  if (isLoading) {
    return <div className="glass-card p-6 text-text-secondary">加载中...</div>;
  }

  return (
    <div className="glass-card p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-text-primary">用户管理</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex items-center gap-2 px-4 py-2 bg-accent-indigo text-white rounded-lg hover:bg-accent-indigo/90 transition-colors"
        >
          <UserPlus className="w-4 h-4" />
          添加用户
        </button>
      </div>

      {showForm && (
        <div className="p-4 bg-surface-overlay rounded-lg border border-surface-border space-y-4">
          <div className="space-y-2">
            <label className="text-sm text-text-secondary">用户名</label>
            <input
              type="text"
              value={newUsername}
              onChange={(e) => setNewUsername(e.target.value)}
              className="w-full px-4 py-2 bg-surface border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
              placeholder="3-32 个字符"
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm text-text-secondary">密码</label>
            <input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="w-full px-4 py-2 bg-surface border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
              placeholder="至少 8 位"
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm text-text-secondary">角色</label>
            <select
              value={newRole}
              onChange={(e) => setNewRole(e.target.value as 'admin' | 'user')}
              className="w-full px-4 py-2 bg-surface border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo outline-none transition-all"
            >
              <option value="user">user（常规操作）</option>
              <option value="admin">admin（含配置与用户管理）</option>
            </select>
          </div>

          {error && (
            <div className="text-red-500 text-sm bg-red-500/10 px-4 py-2 rounded-lg border border-red-500/30">
              {error}
            </div>
          )}

          <div className="flex gap-3">
            <button
              onClick={handleCreate}
              className="px-4 py-2 bg-gradient-to-r from-accent-indigo to-accent-cyan text-white rounded-lg hover:shadow-lg hover:shadow-accent-indigo/30 transition-all"
            >
              创建
            </button>
            <button
              onClick={() => setShowForm(false)}
              className="px-4 py-2 bg-surface-overlay text-text-secondary rounded-lg hover:bg-surface-border transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      )}

      <div className="space-y-3">
        {users.map((user) => (
          <div
            key={user.id}
            className="flex items-center justify-between p-4 bg-surface-overlay rounded-lg border border-surface-border"
          >
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-accent-indigo/10 flex items-center justify-center text-accent-indigo">
                {user.role === 'admin' ? <Shield className="w-5 h-5" /> : <UserIcon className="w-5 h-5" />}
              </div>
              <div>
                <div className="font-medium text-text-primary flex items-center gap-2">
                  {user.username}
                  <span className={`px-2 py-0.5 rounded text-xs border ${
                    user.role === 'admin'
                      ? 'bg-accent-indigo/20 text-accent-indigo border-accent-indigo/30'
                      : 'bg-surface-border text-text-secondary border-surface-border'
                  }`}>
                    {user.role}
                  </span>
                </div>
                <div className="text-xs text-text-muted mt-1">
                  创建于 {new Date(user.created_at).toLocaleString()}
                </div>
              </div>
            </div>

            <button
              onClick={() => handleToggle(user)}
              className={`px-3 py-1 rounded-lg text-sm border transition-colors ${
                user.enabled
                  ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/30 hover:bg-emerald-500/20'
                  : 'bg-red-500/10 text-red-500 border-red-500/30 hover:bg-red-500/20'
              }`}
            >
              {user.enabled ? '已启用' : '已禁用'}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
