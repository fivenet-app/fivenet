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
    tab: Tab;
}

export const useInternetStore = defineStore('internet', {
    state: () =>
        ({
            tab: {
                id: 0,
                label: '',
                domain: urlHomePage,
                path: '/',
                history: [],
                active: true,
            },
        }) as InternetState,
    persist: {
        pick: ['tab'],
    },
    actions: {
        // Navigation
        goTo(domain: string, path: string = '', disableHistory: boolean = false): void {
            logger.debug('goto, domain:', domain, 'path:', path);

            if (!disableHistory) {
                this.tab.history.push(joinURL(this.tab.domain, this.tab.path));
            }

            this.tab.domain = domain;
            this.tab.path = path;
        },
        back(): void {
            const url = this.tab.history.pop();
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
});
