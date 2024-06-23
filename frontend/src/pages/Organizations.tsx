import { useStores } from '@/hooks/useStores';
import { useEffect } from 'react';
import { toast } from '@/components/ui/use-toast';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { observer } from 'mobx-react-lite';
import { Role } from '@/api/models';
import { useAuth } from '@/auth';
import OrganizationsList from '@/components/organizations/OrganizationsList';
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import UsersList from '@/components/organizations/UsersList';
import OrganizationDataUpload from '@/components/organizations/OrganizationDataUpload';
import GoodsList from '@/components/organizations/GoodsList';

const Organizations = observer(() => {
    const { rootStore } = useStores();

    const auth = useAuth();

    const roles = auth.user?.roles;

    useEffect(() => {
        if (roles?.includes(Role.Admin)) {
            rootStore
                .getOrganizations()
                .then((organizations) => {
                    if (organizations.length) {
                        rootStore.setSelectedAdminOrganization(organizations[0].id);
                    }
                })
                .catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось загрузить организации',
                        variant: 'destructive',
                    });
                });
        }
    }, [rootStore, roles]);

    return (
        <div className='organizations'>
            {roles?.includes(Role.Admin) ? (
                rootStore.isOrganizationLoading ? (
                    <>
                        <Skeleton className='bg-slate-200 h-40 w-full' />
                        <Skeleton className='bg-slate-200 h-40 w-full' />
                    </>
                ) : (
                    <div>
                        <OrganizationsList />

                        <div className='mt-8'>
                            <div className='flex items-center'>
                                <h1 className='font-semibold text-lg md:text-2xl'>
                                    Выбор организации
                                </h1>
                            </div>

                            <div className='mt-3'>
                                <Select
                                    onValueChange={(organizationId) =>
                                        rootStore.setSelectedAdminOrganization(+organizationId)
                                    }
                                >
                                    <SelectTrigger className='w-[280px] ml-2'>
                                        <SelectValue placeholder='Выберите организацию' />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectLabel>Организации</SelectLabel>
                                            {rootStore.adminOrganizations?.map((organization) => (
                                                <SelectItem
                                                    key={organization.id}
                                                    value={`${organization.id}`}
                                                >
                                                    {organization.title}
                                                </SelectItem>
                                            ))}
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>

                        {rootStore.selectedAdminOrganization && (
                            <>
                                <UsersList
                                    organizationId={rootStore.selectedAdminOrganization?.id}
                                />

                                <OrganizationDataUpload
                                    organizationId={rootStore.selectedAdminOrganization?.id}
                                />

                                <GoodsList
                                    organizationId={rootStore.selectedAdminOrganization?.id}
                                />
                            </>
                        )}
                    </div>
                )
            ) : (
                <Alert>
                    <AlertTitle>У вас нет доступа</AlertTitle>
                    <AlertDescription>
                        Для доступа к этой странице необходимо быть администратором
                    </AlertDescription>
                </Alert>
            )}
        </div>
    );
});

export default Organizations;
