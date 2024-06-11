import axios from 'axios';

import { API_URL } from '@/config';
import { LoginParams, LoginResponse } from './models';

class AuthApiService {
    public async login(body: LoginParams) {
        const response = await axios.post<LoginResponse>(`${API_URL}/chat/auth/login`, body);

        return response.data;
    }
}

export default new AuthApiService();
