import { defineStore } from 'pinia';
import type { Category } from '~~/gen/ts/resources/documents/category';

export interface DocumentEditorState {
    title: string;
    content: string;
    closed: boolean;
    state: string;
    category?: Category;
}

export const useDocumentEditorStore = defineStore('documentEditor', {
    state: () =>
        ({
            title: '',
            content: '',
            closed: false,
            state: '',
            category: undefined,
        }) as DocumentEditorState,
    persist: true,
    actions: {
        save(doc: DocumentEditorState): void {
            this.title = doc.title;
            this.content = doc.content;
            this.closed = doc.closed;
            this.state = doc.state;
            this.category = doc.category;
        },
        clear(): void {
            this.title = '';
            this.content = '';
            this.closed = false;
            this.state = '';
            this.category = undefined;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useDocumentEditorStore, import.meta.hot));
}
