import { LoginParams, LoginResponse } from '@/api/models';
import { AuthContext } from '.';
import { ReactNode, useEffect, useState } from 'react';
import AuthApiService from '@/api/AuthApiService';

const LOCAL_STORAGE_KEY = 'authUser';
const EXPIRATION_TIME = 24 * 60 * 60 * 1000; // 24 hours

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<LoginResponse | null>(() => {
        const storedUser = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (storedUser) {
            const { user, timestamp } = JSON.parse(storedUser);
            if (Date.now() - timestamp < EXPIRATION_TIME) {
                return user;
            } else {
                localStorage.removeItem(LOCAL_STORAGE_KEY);
            }
        }
        return null;
    });

    useEffect(() => {
        if (user) {
            localStorage.setItem(
                LOCAL_STORAGE_KEY,
                JSON.stringify({ user, timestamp: Date.now() })
            );
        } else {
            localStorage.removeItem(LOCAL_STORAGE_KEY);
        }
    }, [user]);

    const login = async (params: LoginParams, callback: VoidFunction) => {
        const response = await AuthApiService.login(params);
        setUser(response);
        callback();

        return response;
    };

    const logout = (callback: VoidFunction) => {
        setUser(null);
        localStorage.removeItem(LOCAL_STORAGE_KEY);
        callback();
    };

    return <AuthContext.Provider value={{ user, login, logout }}>{children}</AuthContext.Provider>;
}
