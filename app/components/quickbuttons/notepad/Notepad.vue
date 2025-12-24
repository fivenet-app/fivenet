<script lang="ts" setup>
import { z } from 'zod';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import type { Content } from '~/types/history';

const historyStore = useHistoryStore();

const schema = z.object({
    index: z.coerce.number().min(0),
    content: z.coerce.string().min(0).max(1750000),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    index: 0,
    content: '',
});

const content = ref('');

const changed = ref(false);
const saving = ref(false);

// Track last saved string and timestamp
let lastSavedString = '';
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, type = 'quickbutton-notepad'): Promise<void> {
    if (saving.value) return;

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) return;

    saving.value = true;

    historyStore.addVersion<Content>(type, state.index, {
        content: values.content,
        files: [],
    });

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

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
                v-model="content"
                name="content"
                class="mx-auto my-2 h-full w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
                history-type="quickbutton-notepad"
                :saving="saving"
                disable-images
            />
        </ClientOnly>
    </UForm>
</template>
