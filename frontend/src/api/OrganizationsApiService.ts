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
}

export default new FavoritesApiService();
