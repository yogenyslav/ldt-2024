import { useLocation, Navigate } from 'react-router-dom';
import { useAuth } from '.';
import { Pages } from '@/router/constants';

export function RequireUnauth({ children }: { children: JSX.Element }) {
    const auth = useAuth();
    const location = useLocation();

    if (auth.user) {
        return <Navigate to={`/${Pages.Chat}`} state={{ from: location }} replace />;
    }

    return children;
}
