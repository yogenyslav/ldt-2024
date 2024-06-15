export interface SavedPrediction {
    id: number;
    response: OutputJson;
}

export interface OutputJson {
    id: number;
    CustomerId: number;
    rows: DeliveryRow[];
}

export interface DeliveryRow {
    DeliverySchedule: DeliverySchedule;
    entityId: number;
    id: number;
    nmc: number;
    purchaseAmount: number | null;
    spgzCharacteristics: SpgzCharacteristics;
}

export interface DeliverySchedule {
    start_date: string;
    end_date: string;
    year: number;
    deliveryAmount: number | null;
}

export interface SpgzCharacteristics {
    spgzId: number;
    spgzName: string;
    kpgzCode: string;
    kpgzName: string;
}
