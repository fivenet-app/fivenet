import { InjectionKey } from 'vue';
import { createStore, useStore as baseUseStore, Store } from 'vuex';
import VuexPersistence from 'vuex-persist';
import authmodule, { AuthModuleState } from './modules/authmodule';
import clipboardModule, { ClipboardModuleState } from './modules/clipboardmodule';
import loadermodule, { LoaderModuleState } from './modules/loadermodule';

export interface RootState {
    version: string;
    loader?: LoaderModuleState;
    auth?: AuthModuleState;
    clipboard?: ClipboardModuleState;
}

const vuexPersist = new VuexPersistence<RootState>({
    key: 'arpanet',
    storage: window.localStorage,
    modules: ['auth', 'clipboard'],
    reducer: (state: RootState) => ({
        auth: {
            accessToken: state.auth?.accessToken,
            lastCharID: state.auth?.lastCharID,
        },
        clipboard: {
            documents: state.clipboard?.documents,
            users: state.clipboard?.users,
            vehicles: state.clipboard?.vehicles,
            activeStack: {
                documents: state.clipboard?.activeStack.documents,
                users: state.clipboard?.activeStack.users,
                vehicles: state.clipboard?.activeStack.vehicles,
            },
        },
    }),
});

export const key: InjectionKey<Store<RootState>> = Symbol();

export const store = createStore<RootState>({
    plugins: [vuexPersist.plugin],
    modules: {
        auth: authmodule,
        clipboard: clipboardModule,
        loader: loadermodule,
    },
    state: {
        version: '',
    },
    mutations: {},
});

export function useStore() {
    return baseUseStore(key);
}
