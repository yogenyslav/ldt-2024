import ChatApiService from '@/api/ChatApiService';
import { ShortSession } from '@/api/models';
import { makeAutoObservable } from 'mobx';

export class RootStore {
    sessions: ShortSession[] = [];
    sessionsLoading: boolean = false;

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

    deleteSession() {}
}
