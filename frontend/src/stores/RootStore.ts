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
                this.activeSessionId = null;
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
                return {
                    incomingMessage: {
                        body: content.response.body,
                        type: content.query.type as IncomingMessageType,
                        status: content.query.status as IncomingMessageStatus,
                    },
                    outcomingMessage: {
                        prompt: content.query.prompt,
                    },
                };
            }),
        };

        this.connectWebSocket(session.id);
    }

    setActiveSessionId(id: string) {
        if (id !== this.activeSession?.id) {
            // this.disconnectWebSocket();

            this.activeSessionId = id;
        }
    }

    renameSession({ id, title }: RenameSessionParams) {
        return ChatApiService.renameSession({ id, title });
    }

    async createSession() {
        return ChatApiService.createSession().then(async ({ id }) => {
            this.activeDisplayedSession = null;

            this.setActiveSessionId(id);

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
            }

            this.setActiveSessionId(sessionId);
        };

        this.websocket.onmessage = (event) => {
            const wsMessage: WSMessage = JSON.parse(event.data);

            runInAction(() => {
                const data = wsMessage.data;

                console.log(wsMessage);

                if (wsMessage.chunk && !wsMessage.finish) {
                    // this.isModelAnswering = true;
                    // this.isChatDisabled = true;

                    this.processIncomingChunk(data as WSIncomingChunk);
                } else if (!wsMessage.chunk && wsMessage.data) {
                    this.processIncomingQuery(data as WSIncomingQuery);
                }

                if (wsMessage.finish || !wsMessage.chunk) {
                    this.isModelAnswering = false;
                }

                if (wsMessage.finish) {
                    this.isChatDisabled = false;
                }
            });
        };

        this.websocket.onclose = () => {
            console.log('WebSocket connection closed');
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
                title: message.prompt.slice(0, 25),
            });
        }

        this.addMessageToActiveSession(message);
    }

    disconnectWebSocket() {
        if (this.websocket) {
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
                    prompt: message.prompt,
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
                body: lastMessageBody ? lastMessageBody + info : info,
                type: IncomingMessageType.Undefined,
                status: IncomingMessageStatus.Valid,
            };
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
            prompt: '',
            command: ChatCommand.Cancel,
        });

        this.setChatDisabled(false);
        this.setIsModelAnswering(false);
    }

    private isFirstMessageInSession() {
        return !this.activeDisplayedSession?.messages.length;
    }
}
