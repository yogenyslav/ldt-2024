import { observer } from 'mobx-react-lite';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { useStores } from '@/hooks/useStores';

const OrganizationSwitcher = observer(() => {
    const { rootStore } = useStores();

    return (
        <>
            {rootStore.adminOrganizations?.length && rootStore.selectedOrganizationId && (
                <Select
                    defaultValue={`${rootStore.selectedOrganizationId}`}
                    onValueChange={(value) => rootStore.setSelectedOrganizationId(+value)}
                >
                    <SelectTrigger
                        className='flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0'
                        aria-label='Select account'
                    >
                        <SelectValue placeholder='Select an account'>
                            {
                                rootStore.adminOrganizations.find(
                                    (organization) =>
                                        organization.id === rootStore.selectedOrganizationId
                                )?.title
                            }
                        </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                        {rootStore.adminOrganizations.map((organization) => (
                            <SelectItem key={organization.id} value={`${organization.id}`}>
                                <div className='flex items-center gap-3 [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0 [&_svg]:text-foreground'>
                                    {organization.title}
                                </div>
                            </SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            )}
        </>
    );
});

export default OrganizationSwitcher;
