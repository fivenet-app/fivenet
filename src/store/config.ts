import { defineStore, type StoreDefinition } from 'pinia';

export interface ConfigState {
    nuiEnabled: boolean;
    nuiResourceName: string | undefined;
    updateAvailable: false | string;
}

export const useConfigStore = defineStore('config', {
    state: () =>
        ({
            nuiEnabled: false,
            nuiResourceName: undefined,
            updateAvailable: false,
        }) as ConfigState,
    persist: {
        paths: ['nuiEnabled', 'nuiResourceName'],
    },
    actions: {
        async setUpdateAvailable(version: string): Promise<void> {
            this.updateAvailable = version;
        },
        setNuiDetails(enabled: boolean, resourceName: string | undefined): void {
            this.nuiEnabled = enabled;
            this.nuiResourceName = resourceName;
        },
    },
    getters: {
        isNUIAvailable(state): boolean {
            return state.nuiEnabled ?? false;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useConfigStore as unknown as StoreDefinition, import.meta.hot));
}
