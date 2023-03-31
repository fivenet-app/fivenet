import { defineStore } from 'pinia';

export interface DocumentEditorState {
    title: string;
    content: string;
    closed: undefined | { id: number; label: string; closed: boolean };
    state: string;
}

export const useDocumentEditorStore = defineStore('documentEditor', {
    state: () => ({
        title: '',
        content: '',
        closed: undefined,
        state: '',
    } as DocumentEditorState),
    persist: true,
    actions: {
        save(doc: DocumentEditorState) {
            this.title = doc.title;
            this.content = doc.content;
            this.closed = doc.closed;
            this.state = doc.state;
        },
        clear(): void {
            this.title = '';
            this.content = '';
            this.closed = undefined;
            this.state = '';
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useDocumentEditorStore, import.meta.hot));
}
