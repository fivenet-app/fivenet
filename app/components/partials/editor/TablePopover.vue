<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { Editor } from '@tiptap/vue-3';
import { z } from 'zod';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean | undefined;
}>();

const schema = z.object({
    rows: z.number().min(1).max(25).default(2),
    cols: z.number().min(1).max(200).default(3),
    withHeaderRow: z.boolean().default(true),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    rows: 2,
    cols: 3,
    withHeaderRow: true,
});

async function onSubmit(event: FormSubmitEvent<Schema>): Promise<void> {
    const values = event.data;

    unref(props.editor)
        ?.chain()
        .focus()
        .insertTable({
            rows: values.rows,
            cols: values.cols,
            withHeaderRow: values.withHeaderRow,
        })
        .run();
}
</script>

<template>
    <UPopover>
        <UTooltip :text="$t('components.partials.TiptapEditor.table')">
            <UButton
                :class="{ 'bg-gray-300 dark:bg-gray-900': editor.isActive('table') }"
                color="neutral"
                variant="ghost"
                icon="i-mdi-table"
                :disabled="disabled"
            />
        </UTooltip>

        <template #content>
            <div class="p-4">
                <UForm :schema="schema" :state="{}" @submit="onSubmit">
                    <UFormField :label="$t('common.rows')">
                        <UInput v-model="state.rows" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField :label="$t('common.cols')">
                        <UInput v-model="state.cols" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField :label="$t('common.with_header_row')">
                        <USwitch v-model="state.withHeaderRow" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField>
                        <UButton type="submit" :label="$t('common.create')" :disabled="disabled" />
                    </UFormField>
                </UForm>
            </div>
        </template>
    </UPopover>
</template>
