import { Navigate, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';

import type { Role } from '../api/auth';
import { useAuthStore } from '../stores/authStore';

export function RequireRole({ allow, children }: { allow: Role[]; children: ReactNode }) {
  const user = useAuthStore((state) => state.user);
  const location = useLocation();

  if (!user) {
    return <Navigate to="/gallery" state={{ from: location.pathname, loginRequired: true }} replace />;
  }
  if (!allow.includes(user.role)) {
    return <Navigate to="/gallery" replace />;
  }
  return <>{children}</>;
}
