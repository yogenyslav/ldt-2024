import ChatApiService from '@/api/ChatApiService';
import {
    ChatSession,
    DeleteSessionParams,
    GetSessionParams,
    WSOutcomingMessage,
    RenameSessionParams,
    ShortSession,
    DisplayedChat,
    WSIncomingQuery,
    IncomingMessageType,
    IncomingMessageStatus,
    WSMessage,
    WSIncomingChunk,
    ChatCommand,
    ModelResponseType,
    PredictionResponse,
    StockResponse,
    UNAUTHORIZED_ERR,
} from '@/api/models';
import { LOCAL_STORAGE_KEY } from '@/auth/AuthProvider';
import { WS_URL } from '@/config';
import { makeAutoObservable, runInAction } from 'mobx';

export class RootStore {
    sessions: ShortSession[] = [];
    sessionsLoading: boolean = false;

    activeSessionId: string | null = null;
    activeSession: ChatSession | null = null;
    activeDisplayedSession: DisplayedChat | null = null;
    activeSessionLoading: boolean = false;
    isChatDisabled: boolean = false;
    isModelAnswering: boolean = false;
    chatError: string | null = null;
    closedWebSocket: WebSocket | null = null;

    websocket: WebSocket | null = null;

    constructor() {
        makeAutoObservable(this);

        this.sessions = [];
    }

    async getSessions() {
        this.sessionsLoading = true;

        return ChatApiService.getSessions()
            .then(({ sessions }) => {
                this.sessions = sessions;
            })
            .finally(() => {
                this.sessionsLoading = false;
            });
    }

    async deleteSession({ id }: DeleteSessionParams) {
        return ChatApiService.deleteSession({ id }).then(() => {
            if (this.activeSessionId === id) {
                this.setActiveSessionId(null);
                this.activeSession = null;
            }
        });
    }

    async getSession({ id }: GetSessionParams) {
        this.activeSessionLoading = true;

        return ChatApiService.getSession({ id })
            .then((session) => {
                this.setActiveSession(session);
            })
            .finally(() => {
                this.activeSessionLoading = false;
            });
    }

    setActiveSession(session: ChatSession) {
        this.activeSession = session;

        this.activeDisplayedSession = {
            messages: session.content.map((content) => {
                const prediction =
                    content.response.data_type === ModelResponseType.Prediction
                        ? (content.response.data as PredictionResponse)
                        : null;

                const stocks =
                    content.response.data_type === ModelResponseType.Stock
                        ? (content.response.data as StockResponse)
                        : null;

                return {
                    incomingMessage: {
                        body: content.response.body,
                        type: content.query.type as IncomingMessageType,
                        status: content.query.status as IncomingMessageStatus,
                        product: content.query.product,
                        period: content.query.period,
                        prediction: prediction
                            ? {
                                  forecast: prediction.forecast,
                                  history: prediction.history,
                              }
                            : undefined,
                        stocks: stocks?.data,
                        outputJson: prediction?.output_json,
                    },
                    outcomingMessage: {
                        prompt: content.query.prompt,
                    },
                };
            }),
        };

        this.connectWebSocket(session.id);
    }

    setActiveSessionId(id: string | null) {
        if (id !== this.activeSessionId) {
            this.activeSessionId = id;
        }
    }

    renameSession({ id, title }: RenameSessionParams) {
        return ChatApiService.renameSession({ id, title });
    }

    async createSession() {
        return ChatApiService.createSession().then(async ({ id }) => {
            this.activeDisplayedSession = null;

            this.getSessions();

            this.connectWebSocket(id);
        });
    }

    connectWebSocket(sessionId: string) {
        this.disconnectWebSocket();

        this.websocket = new WebSocket(`${WS_URL}/${sessionId}`);

        this.websocket.onopen = () => {
            console.log('WebSocket connection opened');

            if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
                this.websocket.send(
                    JSON.parse(localStorage.getItem(LOCAL_STORAGE_KEY) as string)?.user?.token
                );

                this.setChatDisabled(false);
            }

            this.setActiveSessionId(sessionId);
        };

        this.websocket.onmessage = (event) => {
            const wsMessage: WSMessage = JSON.parse(event.data);

            runInAction(() => {
                const data = wsMessage.data;

                console.log(wsMessage);

                if (wsMessage.err) {
                    this.chatError = wsMessage.err;

                    if (wsMessage.err === UNAUTHORIZED_ERR) {
                        localStorage.removeItem(LOCAL_STORAGE_KEY);

                        window.location.href = '/login';
                    }

                    this.isModelAnswering = false;
                    this.isChatDisabled = false;
                }

                if (wsMessage.chunk && !wsMessage.finish) {
                    // this.isModelAnswering = true;
                    // this.isChatDisabled = true;

                    this.processIncomingChunk(data as WSIncomingChunk);
                } else if (!wsMessage.chunk && wsMessage.data && !wsMessage.data_type) {
                    //!wsMessage.data_type значит, что это ответ модели (prediction или stock)
                    this.processIncomingQuery(data as WSIncomingQuery);
                } else if (wsMessage.data_type === ModelResponseType.Prediction && wsMessage.data) {
                    this.processIncomingPrediction(data as PredictionResponse);
                } else if (wsMessage.data_type === ModelResponseType.Stock && wsMessage.data) {
                    this.processIncomingStock(data as StockResponse);
                }

                if (
                    (wsMessage.finish || !wsMessage.chunk) &&
                    !(wsMessage.data_type === ModelResponseType.Prediction)
                ) {
                    this.isModelAnswering = false;
                }

                if (wsMessage.finish || wsMessage.data_type === ModelResponseType.Stock) {
                    this.isChatDisabled = false;
                }
            });
        };

