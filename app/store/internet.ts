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

export const useInternetStore = defineStore(
    'internet',
    () => {
        // State
        const tab = ref<Tab>({
            id: 0,
            label: '',
            domain: urlHomePage,
            path: '/',
            history: [],
            active: true,
        });

        // Actions
        const goTo = (domain: string, path: string = '', disableHistory: boolean = false): void => {
            if (path.length === 0) {
                path = '/';
            } else if (!path.startsWith('/')) {
                path = '/' + path;
            }

            logger.debug('goto, domain:', domain, 'path:', path);

            if (!disableHistory) {
                tab.value.history.push(joinURL(tab.value.domain, tab.value.path));
            }

            tab.value.domain = domain;
            tab.value.path = path;
        };

        const back = (): void => {
            const url = tab.value.history.pop();
            logger.debug('back, url:', url);
            if (url) {
                const split = splitURL(url);
                if (!split) {
                    return;
                }
                goTo(split.domain, split.path, true);
            }
        };

        // Ads
        const getAds = async (req: GetAdsRequest): Promise<GetAdsResponse> => {
            try {
                const call = getGRPCInternetAdsClient().getAds(req);
                const { response } = await call;
                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        return {
            tab,
            goTo,
            back,
            getAds,
        };
    },
    {
        persist: {
            pick: ['tab'],
        },
    },
);
