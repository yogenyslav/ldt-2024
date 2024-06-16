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
    organization?: string;
    password: string;
    roles: Role[];
    username: string;
}

export interface GetUsersInOrganizationParams {
    organization: string;
}

export interface AddUserToOrganizationParams {
    username: string;
    organization: string;
}

export interface DeleteUserParams {
    username: string;
}

export interface CreateOrganizationResponse {
    id: number;
}
