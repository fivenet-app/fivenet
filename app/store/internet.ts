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
    tabs: Tab[];
}

export const useInternetStore = defineStore('internet', {
    state: () =>
        ({
            selectedTab: undefined,
            tabs: [],
        }) as InternetState,
    persist: {
        pick: ['tabs'],
    },
    actions: {
        // Tabs
        newTab(select?: boolean): void {
            this.addTab({
                domain: urlHomePage,
                icon: 'i-mdi-home',
                active: select,
            });
        },
        addTab(tab: Partial<Tab>, select: boolean = true): number {
            const id = this.tabs.length === 0 ? 1 : this.tabs[this.tabs.length - 1]!.id + 1;
            logger.debug('tab added, id:', id, 'label:', tab.label, 'domain:', tab.domain, 'path:', tab.path);
            this.tabs.push({
                id: id,
                label: tab.label ?? '',
                domain: tab.domain ?? urlHomePage,
                path: tab.path ?? '/',
                icon: tab.icon,
                active: tab.active ?? false,
                history: [],
            });

            if (select) {
                this.selectTab(id);
            }

            return id;
        },
        closeTab(id: number): void {
            const idx = this.tabs.findIndex((t) => t.id === id);
            if (idx === -1) {
                return;
            }

            logger.debug('close tab, id:', id);
            this.tabs.splice(idx, 1);

            // Attempt to find a close by tab if needed
            if (this.activeTab?.id === id) {
                this.selectTab();

                if (this.tabs.length > 0) {
                    const idx = this.tabs.findIndex((t) => t.id === id - 1);
                    if (idx === -1) {
                        this.selectTab(this.tabs[0]?.id);
                    } else {
                        this.selectTab(this.tabs[idx]?.id);
                    }
                }
            }
        },
        selectTab(id?: number): void {
            logger.debug('select tab, id:', id);

            this.tabs.forEach((t) => (t.active = t.id === id));
        },

        // Navigation
        goTo(domain: string, path: string = '', disableHistory: boolean = false): void {
            const tab = this.activeTab;
            if (!tab) {
                return;
            }

            logger.debug('goto, domain:', domain, 'path:', path);

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
        activeTab: (state) => state.tabs.find((t) => t.active),
    },
});
