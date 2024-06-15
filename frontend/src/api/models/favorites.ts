import { PredictionResponse } from './predict';

export interface SavedPrediction {
    id: number;
    response: PredictionResponse;
}
