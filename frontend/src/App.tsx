import { Route, Routes } from 'react-router-dom';
import { Login } from './pages/Login';
import Chat from './pages/Chat';
import { RequireAuth } from './auth/RequireAuth';
import { AuthProvider } from './auth/AuthProvider';
import { RequireUnauth } from './auth/RequireUnauth';
import { Dashboard } from './components/Dashboard';
import { Toaster } from './components/ui/toaster';
import SavedPredictions from './pages/SavedPredictions';
import Organizations from './pages/Organizations';
import { Pages } from './router/constants';
import { useEffect } from 'react';

function App() {
    useEffect(() => {
        // eslint-disable-next-line
        // @ts-ignore
        if (window.Telegram && window.Telegram.WebApp) {
            // eslint-disable-next-line
            // @ts-ignore
            const telegramWebApp = window.Telegram.WebApp;

            telegramWebApp.expand();
        } else {
            console.error('Telegram WebApp is not available');
        }
    }, []);

    return (
        <>
            <Toaster />

            <AuthProvider>
                <Routes>
                    <Route
                        path={`/${Pages.Login}`}
                        element={
                            <RequireUnauth>
                                <Login />
                            </RequireUnauth>
                        }
                    />
                    <Route
                        path={`/${Pages.Chat}/:sessionId`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <Chat />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path={`/${Pages.SavedPredictions}`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <SavedPredictions />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path={`/${Pages.Organizatinos}`}
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <Organizations />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path='*'
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <Chat />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                </Routes>
            </AuthProvider>
        </>
    );
}

export default App;
