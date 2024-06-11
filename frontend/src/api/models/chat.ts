export interface ShortSession {
    id: string;
    title: string;
    created_at: string;
    tg: boolean;
}

export interface ChatQuery {
    id: number;
    prompt: string;
    product: string;
    type: string;
    status: string;
    created_at: string;
}

export interface ChatResponse {
    created_at: string;
    body: string;
    status: string;
}

export interface SessionContent {
    query: ChatQuery;
    response: ChatResponse;
}

export interface ChatSession {
    id: string;
    title: string;
    content: SessionContent[];
    editable: boolean;
    tg: boolean;
}

export interface GetSessionsResponse {
    sessions: ShortSession[];
}

export interface CreateSessionResponse {
    id: string;
}

export interface GetSessionParams {
    id: string;
}

export interface RenameSessionParams {
    id: string;
    title: string;
}

export interface DeleteSessionParams {
    id: string;
}
