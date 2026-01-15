<script lang="ts" setup>
import type { JSONContent } from '@tiptap/core';
import { z } from 'zod';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import type { HistoryContent } from '~/types/history';

const historyStore = useHistoryStore();

const schema = z.object({
    index: z.coerce.number().min(0),
    content: z.custom<JSONContent | string>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    index: 0,
    content: '',
});

const logger = useLogger('🗒️ Notepad');

const changed = ref(false);
const saving = ref(false);

const historyType = 'quickbutton-notepad' as const;
// Track last saved string and timestamp
let lastSavedString: JSONContent | string | undefined = undefined;
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema): Promise<void> {
    if (saving.value) return;

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) return;

    saving.value = true;

    historyStore.addVersion<HistoryContent>(historyType, state.index, {
        content: values.content,
        files: [],
    });

    useTimeoutFn(() => (saving.value = false), 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

onMounted(() => {
    logger.info('Notepad mounted, loading last version...');
    const lastVersion = historyStore.getLastVersion<HistoryContent>(historyType);
    if (lastVersion && lastVersion.content) {
        state.content = lastVersion.content.content;
    }
});

onBeforeUnmount(() => {
    logger.info('Notepad unmounting, force saving history...');
    saveHistory(state);
});

historyStore.handleRefresh(() => saveHistory(state));

watchDebounced(
    state,
    () => {
        if (changed.value) {
            saveHistory(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 1000,
        maxWait: 2500,
    },
);
</script>

<template>
    <UForm
        ref="formRef"
        :schema="schema"
        :state="state"
        class="flex min-h-full w-full max-w-full flex-1 flex-col overflow-y-auto"
    >
        <ClientOnly>
            <TiptapEditor
                v-model="state.content"
                name="content"
                class="mx-auto my-2 h-full w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
                :history-type="historyType"
                :saving="saving"
                disable-images
            />
        </ClientOnly>
    </UForm>
</template>
