import { defineStore } from 'pinia';

export const useDraftsStore = defineStore(
    'drafts',
    () => {
        // State
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const drafts = ref<Record<string, any>>({});

        // Actions
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        function setDraft(key: string, value: any) {
            drafts.value[key] = value;
        }

        function getDraft(key: string) {
            return drafts.value[key];
        }

        function removeDraft(key: string) {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete drafts.value[key];
        }

        function clearDrafts() {
            drafts.value = {};
        }

        return {
            drafts,
            setDraft,
            getDraft,
            removeDraft,
            clearDrafts,
        };
    },
    {
        persist: true,
    },
);
