import { InjectionKey } from 'vue';
import { createStore, useStore as baseUseStore, Store } from 'vuex';
import VuexPersistence from 'vuex-persist';
import authmodule, { AuthModuleState } from './modules/authmodule';
import clipboardModule, { ClipboardModuleState } from './modules/clipboardmodule';
import loadermodule, { LoaderModuleState } from './modules/loadermodule';
import documentEditorModule, { DocumentEditorModuleState } from './modules/documenteditormodule';

export interface RootState {
    version: string;
    loader?: LoaderModuleState;
    auth?: AuthModuleState;
    clipboard?: ClipboardModuleState;
    documentEditor?: DocumentEditorModuleState;
}

const vuexPersist = new VuexPersistence<RootState>({
    key: 'arpanet',
    storage: window.localStorage,
    modules: ['auth', 'clipboard', 'documentEditor'],
    reducer: (state: RootState) => ({
        version: state.version,
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
        documentEditor: {
            title: state.documentEditor?.title,
            content: state.documentEditor?.content,
            closed: state.documentEditor?.closed,
            state: state.documentEditor?.state,
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
        documentEditor: documentEditorModule,
    },
    state: {
        version: __APP_VERSION__,
    },
    strict: import.meta.env.DEV,
});

export function useStore() {
    return baseUseStore(key);
}
