<script setup lang="ts">
import { LazyPartialsCodeDiff } from '#components';
import type { JSONContent } from '@tiptap/core';
import { renderToHTMLString } from '@tiptap/static-renderer';
import type { HistoryContent, Version } from '~/types/history';

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

function prettyPrintHtml(html: string): string {
    // Insert newlines between tags, but not for pairs like <p>...</p>
    const formatted = html.replace(/>\s*</g, '>\n<');
    const lines: string[] = [];
    let buffer = '';

    formatted.split('\n').forEach((line) => {
        // If this line is an opening tag, content, and closing tag all on one line,
        // keep it together
        if (line.match(/^<([a-zA-Z0-9]+)(\s[^>]*)?>.*<\/\1>$/)) {
            lines.push(buffer + line.trim());
            buffer = '';
        } else if (line.trim() === '') {
            // skip empty lines
        } else if (line.startsWith('<') && !line.startsWith('</')) {
            // opening tag
            lines.push(buffer + line.trim());
            buffer = '';
        } else if (line.startsWith('</')) {
            // closing tag
            lines.push(buffer + line.trim());
            buffer = '';
        } else {
            // content between tags
            buffer += line.trim();
        }
    });
    // Now indent
    let indent = 0;
    return lines
        .map((line) => {
            if (line.match(/^<\/\w/)) indent--;
            const pad = '  '.repeat(Math.max(indent, 0));
            if (
                line.match(/^<\w([^>]*)[^/]>.*$/) &&
                !line.match(/^<.*\/>$/) && // Not self-closing
                !line.match(/^<([a-zA-Z0-9]+)(\s[^>]*)?>.*<\/\1>$/) // Not open+close same line
            )
                indent++;
            return pad + line;
        })
        .join('\n');
}

const extensions = useTiptapEditor();

const colorMode = useColorMode();

const prettyCurrent = computed(() =>
    props.currentContent ? prettyPrintHtml(renderToHTMLString({ content: props.currentContent, extensions: extensions })) : '',
);

const selectedContent = computed(() => props.selectedVersion?.content.content ?? '');

const prettySelected = computed(() =>
    prettyPrintHtml(
        typeof selectedContent.value === 'string'
            ? selectedContent.value
            : renderToHTMLString({ content: selectedContent.value, extensions: extensions }),
    ),
);

function applySelected() {
    if (props.selectedVersion) {
        emit('apply', props.selectedVersion);
        emit('close', false);
    }
}
</script>

<template>
    <UModal :title="$t('common.version')" fullscreen>
        <template #body>
            <div class="space-y-4">
                <div>
                    <div class="mb-1 font-semibold">{{ $t('common.content') }}</div>
                    <LazyPartialsCodeDiff
                        class="rounded border"
                        :old-string="prettyCurrent"
                        :new-string="prettySelected"
                        :context="3"
                        output-format="side-by-side"
                        language="text"
                        :theme="colorMode.value === 'dark' ? 'dark' : 'light'"
                    >
                        <template #stat="{ stat }">
                            <span class="diff-stat-added">+{{ stat.additionsNum }} {{ $t('common.additions') }}</span>
                            <span class="diff-stat-deleted">-{{ stat.deletionsNum }} {{ $t('common.deletions') }}</span>
                        </template>
                    </LazyPartialsCodeDiff>
                </div>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton :label="$t('common.apply')" class="flex-1" color="primary" @click="applySelected" />
                <UButton :label="$t('common.cancel')" class="flex-1" color="neutral" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
