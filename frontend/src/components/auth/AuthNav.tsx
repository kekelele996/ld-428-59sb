import { useEffect, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { LoginModal } from './LoginModal';
import { useAuthStore } from '../../stores/authStore';

// AuthNav 顶栏的登录态与策展人入口；路由守卫重定向时自动弹出登录框。
export function AuthNav() {
  const [loginOpen, setLoginOpen] = useState(false);
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if ((location.state as { loginRequired?: boolean } | null)?.loginRequired) {
      setLoginOpen(true);
      navigate(location.pathname, { replace: true, state: null });
    }
  }, [location, navigate]);

  const isCurator = user?.role === 'Curator' || user?.role === 'Admin';

  return (
    <>
      <nav className="flex items-center gap-5 text-sm text-ink/70">
        <Link to="/gallery">画廊</Link>
        <Link to="/studio">工作台</Link>
        {isCurator && <Link to="/curator/reservations" className="font-semibold text-clay">预约管理</Link>}
        {user ? (
          <>
            <span className="text-ink/45">{user.name}（{user.role}）</span>
            <button onClick={logout} className="text-ink/50 underline-offset-2 hover:underline">退出</button>
          </>
        ) : (
          <button onClick={() => setLoginOpen(true)} className="font-semibold text-clay">策展人登录</button>
        )}
      </nav>
      <LoginModal open={loginOpen} onClose={() => setLoginOpen(false)} />
    </>
  );
}
