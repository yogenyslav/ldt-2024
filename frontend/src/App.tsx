import { Route, Routes } from 'react-router-dom';
import { Login } from './pages/Login';
import Chat from './pages/Chat';
import { RequireAuth } from './auth/RequireAuth';
import { AuthProvider } from './auth/AuthProvider';
import { RequireUnauth } from './auth/RequireUnauth';
import { Dashboard } from './components/Dashboard';
import { Toaster } from './components/ui/toaster';
import SavedPredictions from './pages/SavedPredictions';

function App() {
    return (
        <>
            <Toaster />

            <AuthProvider>
                <Routes>
                    <Route
                        path='/login'
                        element={
                            <RequireUnauth>
                                <Login />
                            </RequireUnauth>
                        }
                    />
                    <Route
                        path='/chat/:sessionId'
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <Chat />
                                </Dashboard>
                            </RequireAuth>
                        }
                    />
                    <Route
                        path='/saved'
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <SavedPredictions />
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
