import { joinURL, splitURL, urlHomePage } from '~/components/internet/helper';
import type { GetAdsRequest, GetAdsResponse } from '~~/gen/ts/services/internet/ads';

const logger = useLogger('🌐 Internet');

export type Tab = {
    id: number;
    label: string;
    domain: string;
    path: string;
    icon?: string;
    active?: boolean;
    history: string[];
};

export interface InternetState {
    selectedTab: number | undefined;
    tabs: Tab[];
}

export const useInternetStore = defineStore('internet', {
    state: () =>
        ({
            selectedTab: undefined,
            tabs: [],
        }) as InternetState,
    persist: {
        pick: ['selectedTab', 'tabs'],
    },
    actions: {
        // Tabs
        newTab(select?: boolean): void {
            const id = this.addTab({
                domain: urlHomePage,
                icon: 'i-mdi-home',
                active: select,
            });

            if (select === true) {
                this.selectTab(id);
            }
        },
        addTab(tab: Partial<Tab>): number {
            const id = this.tabs.length === 0 ? 1 : this.tabs.length + 1;
            logger.debug('tab added, label:', tab.label, 'domain:', tab.domain, 'path:', tab.path);
            this.tabs.push({
                id: id,
                label: tab.label ?? '',
                domain: tab.domain ?? urlHomePage,
                path: tab.path ?? '',
                icon: tab.icon,
                active: tab.active ?? false,
                history: [],
            });

            return id;
        },
        closeTab(id: number): void {
            if (this.selectedTab === id) {
                this.selectTab();
            }

            const idx = this.tabs.findIndex((t) => t.id === id);
            if (idx === -1) {
                return;
            }

            logger.debug('close tab, id:', id);
            this.tabs.splice(idx, 1);
        },
        selectTab(id?: number): void {
            logger.debug('select tab, id:', id);
            this.selectedTab = id;

            this.tabs.forEach((t) => (t.active = t.id === id));
        },

        // Navigation
        goTo(domain: string, path: string = '', disableHistory: boolean = false): void {
            const tab = this.activeTab;
            if (!tab) {
                return;
            }

            logger.debug('goto, domain:', tab.domain, 'path:', tab.path);

            if (!disableHistory) {
                tab.history.push(joinURL(tab.domain, tab.path));
            }

            tab.domain = domain;
            tab.path = path;
        },
        back(): void {
            const tab = this.activeTab;
            if (!tab) {
                return;
            }

            const url = tab.history.pop();
            logger.debug('back, url:', url);
            if (url) {
                const split = splitURL(url);
                if (!split) {
                    return;
                }

                this.goTo(split.domain, split.path, true);
            }
        },

        // Ads
        async getAds(req: GetAdsRequest): Promise<GetAdsResponse> {
            try {
                const call = getGRPCInternetAdsClient().getAds(req);
                const { response } = await call;

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
    getters: {
        activeTab: (state) => state.tabs.find((t) => t.id === state.selectedTab),
    },
});
