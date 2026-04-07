<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { Editor } from '@tiptap/core';
import { z } from 'zod';

const props = defineProps<{
    editor: Editor;
    active?: boolean | undefined;
    disabled?: boolean | undefined;
}>();

const schema = z.object({
    rows: z.coerce.number().min(1).max(25).default(2),
    cols: z.coerce.number().min(1).max(200).default(3),
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

    isOpen.value = false;
}

const isOpen = ref(false);
</script>

<template>
    <UPopover v-model:open="isOpen">
        <UTooltip :text="$t('components.partials.tiptap_editor.table')">
            <UButton
                :class="{ 'bg-neutral-300 dark:bg-neutral-900': active }"
                color="neutral"
                variant="ghost"
                icon="i-mdi-table"
                :disabled="disabled"
            />
        </UTooltip>

        <template #content>
            <div class="p-4">
                <div v-if="active" class="grid grid-cols-1 gap-2">
                    <div class="grid grid-cols-2 gap-2">
                        <UFieldGroup orientation="vertical">
                            <UButton
                                class="w-full"
                                :label="$t('components.partials.tiptap_editor.add_column_before')"
                                icon="i-mdi-table-column-plus-before"
                                @click="editor.chain().focus().addColumnBefore().run()"
                            />
                            <UButton
                                class="w-full"
                                :label="$t('components.partials.tiptap_editor.add_column_after')"
                                icon="i-mdi-table-column-plus-after"
                                @click="editor.chain().focus().addColumnAfter().run()"
                            />
                        </UFieldGroup>

                        <UFieldGroup orientation="vertical">
                            <UButton
                                class="w-full"
                                :label="$t('components.partials.tiptap_editor.add_row_before')"
                                icon="i-mdi-table-row-plus-before"
                                @click="editor.chain().focus().addRowBefore().run()"
                            />
                            <UButton
                                class="w-full"
                                :label="$t('components.partials.tiptap_editor.add_row_after')"
                                icon="i-mdi-table-row-plus-after"
                                @click="editor.chain().focus().addRowAfter().run()"
                            />
                        </UFieldGroup>
                    </div>

                    <USeparator />

                    <UFieldGroup>
                        <UButton
                            class="w-full"
                            :label="$t('components.partials.tiptap_editor.delete_column')"
                            icon="i-mdi-table-column-remove"
                            color="red"
                            variant="subtle"
                            @click="editor.chain().focus().deleteColumn().run()"
                        />
                        <UButton
                            class="w-full"
                            :label="$t('components.partials.tiptap_editor.delete_row')"
                            icon="i-mdi-table-row-remove"
                            color="red"
                            variant="subtle"
                            @click="editor.chain().focus().deleteRow().run()"
                        />
                    </UFieldGroup>
                </div>

                <UForm v-else :schema="schema" :state="state" @submit="onSubmit">
                    <UFormField :label="$t('components.partials.tiptap_editor.rows')" name="rows">
                        <UInput v-model="state.rows" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.cols')" name="cols">
                        <UInput v-model="state.cols" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.with_header_row')" name="withHeaderRow">
                        <USwitch v-model="state.withHeaderRow" type="text" :disabled="disabled" />
                    </UFormField>

                    <UFormField class="mt-2">
                        <UButton
                            type="submit"
                            icon="i-mdi-table-plus"
                            class="w-full"
                            :label="$t('common.create')"
                            :disabled="disabled"
                        />
                    </UFormField>
                </UForm>
            </div>
        </template>
    </UPopover>
</template>
