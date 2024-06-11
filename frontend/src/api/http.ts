import axiosInstance from './axiosInstance';

export const get = async <T>(url: string, config = {}) => {
    const response = await axiosInstance.get<T>(url, config);
    return response.data;
};

export const post = async <T>(url: string, data: unknown, config = {}) => {
    const response = await axiosInstance.post<T>(url, data, config);
    return response.data;
};

export const put = async <T>(url: string, data: unknown, config = {}) => {
    const response = await axiosInstance.put<T>(url, data, config);
    return response.data;
};

export const del = async <T>(url: string, config = {}) => {
    const response = await axiosInstance.delete<T>(url, config);
    return response.data;
};
