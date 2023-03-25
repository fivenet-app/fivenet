import { Module } from 'vuex';
import { RootState } from '../store';

export interface ClipboardModuleState {
    documents: Array<ClipboardItem>;
    users: Array<ClipboardItem>;
}

export type ClipboardItemType = 'Char' | 'Document';

export class ClipboardItem {
    private type: ClipboardItemType;
    private id: number;
    private data: Object;

    constructor(type: ClipboardItemType, id: number, data: Object) {
        this.type = type;
        this.id = id;
        this.data = data;
    }

    getType(): ClipboardItemType {
        return this.type;
    }

    getId(): number {
        return this.id;
    }

    getData(): Object {
        return this.data;
    }
}

const clipboardModule: Module<ClipboardModuleState, RootState> = {
    namespaced: true,
    state: {
        documents: new Array<ClipboardItem>(),
        users: new Array<ClipboardItem>(),
    },
    actions: {
        addItem({ commit }, data: ClipboardItem) {
            commit('addItem', data);
        },
        removeItem({ commit }, key: { type: ClipboardItemType; id: number }) {
            commit('removeItem', key);
        },
        clear({ commit }) {
            commit('clear');
        },
    },
    mutations: {
        addItem(state: ClipboardModuleState, item: ClipboardItem ): void {
            switch (item.getType()) {
                case 'Char':
                    const charIdx = state.users.findIndex((o: ClipboardItem) => {
                        return o.getId() === item.getId();
                    });
                    if (charIdx === -1) {
                        console.log("ADDING CHAR TO CLIPBOARD");
                        state.users.push(item);
                    }
                    break;
                case 'Document':
                    const docIdx = state.documents.findIndex((o: ClipboardItem) => {
                        return o.getId() === item.getId();
                    });
                    if (docIdx === -1) {
                        state.documents.push(item);
                    }
                    break;
            }
        },
        removeItem(state: ClipboardModuleState, item: { type: ClipboardItemType; id: number }): void {
            switch (item.type) {
                case 'Char':
                    const charIdx = state.users.findIndex((o: ClipboardItem) => {
                        return o.getId() === item.id;
                    });
                    if (charIdx > 0) {
                        state.users.splice(charIdx, 1);
                    }
                    break;
                case 'Document':
                    const docIdx = state.documents.findIndex((o: ClipboardItem) => {
                        return o.getId() === item.id;
                    });
                    if (docIdx > 0) {
                        state.documents.splice(docIdx, 1);
                    }
            }
        },
        clear(state: ClipboardModuleState) {
            state.users.length = 0;
            state.documents.length = 0;
        },
    },
};

export default clipboardModule;

export interface ClipboardData {
    users: Array<Object>;
    documents: Array<Object>;
}

export function getClipboardData(state: ClipboardModuleState): ClipboardData {
    const data: ClipboardData = {
        documents: [],
        users: [],
    };
    state.documents.forEach((v) => {
        data.documents.push(v.getData());
    })
    state.users.forEach((v) => {
        data.users.push(v.getData());
    })

    return data;
}
