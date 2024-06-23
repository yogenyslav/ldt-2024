import { Role } from './auth';

export interface Organization {
    id: number;
    title: string;
    created_at: string;
}

export interface CreateOrganizationParams {
    title: string;
}

export interface CreateUserParams {
    email: string;
    first_name: string;
    last_name: string;
    organization_ids: number[];
    password: string;
    roles: Role[];
    username: string;
}

export interface GetUsersInOrganizationParams {
    organizationId: number;
}

export interface AddUserToOrganizationParams {
    username: string;
    organization: string;
}

export interface DeleteUserParams {
    organization_id: number;
    username: string;
}

export interface CreateOrganizationResponse {
    id: number;
}

export interface Product {
    name: string;
    regular: boolean;
    segment: string;
}

export interface GetProductsResponse {
    codes: Product[];
}

export interface UserInOrganization {
    username: string;
    notifications: boolean;
}
