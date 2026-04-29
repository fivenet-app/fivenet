<script setup lang="ts">
import { ComparePlugin } from '@samemichaeltadele/tiptap-compare';
import type { JSONContent } from '@tiptap/core';
import type { HistoryContent, Version } from '~/types/history';
import TiptapEditor from './TiptapEditor.vue';

const props = defineProps<{
    // Current content for diff against selected version
    currentContent: JSONContent | undefined;
    selectedVersion: Version<HistoryContent> | null;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'apply', version: Version<HistoryContent>): void;
    (e: 'update:modelValue', value: boolean): void;
}>();

const selectedContent = computed(() => props.selectedVersion?.content.content ?? '');
watch(selectedContent, (newContent) => {
    editorRef.value?.editor?.commands.setComparisonContent(newContent);
});

function applySelected() {
    if (props.selectedVersion) {
        emit('apply', props.selectedVersion);
        emit('close', false);
    }
}

const editorRef = useTemplateRef('editorRef');
</script>

<template>
    <UModal :title="$t('common.version')" fullscreen>
        <template #body>
            <div class="space-y-4">
                <div>
                    <div class="mb-1 font-semibold">{{ $t('common.content') }}</div>

                    <TiptapEditor
                        ref="editorRef"
                        :model-value="props.currentContent"
                        disabled
                        disable-images
                        hide-toolbar
                        :extensions="[
                            ComparePlugin.configure({
                                comparisonContent: selectedContent,
                            }),
                        ]"
                    />
                </div>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" :label="$t('common.apply')" color="primary" @click="applySelected" />
                <UButton class="flex-1" :label="$t('common.cancel')" color="neutral" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
