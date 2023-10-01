import { NuxtError } from 'nuxt/app';
import { StoreDefinition, defineStore } from 'pinia';

export interface ConfigState {
    fetched: boolean;
    appConfig: AppConfig;
    clientConfig: ClientConfig;
    updateAvailable: false | string;
}

type AppConfig = {
    sentryDSN?: string;
    login: LoginConfig;
    discord: DiscordConfig;
};

type LoginConfig = {
    signupEnabled: boolean;
    providers: ProviderConfig[];
};

type ProviderConfig = {
    name: string;
    label: string;
};

type DiscordConfig = {
    botInviteURL?: string;
};

type ClientConfig = {
    NUIEnabled: boolean;
    NUIResourceName?: string;
};

export const useConfigStore = defineStore('config', {
    state: () =>
        ({
            fetched: false,
            appConfig: {
                sentryDSN: '',
                login: {
                    signupEnabled: true,
                    providers: [],
                },
                discord: {
                    botInviteURL: '',
                },
            } as AppConfig,
            clientConfig: {
                NUIEnabled: false,
            } as ClientConfig,
            updateAvailable: false,
        }) as ConfigState,
    persist: {
        paths: ['clientConfig'],
    },
    actions: {
        async loadConfig(): Promise<void> {
            return new Promise(async (res, rej) => {
                if (this.fetched) {
                    return res();
                }

                try {
                    // 6 seconds should be enough
                    const abort = new AbortController();
                    const tId = setTimeout(() => abort.abort(), 8000);

                    const resp = await fetch('/api/config', {
                        method: 'POST',
                        signal: abort.signal,
                    });
                    clearTimeout(tId);

                    if (!resp.ok) {
                        const text = await resp.text();
                        throw createError({
                            statusCode: 500,
                            statusMessage: 'Failed to get FiveNet config from backend',
                            message: text,
                            fatal: true,
                            unhandled: false,
                        });
                    }
                    const data = (await resp.json()) as AppConfig;
                    this.appConfig = data;

                    this.fetched = true;

                    return res();
                } catch (e) {
                    showError(e as NuxtError);
                    return rej(e);
                }
            });
        },
        setUpdateAvailable(version: string): void {
            this.updateAvailable = version;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useConfigStore as unknown as StoreDefinition, import.meta.hot));
}
