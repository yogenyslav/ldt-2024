import { Route, Routes } from 'react-router-dom';
import { Login } from './pages/Login';
import Chat from './pages/Chat';
import { RequireAuth } from './auth/RequireAuth';
import { AuthProvider } from './auth/AuthProvider';
import { RequireUnauth } from './auth/RequireUnauth';

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
                                <Chat />
                            </RequireAuth>
                        }
                    />
                    <Route
                        path='*'
                        element={
                            <RequireAuth>
                                <Chat />
                            </RequireAuth>
                        }
                    />
                </Routes>
            </AuthProvider>
        </>
    );
}

export default App;
