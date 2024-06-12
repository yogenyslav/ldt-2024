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

export interface WSOutcomingMessage {
    prompt: string;
    command?: ChatCommand;
    period?: string;
    product?: string;
}

export enum ChatCommand {
    Valid = 'valid',
    Invalid = 'invalid',
    Cancel = 'cancel',
}

export enum IncomingMessageType {
    Stock = 'STOCK',
    Prediction = 'PREDICTION',
    Undefined = 'UNDEFINED',
}

export enum IncomingMessageStatus {
    Pending = 'PENDING',
    Valid = 'VALID',
}

export interface WSMessage {
    data: WSIncomingQuery | WSIncomingChunk;
    finish: boolean;
    chunk: boolean;
}

export interface WSIncomingQuery {
    created_at: string;
    prompt: string;
    period: string;
    product: string;
    type: IncomingMessageType;
    status: string;
    id: number;
}

export interface WSIncomingChunk {
    info: string;
}

export interface ChatConversation {
    outcomingMessage?: DisplayedOutcomingMessage;
    incomingMessage?: DisplayedIncomingMessage;
}

export interface DisplayedChat {
    messages: ChatConversation[];
}

export interface DisplayedOutcomingMessage {
    prompt: string;
}

export interface DisplayedIncomingMessage {
    type: IncomingMessageType;
    status: IncomingMessageStatus;
    body: string;
    product?: string;
    period?: string;
}
