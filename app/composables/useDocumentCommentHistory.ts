import type { JSONContent } from '@tiptap/core';
import type { MaybeRefOrGetter } from 'vue';
import type { HistoryContent } from '~/types/history';
import { isEmptyDoc } from '~/utils/tiptap';

export type DocumentCommentEditorState = {
    content?: JSONContent | string;
};

export function useDocumentCommentHistory(
    state: DocumentCommentEditorState,
    documentId: MaybeRefOrGetter<number>,
    options: { dirty?: MaybeRefOrGetter<boolean> } = {},
) {
    const { t } = useI18n();
    const historyStore = useHistoryStore();

    const serializeCommentContent = (value: DocumentCommentEditorState): string => {
        if (isEmptyDoc(value.content)) return '__empty__';
        return JSON.stringify(value.content);
    };

    const { snapshotDirty, hasUnsavedChanges, syncSnapshot } = useSnapshotChanges(state, {
        dirty: options.dirty,
        serializer: serializeCommentContent,
    });

    const saving = ref(false);
    let lastSavedContent: JSONContent | string | undefined;
    let lastSaveTimestamp = 0;

    async function saveHistory(values: DocumentCommentEditorState, type = 'document_comments'): Promise<void> {
        if (saving.value) return;

        const now = Date.now();
        if (values.content === lastSavedContent || now - lastSaveTimestamp < 5000) return;

        saving.value = true;

        historyStore.addVersion<HistoryContent>(
            type,
            toValue(documentId),
            {
                content: values.content,
                files: [],
            },
            `${t('common.comment')}: DOC-${toValue(documentId)}`,
        );

        useTimeoutFn(() => {
            saving.value = false;
        }, 1750);

        lastSavedContent = values.content;
        lastSaveTimestamp = now;
    }

    historyStore.handleRefresh(() => saveHistory(state));

    return {
        hasUnsavedChanges,
        snapshotDirty,
        saving,
        syncSnapshot,
        saveHistory,
    };
}
