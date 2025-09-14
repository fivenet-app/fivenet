<script setup lang="ts">
import { CodeDiff } from 'v-code-diff';
import { computed } from 'vue';
import type { Content, Version } from '~/types/history';

const props = defineProps<{
    currentContent: string; // String for diffing
    selectedVersion: Version<Content> | null;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'apply', version: Version<Content>): void;
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

const colorMode = useColorMode();

const selectedContent = computed(() => props.selectedVersion?.content.content ?? '');

const prettyCurrent = computed(() => prettyPrintHtml(props.currentContent));
const prettySelected = computed(() => prettyPrintHtml(selectedContent.value));

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
                    <CodeDiff
                        class="rounded border"
                        :old-string="prettyCurrent"
                        :new-string="prettySelected"
                        :context="3"
                        output-format="side-by-side"
                        language="text"
                        :theme="colorMode.value === 'dark' ? 'dark' : 'light'"
                    >
                        <!-- @vue-expect-error v-code-diff doesn't type the slot vars not even the slots currently -->
                        <template #stat="{ stat }">
                            <span class="diff-stat-added">+{{ stat.additionsNum }} additions</span>
                            <span class="diff-stat-deleted">-{{ stat.deletionsNum }} deletions</span>
                        </template>
                    </CodeDiff>
                </div>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="primary" @click="applySelected">{{ $t('common.apply') }}</UButton>
                <UButton class="flex-1" color="neutral" @click="$emit('close', false)">{{ $t('common.cancel') }}</UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
