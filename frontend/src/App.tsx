import { Route, Routes } from 'react-router-dom';
import { Login } from './pages/Login';
import Chat from './pages/Chat';
import { RequireAuth } from './auth/RequireAuth';
import { AuthProvider } from './auth/AuthProvider';
import { RequireUnauth } from './auth/RequireUnauth';
import { Dashboard } from './components/Dashboard';

function App() {
    return (
        <>
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
                        path='/chat'
                        element={
                            <RequireAuth>
                                <Dashboard>
                                    <Chat />
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
