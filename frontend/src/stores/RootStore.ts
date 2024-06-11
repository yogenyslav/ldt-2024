import ChatApiService from '@/api/ChatApiService';
import {
    ChatSession,
    DeleteSessionParams,
    GetSessionParams,
    RenameSessionParams,
    ShortSession,
} from '@/api/models';
import { makeAutoObservable } from 'mobx';

export class RootStore {
    sessions: ShortSession[] = [];
    sessionsLoading: boolean = false;

    activeSessionId: string | null = null;
    activeSession: ChatSession | null = null;
    activeSessionLoading: boolean = false;

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
    }

    setActiveSessionId(id: string) {
        if (id !== this.activeSession?.id) {
            this.activeSessionId = id;
        }
    }

    renameSession({ id, title }: RenameSessionParams) {
        return ChatApiService.renameSession({ id, title });
    }

    async createSession() {
        return ChatApiService.createSession().then(({ id }) => {
            this.getSessions();

            // activate web socket connection and pass token as a first
            // this.setActiveSessionId(id);
        });
    }
}
