import { get, post, del, put } from './http';
import { SavedPrediction } from './models';

class FavoritesApiService {
    public async createFavorite(savedPrediction: SavedPrediction) {
        const response = await post<void>('/chat/favorite', savedPrediction);

        return response;
    }

    public async getFavorites() {
        const response = await get<SavedPrediction[]>('/chat/favorite/list');

        return response;
    }

    public async deleteFavorite(id: number) {
        await del(`/chat/favorite/${id}`);
    }

    public async getFavorite(id: number) {
        const response = await get<SavedPrediction>(`/chat/favorite/${id}`);

        return response;
    }

    public async updateFavorite(savedPrediction: SavedPrediction) {
        await put(`/chat/favorite`, savedPrediction);
    }
}

export default new FavoritesApiService();
