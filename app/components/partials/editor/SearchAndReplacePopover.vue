<script lang="ts" setup>
import type { Editor, Range } from '@tiptap/core';

const props = defineProps<{
    editor: Editor;
    disabled: boolean;
}>();

const searchAndReplace = reactive<{
    search: string;
    replace: string;
    caseSensitive: boolean;
}>({
    search: '',
    replace: '',
    caseSensitive: false,
});

type SearchState = readonly [search: string, replace: string, caseSensitive: boolean];

const goToSelection = () => {
    if (!props.editor) return;

    const { results, resultIndex } = props.editor!.storage.searchAndReplace;
    const position: Range | undefined = results[resultIndex];

    if (!position) return;

    props.editor?.commands.setTextSelection(position);

    const { node } = props.editor!.view.domAtPos(props.editor!.state.selection.anchor);
    node instanceof HTMLElement && node.scrollIntoView({ behavior: 'smooth', block: 'center' });
};

watch(
    () => [searchAndReplace.search.trim(), searchAndReplace.replace, searchAndReplace.caseSensitive] as SearchState,
    ([search, replace, caseSensitive], previous) => {
        if (!props.editor) return;

        const [previousSearch, previousReplace, previousCaseSensitive] = previous ?? ['', '', false];
        const searchChanged = search !== previousSearch;
        const replaceChanged = replace !== previousReplace;
        const caseSensitiveChanged = caseSensitive !== previousCaseSensitive;

        if (!search) {
            if (previousSearch) {
                clear();
                return;
            }

            if (replaceChanged) props.editor.commands.setReplaceTerm(searchAndReplace.replace ?? '');
            if (caseSensitiveChanged) props.editor.commands.setCaseSensitive(searchAndReplace.caseSensitive);
            return;
        }

        if (searchChanged || caseSensitiveChanged) {
            props.editor.commands.resetIndex();
        }

        if (searchChanged) props.editor.commands.setSearchTerm(searchAndReplace.search);
        if (replaceChanged) props.editor.commands.setReplaceTerm(searchAndReplace.replace ?? '');
        if (caseSensitiveChanged) props.editor.commands.setCaseSensitive(searchAndReplace.caseSensitive);
    },
);

const replace = () => {
    props.editor?.commands.replace();
    goToSelection();
};

const next = () => {
    props.editor?.commands.nextSearchResult();
    goToSelection();
};

const previous = () => {
    props.editor?.commands.previousSearchResult();
    goToSelection();
};

const clear = () => {
    searchAndReplace.search = searchAndReplace.replace = '';
    props.editor?.commands.resetIndex();
};

const replaceAll = () => props.editor?.commands.replaceAll();
</script>

<template>
    <UPopover>
        <UTooltip :text="$t('components.partials.tiptap_editor.search_and_replace')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-text-search" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex gap-0.5 p-4">
                <UForm class="flex flex-col gap-2" :state="searchAndReplace">
                    <UFormField name="search" :label="$t('common.search')">
                        <UInput v-model="searchAndReplace.search" class="w-full" :disabled="disabled" />
                    </UFormField>

                    <UFormField name="replace" :label="$t('components.partials.tiptap_editor.replace')">
                        <UInput v-model="searchAndReplace.replace" class="w-full" :disabled="disabled" />
                    </UFormField>

                    <UFormField name="caseSensitive" :label="$t('common.case_sensitive')">
                        <USwitch v-model="searchAndReplace.caseSensitive" class="w-full" :disabled="disabled" />
                    </UFormField>

                    <UFormField class="flex flex-col lg:flex-row">
                        <UFieldGroup class="w-full">
                            <UButton
                                color="error"
                                variant="outline"
                                :label="$t('components.partials.tiptap_editor.clear')"
                                :disabled="disabled"
                                @click="() => clear()"
                            />
                            <UButton
                                color="neutral"
                                variant="outline"
                                :label="$t('components.partials.tiptap_editor.previous')"
                                :disabled="disabled"
                                @click="() => previous()"
                            />
                            <UButton
                                color="neutral"
                                variant="outline"
                                :label="$t('components.partials.tiptap_editor.next')"
                                :disabled="disabled"
                                @click="() => next()"
                            />
                            <UButton
                                color="neutral"
                                variant="outline"
                                :label="$t('components.partials.tiptap_editor.replace')"
                                :disabled="disabled"
                                @click="() => replace()"
                            />
                            <UButton
                                color="neutral"
                                variant="outline"
                                :label="$t('components.partials.tiptap_editor.replace_all')"
                                :disabled="disabled"
                                @click="
                                    () => {
                                        replaceAll();
                                    }
                                "
                            />
                        </UFieldGroup>

                        <div class="mt-1 block text-sm">
                            <span class="font-semibold">{{ $t('common.result', 2) }}</span
                            >:
                            {{
                                editor?.storage?.searchAndReplace?.results?.length > 0 &&
                                (editor?.storage?.searchAndReplace?.resultIndex ?? -1) >= 0
                                    ? editor?.storage?.searchAndReplace?.resultIndex + 1
                                    : 0
                            }}
                            /
                            {{ editor?.storage?.searchAndReplace?.results?.length }}
                        </div>
                    </UFormField>
                </UForm>
            </div>
        </template>
    </UPopover>
</template>
