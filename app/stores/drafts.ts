import { defineStore } from 'pinia';

export const useDraftsStore = defineStore(
    'drafts',
    () => {
        // State
        const drafts = ref<Record<string, any>>({});

        // Actions
        function setDraft(key: string, value: any) {
            drafts.value[key] = value;
        }

        function getDraft(key: string) {
            return drafts.value[key];
        }

        function removeDraft(key: string) {
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
