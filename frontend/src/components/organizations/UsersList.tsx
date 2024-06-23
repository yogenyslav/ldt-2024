import { Loader2, TrashIcon } from 'lucide-react';
import { Button } from '../ui/button';
import { Checkbox } from '../ui/checkbox';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '../ui/dialog';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { LoaderButton } from '../ui/loader-button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../ui/table';
import { useEffect, useState } from 'react';
import { useStores } from '@/hooks/useStores';
import { toast } from '../ui/use-toast';
import OrganizationsApiService from '@/api/OrganizationsApiService';
import { Role } from '@/api/models';
import { Badge } from '../ui/badge';
import { observer } from 'mobx-react-lite';

type Props = {
    organizationId: number;
};

const UsersList = observer(({ organizationId }: Props) => {
    const { rootStore } = useStores();
    const [isUserCreating, setIsUserCreating] = useState(false);
    const [isUserDialogOpen, setIsUserDialogOpen] = useState(false);
    const [isDeledingUser, setIsDeletingUser] = useState(false);

    useEffect(() => {
        rootStore.getUsersInOrganization({ organizationId }).catch(() => {
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить пользователей',
                variant: 'destructive',
            });
        });
    }, [rootStore, organizationId]);

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
        const roleAdmin = formData.get('role-admin') as 'on' | null;
        const roleAnalyst = formData.get('analyst') as 'on' | null;
        const roleBuyer = formData.get('buyer') as 'on' | null;

        const roles = [
            roleAdmin ? Role.Admin : null,
            roleAnalyst ? Role.Analyst : null,
            roleBuyer ? Role.Buyer : null,
        ].filter((role) => role !== null) as Role[];

        console.log(firstName, lastName, email, username, password, roleAdmin);

        OrganizationsApiService.createUser({
            email,
            first_name: firstName,
            last_name: lastName,
            password,
            username,
            organization_ids: [organizationId],
            roles,
        })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Пользователь успешно создан',
                    variant: 'default',
                });

                if (rootStore.adminOrganizations) {
                    rootStore.getUsersInOrganization({
                        organizationId,
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

        OrganizationsApiService.deleteUser({ username, organization_id: organizationId })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Пользователь успешно удален',
                    variant: 'default',
                });

                if (rootStore.adminOrganizations) {
                    rootStore.getUsersInOrganization({
                        organizationId,
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
        <div className='mt-8'>
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
                                    <Input name='firstName' id='firstName' required />
                                </div>
                                <div className='space-y-2'>
                                    <Label htmlFor='lastName'>Фамилия</Label>
                                    <Input name='lastName' id='lastName' required />
                                </div>
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='email'>Email</Label>
                                <Input name='email' id='email' type='email' required />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='username'>Имя пользователя</Label>
                                <Input name='username' id='username' required />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='password'>Пароль</Label>
                                <Input name='password' id='password' type='password' required />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='roles'>Роли</Label>
                                <div className='grid gap-2'>
                                    <Label className='flex items-center gap-2 font-normal'>
                                        <Checkbox name='role-admin' id='role-admin' /> Администратор
                                    </Label>
                                    <Label className='flex items-center gap-2 font-normal'>
                                        <Checkbox name='analyst' id='analyst' /> Аналитик
                                    </Label>
                                    <Label className='flex items-center gap-2 font-normal'>
                                        <Checkbox name='buyer' id='buyer' /> Закупщик
                                    </Label>
                                </div>
                            </div>
                            <DialogFooter>
                                <LoaderButton isLoading={isUserCreating} type='submit'>
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
                            <TableHead>E-mail</TableHead>
                            <TableHead>Организация</TableHead>
                            <TableHead>Действия</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {rootStore.usersInOrganization.map((user, index) => (
                            <TableRow key={index}>
                                <TableCell>
                                    <div className='font-medium'>{user.username}</div>
                                </TableCell>
                                <TableCell>
                                    <div className='font-medium'>{user.email}</div>
                                </TableCell>
                                <TableCell>
                                    <Badge className='bg-blue-100 text-blue-900 dark:bg-blue-900 dark:text-blue-100'>
                                        {
                                            rootStore.adminOrganizations?.find(
                                                (organization) => organization.id === organizationId
                                            )?.title
                                        }
                                    </Badge>
                                </TableCell>
                                <TableCell>
                                    <div className='flex items-center gap-2'>
                                        <Button
                                            variant='outline'
                                            size='icon'
                                            className='text-red-500'
                                            onClick={() => {
                                                deleteUser(user.username);
                                            }}
                                        >
                                            {isDeledingUser ? (
                                                <Loader2 className='h-4 w-4 animate-spin' />
                                            ) : (
                                                <TrashIcon className='h-4 w-4' />
                                            )}

                                            <span className='sr-only'>Удалить пользователя</span>
                                        </Button>
                                    </div>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>
        </div>
    );
});

export default UsersList;
