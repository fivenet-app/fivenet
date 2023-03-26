import { Module } from 'vuex';
import { RootState } from '../store';

export interface LoaderModuleState {
    loading: number,
}

const loaderModule: Module<LoaderModuleState, RootState> = {
    namespaced: true,
    state: {
        loading: 0,
    },
    actions: {
        show({ commit }) {
            commit('show');
        },
        hide({ commit }) {
            commit('hide');
        },
    },
    mutations: {
        show(state: LoaderModuleState) {
            state.loading++;
        },
        hide(state: LoaderModuleState) {
            if (state.loading > 0) {
                state.loading--;
            }
        },
    },
};

export default loaderModule;
