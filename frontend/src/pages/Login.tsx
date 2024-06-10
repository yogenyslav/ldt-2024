import { useAuth } from '@/auth';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LoaderButton } from '@/components/ui/loader-button';
import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

export function Login() {
    const navigate = useNavigate();
    const location = useLocation();
    const auth = useAuth();
    const [loading, setLoading] = useState<boolean>(false);

    const from = location.state?.from?.pathname || '/';

    function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        setLoading(true);

        const formData = new FormData(event.currentTarget);
        const login = formData.get('login') as string;
        const password = formData.get('password') as string;

        auth.login({ username: login, password }, () => {
            navigate(from, { replace: true });
        })
            .catch(() => {})
            .finally(() => {
                setLoading(false);
            });
    }

    return (
        <div className='flex items-center h-screen'>
            <Card className='mx-auto max-w-sm min-w-96'>
                <CardHeader>
                    <CardTitle className='text-2xl'>Вход</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className='grid gap-4'>
                        <div className='grid gap-2'>
                            <Label htmlFor='login'>Логин</Label>
                            <Input id='login' name='login' required />
                        </div>
                        <div className='grid gap-2'>
                            <div className='flex items-center'>
                                <Label htmlFor='password'>Пароль</Label>
                            </div>
                            <Input id='password' name='password' type='password' required />
                        </div>
                        <LoaderButton isLoading={loading} type='submit' className='w-full'>
                            Войти
                        </LoaderButton>
                    </form>
                    <div className='mt-4 text-center text-sm'>
                        Обратитесь к администратору для получения доступа.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
