import FavoritesApiService from '@/api/FavoritesApiService';
import { SavedPrediction } from '@/api/models';
import SavedOutputJson from '@/components/SavedOutputJson';
import { toast } from '@/components/ui/use-toast';
import { useEffect, useState } from 'react';

const SavedPredictions = () => {
    const [savedPredictions, setSavedPredictions] = useState<SavedPrediction[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        setIsLoading(true);

        FavoritesApiService.getFavorites()
            .then((savedPrediction) => {
                setSavedPredictions(savedPrediction);
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить сохраненные планы закупок',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsLoading(false);
            });
    }, []);

    return (
        <div className='saved flex gap-4 flex-col pb-10'>
            {isLoading ? (
                <div>Loading...</div>
            ) : savedPredictions.length ? (
                savedPredictions.map((savedPrediction) => (
                    <SavedOutputJson
                        key={savedPrediction.id}
                        outputJson={savedPrediction.response}
                        id={savedPrediction.id}
                    />
                ))
            ) : (
                <div>Нет сохраненных планов закупок</div>
            )}
        </div>
    );
};

export default SavedPredictions;
