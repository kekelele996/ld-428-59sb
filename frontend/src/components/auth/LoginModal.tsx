import { useState } from 'react';

import { useAuthStore } from '../../stores/authStore';
import { ApiError } from '../../utils/request';

export function LoginModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const login = useAuthStore((state) => state.login);
  const [username, setUsername] = useState('curator');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (!open) return null;

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(username, password);
      onClose();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '登录失败，请稍后再试。');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-ink/50 p-4" onClick={onClose}>
      <form
        className="w-full max-w-sm border border-ink/15 bg-rice p-7 shadow-xl"
        onClick={(event) => event.stopPropagation()}
        onSubmit={submit}
      >
        <h2 className="font-display text-3xl">策展人登录</h2>
        <p className="mt-2 text-xs leading-5 text-ink/50">演示账号 curator / Curator@123，登录后可管理场次与预约记录。</p>
        <label className="mt-5 block text-xs font-semibold uppercase tracking-wider text-ink/55">
          用户名
          <input
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            className="mt-1 w-full border border-ink/20 bg-white px-4 py-3 text-sm normal-case tracking-normal outline-none focus:border-clay"
          />
        </label>
        <label className="mt-4 block text-xs font-semibold uppercase tracking-wider text-ink/55">
          密码
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            className="mt-1 w-full border border-ink/20 bg-white px-4 py-3 text-sm normal-case tracking-normal outline-none focus:border-clay"
          />
        </label>
        {error && <p className="mt-3 text-sm text-clay">{error}</p>}
        <div className="mt-6 flex gap-3">
          <button type="submit" disabled={loading} className="flex-1 bg-ink px-5 py-3 text-sm font-semibold text-rice hover:bg-clay disabled:opacity-40">
            {loading ? '登录中…' : '登录'}
          </button>
          <button type="button" onClick={onClose} className="border border-ink/20 px-5 py-3 text-sm hover:border-ink">
            取消
          </button>
        </div>
      </form>
    </div>
  );
}
