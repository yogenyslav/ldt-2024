import { get, post, del } from './http';
import {
    AddUserToOrganizationParams,
    CreateOrganizationParams,
    CreateOrganizationResponse,
    CreateUserParams,
    DeleteUserParams,
    GetProductsResponse,
    GetUsersInOrganizationParams,
    Organization,
    SetUserNotificationParams,
    UserInOrganization,
} from './models/organizations';

class FavoritesApiService {
    public async createOrganization({ title }: CreateOrganizationParams) {
        const response = await post<CreateOrganizationResponse>('/admin/organization', { title });

        return response;
    }

    public async getOrganizations() {
        const response = await get<Organization[]>('/admin/organization');

        return response;
    }

    public async createUser(params: CreateUserParams) {
        const response = await post<Organization>('/admin/user', params);

        return response;
    }

    public async getUsersInOrganization({ organizationId }: GetUsersInOrganizationParams) {
        const response = await get<UserInOrganization[]>(`/admin/user/${organizationId}`);

        return response;
    }

    public async addUserToOrganization({ organization, username }: AddUserToOrganizationParams) {
        const response = await post<void>('/admin/user/organization', { organization, username });

        return response;
    }

    public async deleteUser({ organization_id, username }: DeleteUserParams) {
        const response = await del<void>(`/admin/user`, { data: { organization_id, username } });

        return response;
    }

    public async uploadFile(file: File, organization_id: number): Promise<string> {
        const formData = new FormData();
        formData.append('data', file);
        formData.append('organization_id', organization_id.toString());

        const response = await post<string>(`/admin/organization/import`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
            timeout: 510000,
        });

        return response;
    }

    public async getProducts(organization_id: number) {
        const response = await get<GetProductsResponse>(
            `/chat/stock/unique_codes/${organization_id}`
        );

        return response;
    }

    public async setUserNotifications({
        active,
        organization_id,
        username,
    }: SetUserNotificationParams) {
        const response = await post<void>('/admin/notification/switch', {
            active,
            organization_id,
            username,
        });

        return response;
    }
}

export default new FavoritesApiService();
