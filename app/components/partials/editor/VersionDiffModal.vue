<script setup lang="ts">
import { CodeDiff } from 'v-code-diff';
import { computed } from 'vue';
import type { Content, Version } from '~/types/history';

const props = defineProps<{
    modelValue: boolean;
    currentContent: string; // String for diffing
    selectedVersion: Version<Content> | null;
}>();

const emit = defineEmits<{
    (e: 'apply', version: Version<Content>): void;
    (e: 'update:modelValue', value: boolean): void;
}>();

const open = computed({
    get: () => props.modelValue,
    set: (v) => emit('update:modelValue', v),
});

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
        open.value = false;
    }
}
function close() {
    open.value = false;
}
</script>

<template>
    <UCard
        class="flex flex-1 flex-col"
        :ui="{
            body: {
                base: 'flex-1 min-h-[calc(100dvh-(2*var(--ui-header-height)))] max-h-[calc(100dvh-(2*var(--ui-header-height)))] overflow-y-auto',
                padding: 'px-1 py-2 sm:p-2',
            },
        }"
    >
        <template #header>
            <div class="flex items-center justify-between">
                <h3 class="text-2xl font-semibold leading-6">
                    {{ $t('common.version_history') }}
                </h3>

                <UButton class="-my-1" color="neutral" variant="ghost" icon="i-mdi-window-close" @click="close" />
            </div>
        </template>

        <div>
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
                        <template #stat="{ stat }">
                            <span class="diff-stat-added">+{{ stat.additionsNum }} additions</span>
                            <span class="diff-stat-deleted">-{{ stat.deletionsNum }} deletions</span>
                        </template>
                    </CodeDiff>
                </div>
            </div>
        </div>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="primary" @click="applySelected">{{ $t('common.apply') }}</UButton>
                <UButton class="flex-1" color="neutral" @click="close">{{ $t('common.cancel') }}</UButton>
            </UButtonGroup>
        </template>
    </UCard>
</template>
