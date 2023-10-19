import { StoreDefinition, defineStore } from 'pinia';
import { Category } from '~~/gen/ts/resources/documents/category';

export interface DocumentEditorState {
    title: string;
    content: string;
    closed: undefined | { id: number; label: string; closed: boolean };
    state: string;
    category?: Category;
}

export const useDocumentEditorStore = defineStore('documentEditor', {
    state: () =>
        ({
            title: '',
            content: '',
            closed: undefined,
            state: '',
            category: undefined,
        }) as DocumentEditorState,
    persist: {
        serializer: jsonSerializer,
    },
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
            this.closed = undefined;
            this.state = '';
            this.category = undefined;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useDocumentEditorStore as unknown as StoreDefinition, import.meta.hot));
}
