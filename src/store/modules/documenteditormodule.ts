import { Module } from 'vuex';
import { RootState } from '../store';

export interface DocumentEditorModuleState {
    title: string;
    content: string;
    closed: { id: number, label: string, closed: boolean } | undefined;
    state: string;
}

const documentEditorModule: Module<DocumentEditorModuleState, RootState> = {
    namespaced: true,
    state: {
        title: '',
        content: '',
        closed: undefined,
        state: '',
    },
    actions: {
        save({ commit }, doc: DocumentEditorModuleState) {
            commit('save', doc);
        },
        clear({ commit }): void {
            commit('clear');
        }
    },
    mutations: {
        save(state: DocumentEditorModuleState, doc: DocumentEditorModuleState) {
            state.title = doc.title;
            state.content = doc.content;
            state.closed = doc.closed;
            state.state = doc.state;
        },
        clear(state: DocumentEditorModuleState): void {
            state.title = '';
            state.content = '';
            state.closed = undefined;
            state.state = '';
        }
    },
};

export default documentEditorModule;
