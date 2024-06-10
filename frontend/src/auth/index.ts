import { LoginParams, LoginResponse } from '@/api/models';
import React from 'react';

/**
 * This represents some generic auth provider API, like Firebase.
 */

export interface AuthContextType {
    user: LoginResponse | null;
    login: (params: LoginParams, callback: VoidFunction) => Promise<LoginResponse>;
    logout: (callback: VoidFunction) => void;
}

export const AuthContext = React.createContext<AuthContextType>(null!);

export function useAuth() {
    return React.useContext(AuthContext);
}
