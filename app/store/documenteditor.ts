import { defineStore } from 'pinia';
import type { Category } from '~~/gen/ts/resources/documents/category';

export interface DocumentEditorState {
    title: string;
    content: string;
    closed: boolean;
    state: string;
    category?: Category;
}

export const useDocumentEditorStore = defineStore(
    'documentEditor',
    () => {
        // State
        const title = ref<string>('');
        const content = ref<string>('');
        const closed = ref<boolean>(false);
        const state = ref<string>('');
        const category = ref<Category | undefined>(undefined);

        // Actions
        const save = (doc: DocumentEditorState): void => {
            title.value = doc.title;
            content.value = doc.content;
            closed.value = doc.closed;
            state.value = doc.state;
            category.value = doc.category;
        };

        const clear = (): void => {
            title.value = '';
            content.value = '';
            closed.value = false;
            state.value = '';
            category.value = undefined;
        };

        return {
            title,
            content,
            closed,
            state,
            category,
            save,
            clear,
        };
    },
    {
        persist: true,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useDocumentEditorStore, import.meta.hot));
}
