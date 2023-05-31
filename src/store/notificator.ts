import { StoreDefinition, defineStore } from 'pinia';

export interface NotificatorState {
    lastId: string;
}

export const useNotificatorStore = defineStore('notificator', {
    state: () =>
        ({
            lastId: '0',
        } as NotificatorState),
    persist: true,
    actions: {
        setLastId(lastId: bigint): void {
            this.lastId = lastId.toString();
        },
    },
    getters: {
        getLastId: (state): bigint => BigInt(state.lastId),
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificatorStore as unknown as StoreDefinition, import.meta.hot));
}
