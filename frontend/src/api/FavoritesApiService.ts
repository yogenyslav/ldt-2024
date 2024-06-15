import { get, post, del, put } from './http';
import { SavedPrediction } from './models';

class FavoritesApiService {
    public async createFavorite(savedPrediction: SavedPrediction) {
        const response = await post<void>('/favorite', savedPrediction);

        return response;
    }

    public async getFavorites() {
        const response = await get<SavedPrediction[]>('/favorite/list');

        return response;
    }

    public async deleteFavorite(id: number) {
        await del(`/favorite/${id}`);
    }

    public async getFavorite(id: number) {
        const response = await get<SavedPrediction>(`/favorite/${id}`);

        return response;
    }

    public async updateFavorite(savedPrediction: SavedPrediction) {
        await put(`/favorite/${savedPrediction.id}`, savedPrediction);
    }
}

export default new FavoritesApiService();
