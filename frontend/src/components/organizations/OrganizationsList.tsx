import { useStores } from '@/hooks/useStores';
import { Dialog, DialogTrigger } from '@radix-ui/react-dialog';
import { Button } from '../ui/button';
import {
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '../ui/dialog';
import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { LoaderButton } from '../ui/loader-button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../ui/table';
import { Alert, AlertDescription, AlertTitle } from '../ui/alert';
import { EditIcon, Info } from 'lucide-react';
import OrganizationsApiService from '@/api/OrganizationsApiService';
import { toast } from '../ui/use-toast';
import { useState } from 'react';

const OrganizationsList = () => {
    const { rootStore } = useStores();
    const [isOrganizatiionCreating, setIsOrganizationCreating] = useState(false);
    const [isOgranizationDialogOpen, setIsOrganizationDialogOpen] = useState(false);
    const [organizationName, setOrganizationName] = useState('');

    const [isEditOrganizationDialogOpen, setIsEditOrganizationDialogOpen] = useState(false);
    const [editOrganizationId, setEditOrganizationId] = useState<number | null>(null);
    const [editOrganizationName, setEditOrganizationName] = useState('');
    const [isOrganizationEditing, setIsOrganizationEditing] = useState(false);

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

                rootStore.getOrganizations();
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

    function handleEditOrganizationSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        setIsOrganizationEditing(true);

        if (editOrganizationId) {
            OrganizationsApiService.editOrganization({
                id: editOrganizationId,
                title: editOrganizationName,
            })
                .then(() => {
                    toast({
                        title: 'Успех',
                        description: 'Организация успешно отредактирована',
                        variant: 'default',
                    });

                    rootStore.getOrganizations();
                })
                .catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось отредактировать организацию',
                        variant: 'destructive',
                    });
                })
                .finally(() => {
                    setIsOrganizationEditing(false);
                    setIsEditOrganizationDialogOpen(false);
                });
        }
    }

    return (
        <>
            <div>
                <div className='flex items-center'>
                    <h1 className='font-semibold text-lg md:text-2xl'>Организация</h1>

                    <Dialog
                        open={isEditOrganizationDialogOpen}
                        onOpenChange={setIsEditOrganizationDialogOpen}
                    >
                        <DialogContent className='sm:max-w-[425px]'>
                            <DialogHeader>
                                <DialogTitle>Редактирование организации</DialogTitle>
                                <DialogDescription>
                                    Отредактируйте название организации
                                </DialogDescription>
                            </DialogHeader>
                            <form onSubmit={handleEditOrganizationSubmit}>
                                <div className='grid gap-4 py-4'>
                                    <div className='grid items-center grid-cols-4 gap-4'>
                                        <Label htmlFor='name' className='text-right'>
                                            Название
                                        </Label>
                                        <Input
                                            value={editOrganizationName}
                                            onChange={(e) => {
                                                setEditOrganizationName(e.target.value);
                                            }}
                                            required
                                            id='name'
                                            name='name'
                                            placeholder='Название организации'
                                            className='col-span-3'
                                        />
                                    </div>
                                </div>
                                <DialogFooter>
                                    <LoaderButton isLoading={isOrganizationEditing} type='submit'>
                                        Изменить
                                    </LoaderButton>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>

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
                                    Введите данные для новой организации.
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
                                                setOrganizationName(e.target.value);
                                            }}
                                            required
                                            id='name'
                                            name='name'
                                            placeholder='Название организации'
                                            className='col-span-3'
                                        />
                                    </div>
                                </div>
                                <DialogFooter>
                                    <LoaderButton isLoading={isOrganizatiionCreating} type='submit'>
                                        Создать
                                    </LoaderButton>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                </div>
                {rootStore.adminOrganizations?.length ? (
                    <div className='border rounded-lg mt-4'>
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Название</TableHead>
                                    <TableHead>Действия</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {rootStore.adminOrganizations.map((organization) => (
                                    <TableRow key={organization.id}>
                                        <TableCell>{organization.title}</TableCell>
                                        <TableCell>
                                            <Button
                                                variant='outline'
                                                size='icon'
                                                onClick={() => {
                                                    setEditOrganizationId(organization.id);
                                                    setEditOrganizationName(organization.title);
                                                    setIsEditOrganizationDialogOpen(true);
                                                }}
                                            >
                                                <EditIcon className='h-4 w-4' />
                                                <span className='sr-only'>
                                                    Удалить пользователя
                                                </span>
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </div>
                ) : (
                    <Alert className='mt-4'>
                        <Info className='h-4 w-4' />
                        <AlertTitle>Организация еще не создана</AlertTitle>
                        <AlertDescription>
                            Для загрузки данных в систему и создания пользователей необходимо
                            создать организацию
                        </AlertDescription>
                    </Alert>
                )}
            </div>
        </>
    );
};

export default OrganizationsList;
