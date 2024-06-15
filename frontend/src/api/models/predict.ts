import { OutputJson } from './favorites';

export interface PurchasePlan {
    date: string;
    value: number; // стоимость закупки в рублях
    volume: number | null; // объем закупки в штуках
}

export enum ModelResponseType {
    Prediction = 'PREDICTION',
    Stock = 'STOCK',
}

export interface PredictionResponse {
    code: string;
    code_name: string;
    example_contracts_in_code: ExampleContract[];
    forecast: PurchasePlan[];
    history: PurchasePlan[];
    is_regular: boolean;
    mean_ref_price: number | null;
    mean_start_to_execute_days: number;
    median_execution_days: number;
    top5_providers: number[];
    output_json: OutputJson;
}

export interface ExampleContract {
    characteristics_name: string;
    conclusion_date: string;
    end_date_of_validity: string;
    execution_term_from: string;
    execution_term_until: string;
    final_code_kpgz: string;
    final_name_kpgz: string;
    gk_price_rub: number;
    id: number;
    id_spgz: number;
    item_name_gk: string;
    name_spgz: string;
    name_ste: string;
    paid_rub: number;
    provider: number;
    ref_price: number | null;
    registry_number_in_rk: string;
}

export interface StockResponse {
    data: Stock[];
}

export interface Stock {
    amount: number;
    name: string;
    price: number;
    quarter: number;
    sum: number;
    year: number;
}
