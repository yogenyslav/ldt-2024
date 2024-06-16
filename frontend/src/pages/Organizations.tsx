import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Dialog,
    DialogTrigger,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import {
    Table,
    TableHeader,
    TableRow,
    TableHead,
    TableBody,
    TableCell,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Info, Loader2, TrashIcon } from 'lucide-react';
import { useStores } from '@/hooks/useStores';
import { useEffect, useState } from 'react';
import { toast } from '@/components/ui/use-toast';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { LoaderButton } from '@/components/ui/loader-button';
import OrganizationsApiService from '@/api/OrganizationsApiService';
import { observer } from 'mobx-react-lite';
import { Checkbox } from '@/components/ui/checkbox';
import { Role } from '@/api/models';
import { useAuth } from '@/auth';

const Organizations = observer(() => {
    const { rootStore } = useStores();
    const [isOrganizatiionCreating, setIsOrganizationCreating] = useState(false);
    const [isUserCreating, setIsUserCreating] = useState(false);
    const [isOgranizationDialogOpen, setIsOrganizationDialogOpen] = useState(false);
    const [isUserDialogOpen, setIsUserDialogOpen] = useState(false);
    const [organizationName, setOrganizationName] = useState('');
    const [isDeledingUser, setIsDeletingUser] = useState(false);

    const auth = useAuth();

    useEffect(() => {
        rootStore
            .getOrganization()
            .then((organization) => {
                rootStore.getUsersInOrganization({ organization: organization.title }).catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось загрузить пользователей',
                        variant: 'destructive',
                    });
                });
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить организации',
                    variant: 'destructive',
                });
            });
    }, [rootStore]);

    function handleCreateOrganizationSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        console.log(event);

        setIsOrganizationCreating(true);

        const formData = new FormData(event.currentTarget);
        const name = formData.get('name') as string;

        OrganizationsApiService.createOrganization({ title: name })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Организация успешно создана',
                    variant: 'default',
                });

                rootStore.getOrganization();
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось создать организацию',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsOrganizationCreating(false);
                setIsOrganizationDialogOpen(false);
            });
    }

    function handleCreateUserSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        console.log(event);

        setIsUserCreating(true);

        const formData = new FormData(event.currentTarget);
        const firstName = formData.get('firstName') as string;
        const lastName = formData.get('lastName') as string;
        const email = formData.get('email') as string;
        const username = formData.get('username') as string;
        const password = formData.get('password') as string;
        const roleAdmin = formData.get('role-admin') as 'on' | 'off';

        console.log(firstName, lastName, email, username, password, roleAdmin);

        OrganizationsApiService.createUser({
            email,
            first_name: firstName,
            last_name: lastName,
            password,
            username,
            organization: rootStore.organization?.title,
            roles: [Role.Admin],
        })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Пользователь успешно создан',
                    variant: 'default',
                });

                if (rootStore.organization) {
                    rootStore.getUsersInOrganization({
                        organization: rootStore.organization.title,
                    });
                }
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось создать пользователя',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsUserDialogOpen(false);
                setIsUserCreating(false);
            });
    }

    const deleteUser = (username: string) => {
        setIsDeletingUser(true);

        OrganizationsApiService.deleteUser({ username })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Пользователь успешно удален',
                    variant: 'default',
                });

                if (rootStore.organization) {
                    rootStore.getUsersInOrganization({
                        organization: rootStore.organization.title,
                    });
                }
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось удалить пользователя',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsDeletingUser(false);
            });
    };

    return (
        <>
            {rootStore.isOrganizationLoading ? (
                <>
                    <Skeleton className='bg-slate-200 h-40 w-full' />
                    <Skeleton className='bg-slate-200 h-40 w-full' />
                </>
            ) : (
                <>
                    <div>
                        <div className='flex items-center'>
                            <h1 className='font-semibold text-lg md:text-2xl'>Организация</h1>
                            {!rootStore.organization && (
                                <Dialog
                                    open={isOgranizationDialogOpen}
                                    onOpenChange={setIsOrganizationDialogOpen}
                                >
                                    <DialogTrigger asChild>
                                        <Button
                                            onClick={() => {
                                                setIsOrganizationDialogOpen(true);
                                            }}
                                            className='ml-auto'
                                            size='sm'
                                        >
                                            Создать организацию
                                        </Button>
                                    </DialogTrigger>
                                    <DialogContent className='sm:max-w-[425px]'>
                                        <DialogHeader>
                                            <DialogTitle>Создать организацию</DialogTitle>
                                            <DialogDescription>
                                                Введите данные для новой организации. Название
                                                должно содержать только латинские буквы.
                                            </DialogDescription>
                                        </DialogHeader>
                                        <form onSubmit={handleCreateOrganizationSubmit}>
                                            <div className='grid gap-4 py-4'>
                                                <div className='grid items-center grid-cols-4 gap-4'>
                                                    <Label htmlFor='name' className='text-right'>
                                                        Название
                                                    </Label>
                                                    <Input
                                                        value={organizationName}
                                                        onChange={(e) => {
                                                            const value = e.target.value;

                                                            const filteredValue = value.replace(
                                                                /[^a-zA-Z]/g,
                                                                ''
                                                            );
                                                            setOrganizationName(filteredValue);
                                                        }}
                                                        required
                                                        id='name'
                                                        name='name'
                                                        placeholder='Enter organization name'
                                                        className='col-span-3'
                                                    />
                                                </div>
                                            </div>
                                            <DialogFooter>
                                                <LoaderButton
                                                    isLoading={isOrganizatiionCreating}
                                                    type='submit'
                                                >
                                                    Создать
                                                </LoaderButton>
                                            </DialogFooter>
                                        </form>
                                    </DialogContent>
                                </Dialog>
                            )}
                        </div>
                        {rootStore.organization ? (
                            <div className='border rounded-lg mt-4'>
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Название</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        <TableRow>
                                            <TableCell>
                                                <div className='font-medium'>
                                                    {rootStore.organization.title}
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    </TableBody>
                                </Table>
                            </div>
                        ) : (
                            <Alert className='mt-4'>
                                <Info className='h-4 w-4' />
                                <AlertTitle>Организация еще не создана</AlertTitle>
                                <AlertDescription>
                                    Для загрузки данных в систему и создания пользователей
                                    необходимо создать организацию
                                </AlertDescription>
                            </Alert>
                        )}
                    </div>
                    {rootStore.organization && (
                        <div>
                            <div className='flex items-center'>
                                <h1 className='font-semibold text-lg md:text-2xl'>Пользователи</h1>
                                <Dialog open={isUserDialogOpen} onOpenChange={setIsUserDialogOpen}>
                                    <DialogTrigger asChild>
                                        <Button className='ml-auto' size='sm'>
                                            Создать пользователя
                                        </Button>
                                    </DialogTrigger>
                                    <DialogContent className='sm:max-w-[425px]'>
                                        <DialogHeader>
                                            <DialogTitle>Создать пользователя</DialogTitle>
                                            <DialogDescription>
                                                Введите данные нового пользователя.
                                            </DialogDescription>
                                        </DialogHeader>
                                        <form onSubmit={handleCreateUserSubmit}>
                                            <div className='grid grid-cols-2 gap-4'>
                                                <div className='space-y-2'>
                                                    <Label htmlFor='firstName'>Имя</Label>
                                                    <Input
                                                        name='firstName'
                                                        id='firstName'
                                                        required
                                                    />
                                                </div>
                                                <div className='space-y-2'>
                                                    <Label htmlFor='lastName'>Фамилия</Label>
                                                    <Input name='lastName' id='lastName' required />
                                                </div>
                                            </div>
                                            <div className='space-y-2'>
                                                <Label htmlFor='email'>Email</Label>
                                                <Input
                                                    name='email'
                                                    id='email'
                                                    type='email'
                                                    required
                                                />
                                            </div>
                                            <div className='space-y-2'>
                                                <Label htmlFor='username'>Имя пользователя</Label>
                                                <Input name='username' id='username' required />
                                            </div>
                                            <div className='space-y-2'>
                                                <Label htmlFor='password'>Пароль</Label>
                                                <Input
                                                    name='password'
                                                    id='password'
                                                    type='password'
                                                    required
                                                />
                                            </div>
                                            <div className='space-y-2'>
                                                <Label htmlFor='roles'>Роли</Label>
                                                <div className='grid gap-2'>
                                                    <Label className='flex items-center gap-2 font-normal'>
                                                        <Checkbox
                                                            name='role-admin'
                                                            id='role-admin'
                                                        />{' '}
                                                        Администратор
                                                    </Label>
                                                    <Label className='flex items-center gap-2 font-normal'>
                                                        <Checkbox id='role-editor' /> Редактор
                                                    </Label>
                                                    <Label className='flex items-center gap-2 font-normal'>
                                                        <Checkbox id='role-viewer' /> Просмотрщик
                                                    </Label>
                                                </div>
                                            </div>
                                            <DialogFooter>
                                                <LoaderButton
                                                    isLoading={isUserCreating}
                                                    type='submit'
                                                >
                                                    Создать
                                                </LoaderButton>
                                            </DialogFooter>
                                        </form>
                                    </DialogContent>
                                </Dialog>
                            </div>
                            <div className='border shadow-sm rounded-lg mt-4'>
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Имя пользователя</TableHead>
                                            <TableHead>Организация</TableHead>
                                            <TableHead>Действия</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {rootStore.usersInOrganization.map((user) => (
                                            <TableRow>
                                                <TableCell>
                                                    <div className='font-medium'>{user}</div>
                                                </TableCell>
                                                <TableCell>
                                                    <Badge className='bg-blue-100 text-blue-900 dark:bg-blue-900 dark:text-blue-100'>
                                                        {rootStore.organization?.title}
                                                    </Badge>
                                                </TableCell>
                                                <TableCell>
                                                    <div className='flex items-center gap-2'>
                                                        <Button
                                                            variant='outline'
                                                            size='icon'
                                                            className='text-red-500'
                                                            onClick={() => {
                                                                deleteUser(user);
                                                            }}
                                                        >
                                                            {isDeledingUser ? (
                                                                <Loader2 className='h-4 w-4 animate-spin' />
                                                            ) : (
                                                                <TrashIcon className='h-4 w-4' />
                                                            )}

                                                            <span className='sr-only'>
                                                                Удалить пользователя
                                                            </span>
                                                        </Button>
                                                    </div>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        </div>
                    )}
                </>
            )}
        </>
    );
});

export default Organizations;
