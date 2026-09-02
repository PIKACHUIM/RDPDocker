import { useState } from 'react';
import { KeyRound } from 'lucide-react';

export default function PasswordTab() {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSubmit = async () => {
    setMessage(null);
    if (newPassword.length < 8) {
      setMessage({ type: 'error', text: '新密码至少 8 位' });
      return;
    }
    if (newPassword !== confirmPassword) {
      setMessage({ type: 'error', text: '两次输入的新密码不一致' });
      return;
    }

    setIsSubmitting(true);
    try {
      const { api } = await import('@/lib/api');
      await api.post('/auth/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      });
      setMessage({ type: 'success', text: '密码修改成功' });
      setOldPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err: any) {
      setMessage({ type: 'error', text: err.response?.data?.message || '修改失败' });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="glass-card p-6 space-y-6">
      <div className="flex items-center gap-3">
        <KeyRound className="w-5 h-5 text-accent-indigo" />
        <h2 className="text-xl font-semibold text-text-primary">修改密码</h2>
      </div>

      <div className="space-y-4 max-w-md">
        <div className="space-y-2">
          <label className="text-sm text-text-secondary">当前密码</label>
          <input
            type="password"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            className="w-full px-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
            placeholder="输入当前密码"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm text-text-secondary">新密码</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="w-full px-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
            placeholder="至少 8 位"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm text-text-secondary">确认新密码</label>
          <input
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="w-full px-4 py-2 bg-surface-overlay border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
            placeholder="再次输入新密码"
          />
        </div>

        {message && (
          <div className={`text-sm px-4 py-2 rounded-lg border ${
            message.type === 'success'
              ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/30'
              : 'bg-red-500/10 text-red-500 border-red-500/30'
          }`}>
            {message.text}
          </div>
        )}

        <button
          onClick={handleSubmit}
          disabled={isSubmitting}
          className="px-6 py-2 bg-gradient-to-r from-accent-indigo to-accent-cyan text-white rounded-lg hover:shadow-lg hover:shadow-accent-indigo/30 transition-all disabled:opacity-50"
        >
          {isSubmitting ? '提交中...' : '确认修改'}
        </button>
      </div>
    </div>
  );
}
