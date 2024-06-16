import { get, post, del } from './http';
import {
    AddUserToOrganizationParams,
    CreateOrganizationParams,
    CreateOrganizationResponse,
    CreateUserParams,
    DeleteUserParams,
    GetUsersInOrganizationParams,
    Organization,
} from './models/organizations';

class FavoritesApiService {
    public async createOrganization({ title }: CreateOrganizationParams) {
        const response = await post<CreateOrganizationResponse>('/admin/organization', { title });

        return response;
    }

    public async getOrganization() {
        const response = await get<Organization>('/admin/organization');

        return response;
    }

    public async createUser(params: CreateUserParams) {
        const response = await post<Organization>('/admin/user', params);

        return response;
    }

    public async getUsersInOrganization({ organization }: GetUsersInOrganizationParams) {
        const response = await get<string[]>(`/admin/user/${organization}`);

        return response;
    }

    public async addUserToOrganization({ organization, username }: AddUserToOrganizationParams) {
        const response = await post<void>('/admin/user/organization', { organization, username });

        return response;
    }

    public async deleteUser({ username }: DeleteUserParams) {
        const response = await del<void>(`/admin/user/${username}`);

        return response;
    }

    public async uploadFile(file: File): Promise<string> {
        const formData = new FormData();
        formData.append('data', file);

        const response = await post<string>(`/admin/organization/import`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
            timeout: 510000,
        });

        return response;
    }
}

export default new FavoritesApiService();
