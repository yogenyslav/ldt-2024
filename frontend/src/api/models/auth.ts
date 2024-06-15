export interface LoginParams {
    username: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    roles: Role[];
}

export enum Role {
    Admin = 'ADMIN',
    Analyst = 'ANALYST',
    Buyer = 'BUYER',
}
