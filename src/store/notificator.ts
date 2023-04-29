import { StatusCode } from 'grpc-web';
import { StoreDefinition, defineStore } from 'pinia';

export interface NotificatorState {
    lastId: number;
}

export const useNotificatorStore = defineStore('notificator', {
    state: () =>
    ({
        lastId: 0,
    } as NotificatorState),
    persist: true,
    actions: {
        setLastId(lastId: number): void {
            this.lastId = lastId;
        },
    },
    getters: {
        getLastId: (state): number => state.lastId,
    }
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificatorStore as unknown as StoreDefinition, import.meta.hot));
}
