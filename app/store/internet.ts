import { urlHomePage } from '~/components/internet/helper';
import type { GetAdsRequest, GetAdsResponse } from '~~/gen/ts/services/internet/ads';

const logger = useLogger('🌐 Internet');

export type Tab = {
    id: number;
    label: string;
    url: string;
    icon?: string;
    active?: boolean;
};

export interface InternetState {
    selectedTab: number | undefined;
    tabs: Tab[];
    history: string[];
}

export const useInternetStore = defineStore('internet', {
    state: () =>
        ({
            selectedTab: 0,
            tabs: [],
            history: [],
        }) as InternetState,
    persist: {
        pick: ['selectedTab', 'tabs'],
    },
    actions: {
        // Tabs
        async newTab(select?: boolean): Promise<void> {
            const id = await this.addTab({
                label: '',
                url: urlHomePage,
                icon: 'i-mdi-home',
                active: select,
            });

            if (select === true) {
                this.selectTab(id);
            }
        },
        async addTab(tab: Partial<Tab>): Promise<number> {
            const id = this.tabs.length;
            this.tabs.push({
                id: id,
                label: tab.label ?? '',
                url: tab.url ?? '',
                icon: tab.icon,
                active: tab.active ?? false,
            });

            return id;
        },
        async closeTab(id: number): Promise<void> {
            if (this.selectedTab === id) {
                this.selectTab();
            }

            const idx = this.tabs.findIndex((t) => t.id === id);
            if (idx === -1) {
                return;
            }

            this.tabs.splice(idx, 1);
        },
        async selectTab(id?: number): Promise<void> {
            this.selectedTab = id;
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
});
