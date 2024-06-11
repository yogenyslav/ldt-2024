import { CircleUser, Menu, Package2, Search, SquarePen } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '@/auth';
import SessionsHistory from './SessionsHistory';

const menuItems = [
    {
        icon: SquarePen,
        title: 'Новый чат',
        to: '/',
    },
];

type DashboardProps = {
    children: React.ReactNode;
};

export function Dashboard({ children }: DashboardProps) {
    const navigate = useNavigate();
    const auth = useAuth();
    const location = useLocation();

    const from = location.state?.from?.pathname || '/';

    return (
        <>
            <div className='grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]'>
                <div className='hidden border-r bg-muted/40 md:block'>
                    <div className='flex h-full max-h-screen flex-col gap-2'>
                        <div className='flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6'>
                            <Link to='/' className='flex items-center gap-2 font-semibold'>
                                <Package2 className='h-6 w-6' />
                                <span className=''>misis.tech</span>
                            </Link>
                        </div>
                        <div className='flex-1'>
                            <nav className='grid items-start px-2 text-sm font-medium lg:px-4'>
                                <Navigation />
                            </nav>

                            <SessionsHistory />
                        </div>
                    </div>
                </div>
                <div className='flex flex-col'>
                    <header className='flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6'>
                        <Sheet>
                            <SheetTrigger asChild>
                                <Button
                                    variant='outline'
                                    size='icon'
                                    className='shrink-0 md:hidden'
                                >
                                    <Menu className='h-5 w-5' />
                                    <span className='sr-only'>Открыть меню</span>
                                </Button>
                            </SheetTrigger>
                            <SheetContent side='left' className='flex flex-col'>
                                <nav className='grid gap-2 text-lg font-medium'>
                                    <Link
                                        to='/'
                                        className='flex items-center gap-2 text-lg font-semibold'
                                    >
                                        <Package2 className='h-6 w-6' />
                                        misis.tech
                                    </Link>

                                    <SessionsHistory />
                                </nav>

                                <Navigation />
                            </SheetContent>
                        </Sheet>
                        <div className='w-full flex-1'>
                            <form>
                                <div className='relative'>
                                    <Search className='absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground' />
                                    <Input
                                        type='search'
                                        placeholder='Поиск...'
                                        className='w-full appearance-none bg-background pl-8 shadow-none md:w-2/3 lg:w-1/3'
                                    />
                                </div>
                            </form>
                        </div>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant='secondary' size='icon' className='rounded-full'>
                                    <CircleUser className='h-5 w-5' />
                                    <span className='sr-only'>Toggle user menu</span>
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align='end'>
                                <DropdownMenuLabel>Мой аккаунт</DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem
                                    onClick={() =>
                                        auth.logout(() => {
                                            navigate(from, { replace: true });
                                        })
                                    }
                                >
                                    Выйти
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </header>
                    <main className='flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6'>
                        {children}
                    </main>
                </div>
            </div>
        </>
    );
}

const Navigation = () => {
    return (
        <>
            {menuItems.map((item, index) => (
                <Link
                    key={index}
                    to={item.to}
                    className='flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary'
                >
                    <item.icon className='h-4 w-4' />
                    {item.title}
                </Link>
            ))}
        </>
    );
};