        this.websocket.onclose = () => {
            console.log('WebSocket connection closed');

            this.isChatDisabled = true;
            this.closedWebSocket = this.websocket;

            this.reconnectWebSocket();
        };

        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    sendMessage(message: WSOutcomingMessage) {
        console.log('sendMessage', message);

        this.setIsModelAnswering(true);
        this.setChatDisabled(true);

        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify(message));
        }

        if (this.isFirstMessageInSession()) {
            this.renameSession({
                id: this.activeSessionId as string,
                title: message.prompt?.slice(0, 60) || 'Без названия',
            });
        }

        this.addMessageToActiveSession(message);
    }

    disconnectWebSocket() {
        if (this.websocket) {
            this.setActiveSessionId(null);
            this.websocket.close();
        }
    }

    addMessageToActiveSession(message: WSOutcomingMessage) {
        if (!this.activeSessionId) {
            return;
        }

        runInAction(() => {
            if (!this.activeDisplayedSession) {
                this.activeDisplayedSession = { messages: [] };
            }

            this.activeDisplayedSession?.messages.push({
                outcomingMessage: {
                    prompt: message.prompt || '',
                },
            });
        });
    }

    private processIncomingQuery(query: WSIncomingQuery) {
        console.log('processIncomingQuery', query);

        if (this.activeSessionId && this.activeDisplayedSession?.messages.length) {
            this.activeDisplayedSession.messages[
                this.activeDisplayedSession.messages.length - 1
            ].incomingMessage = {
                body: query.prompt,
                type: query.type as IncomingMessageType,
                status: query.status as IncomingMessageStatus,
                product: query.product,
                period: query.period,
            };
        }
    }

    private processIncomingChunk({ info }: WSIncomingChunk) {
        if (this.activeSessionId && this.activeDisplayedSession?.messages.length) {
            const lastMessageIndex = this.activeDisplayedSession.messages.length - 1;
            const lastMessageBody =
                this.activeDisplayedSession.messages[lastMessageIndex].incomingMessage?.body;

            this.activeDisplayedSession.messages[lastMessageIndex].incomingMessage = {
                ...this.activeDisplayedSession.messages[lastMessageIndex].incomingMessage,
                body: lastMessageBody ? lastMessageBody + info : info,
                type: IncomingMessageType.Undefined,
                status: IncomingMessageStatus.Valid,
            };
        }
    }

    private processIncomingPrediction(data: PredictionResponse) {
        const { forecast, history } = data;
        console.log('processIncomingPrediction', forecast, history);

        const session = this.activeDisplayedSession;
        if (!this.activeSessionId || !session?.messages.length) return;

        const lastMessageIndex = session.messages.length - 1;
        const lastMessage = session.messages[lastMessageIndex];

        const incomingMessage = lastMessage.incomingMessage || {
            body: '',
            type: IncomingMessageType.Undefined,
            status: IncomingMessageStatus.Valid,
            prediction: { forecast, history },
        };

        incomingMessage.prediction = { forecast, history };
        lastMessage.incomingMessage = incomingMessage;
        lastMessage.incomingMessage.outputJson = data.output_json;

        if (this.activeDisplayedSession) {
            this.activeDisplayedSession.messages[lastMessageIndex] = lastMessage;
        }
    }

    private processIncomingStock({ data }: StockResponse) {
        console.log('processIncomingStock', data);

        const session = this.activeDisplayedSession;
        if (!this.activeSessionId || !session?.messages.length) return;

        const lastMessageIndex = session.messages.length - 1;
        const lastMessage = session.messages[lastMessageIndex];

        const incomingMessage = lastMessage.incomingMessage || {
            body: '',
            type: IncomingMessageType.Undefined,
            status: IncomingMessageStatus.Valid,
        };

        incomingMessage.stocks = data;
        lastMessage.incomingMessage = incomingMessage;

        if (this.activeDisplayedSession) {
            this.activeDisplayedSession.messages[lastMessageIndex] = lastMessage;
        }
    }

    setChatDisabled(isDisabled: boolean) {
        this.isChatDisabled = isDisabled;
    }

    setIsModelAnswering(isAnswering: boolean) {
        this.isModelAnswering = isAnswering;
    }

    cancelRequest() {
        this.sendMessage({
            command: ChatCommand.Cancel,
        });

        this.setChatDisabled(false);
        this.setIsModelAnswering(false);
    }

    private isFirstMessageInSession() {
        return !this.activeDisplayedSession?.messages.length;
    }

    private reconnectWebSocket() {
        if (this.activeSessionId) {
            this.connectWebSocket(this.activeSessionId);
        }
    }
}
