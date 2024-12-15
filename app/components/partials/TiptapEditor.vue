<script lang="ts" setup>
import { Blockquote } from '@tiptap/extension-blockquote';
import { Bold } from '@tiptap/extension-bold';
import { BulletList } from '@tiptap/extension-bullet-list';
import CharacterCount from '@tiptap/extension-character-count';
import { Code } from '@tiptap/extension-code';
import { CodeBlock } from '@tiptap/extension-code-block';
import { Color } from '@tiptap/extension-color';
import { Document } from '@tiptap/extension-document';
import { Dropcursor } from '@tiptap/extension-dropcursor';
import FontFamily from '@tiptap/extension-font-family';
import { Gapcursor } from '@tiptap/extension-gapcursor';
import { HardBreak } from '@tiptap/extension-hard-break';
import { Heading } from '@tiptap/extension-heading';
import Highlight from '@tiptap/extension-highlight';
import { History } from '@tiptap/extension-history';
import { HorizontalRule } from '@tiptap/extension-horizontal-rule';
import { Italic } from '@tiptap/extension-italic';
import Link from '@tiptap/extension-link';
import { ListItem } from '@tiptap/extension-list-item';
import ListKeymap from '@tiptap/extension-list-keymap';
import { OrderedList } from '@tiptap/extension-ordered-list';
import { Paragraph } from '@tiptap/extension-paragraph';
import Placeholder from '@tiptap/extension-placeholder';
import { Strike } from '@tiptap/extension-strike';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import Table from '@tiptap/extension-table';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import TableRow from '@tiptap/extension-table-row';
import TaskItem from '@tiptap/extension-task-item';
import TaskList from '@tiptap/extension-task-list';
import { Text } from '@tiptap/extension-text';
import TextAlign from '@tiptap/extension-text-align';
import TextStyle from '@tiptap/extension-text-style';
import Underline from '@tiptap/extension-underline';
import { ImageResize } from '~/composables/tiptap/extensions/imageResize';
// @ts-expect-error project doesn't have types
import UniqueId from 'tiptap-unique-id';

import type { Range } from '@tiptap/core';
import SearchAndReplace from '~/composables/tiptap/extensions/searchAndReplace';

const props = withDefaults(
    defineProps<{
        modelValue: string;
        limit?: number;
        disabled?: boolean;
        placeholder?: string;
        splitScreen?: boolean;
        hideToolbar?: boolean;
        wrapperClass?: string;
    }>(),
    {
        limit: undefined,
        disabled: false,
        placeholder: undefined,
        splitScreen: false,
        hideToolbar: false,
        wrapperClass: '',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const content = useVModel(props, 'modelValue', emit);

const editor = useEditor({
    content: '',
    editorProps: {
        attributes: {
            class: 'prose prose-sm sm:prose-base lg:prose-lg xl:prose-2xl m-5 focus:outline-none dark:prose-invert max-w-full break-words',
        },
    },
    editable: !props.disabled,
    extensions: [
        UniqueId.configure({
            attributeName: 'id',
            types: ['heading'],
            createId: () => window.crypto.randomUUID(),
        }),
        // Basics
        Blockquote,
        Bold,
        BulletList,
        Code,
        CodeBlock,
        Color,
        Document,
        Dropcursor,
        FontFamily,
        Gapcursor,
        HardBreak,
        Heading,
        Highlight.configure({
            multicolor: true,
        }),
        History,
        HorizontalRule,
        Italic,
        Link.configure({
            openOnClick: false,
            defaultProtocol: 'https',
        }),
        ListItem,
        ListKeymap,
        OrderedList,
        Paragraph,
        Strike,
        Subscript,
        Superscript,
        Text,
        TextAlign.configure({
            types: ['heading', 'paragraph', 'image'],
        }),
        TextStyle,
        Underline,
        // Table
        Table.configure({
            resizable: true,
            allowTableNodeSelection: true,
            HTMLAttributes: {
                class: 'border border-collapse border-solid border-gray-500',
            },
        }),
        TableRow,
        TableHeader.configure({
            HTMLAttributes: {
                class: 'border border-solid border-gray-600 bg-gray-100 dark:bg-gray-800',
            },
        }),
        TableCell.configure({
            HTMLAttributes: {
                class: 'border border-solid border-gray-500',
            },
        }),
        // Misc
        ImageResize.configure({
            inline: true,
            allowBase64: true,
        }),
        SearchAndReplace,
        TaskList,
        TaskItem.configure({
            nested: true,
        }),
        CharacterCount.configure({
            limit: props.limit,
        }),
        Placeholder.configure({
            placeholder: props.placeholder ?? '',
        }),
    ],
    onUpdate: () => {
        emit('update:modelValue', unref(editor)?.getHTML() ?? '');
    },
});

const fonts = [
    {
        label: 'Default',
        value: 'DM Sans',
    },
    {
        label: 'Arial',
        value: 'arial, helvetica, sans-serif',
    },
    {
        label: 'Serif',
        value: 'serif',
    },
    {
        label: 'Times New Romain',
        value: 'times new roman, times, serif',
    },
    {
        label: 'Comic Sans',
        value: 'Comic Sans MS, Comic Sans',
    },
    {
        label: 'Monospace',
        value: 'monospace',
    },
];

const fontColors = [
    {
        label: 'White',
        value: '#ffffff',
    },
    {
        label: 'Black',
        value: '#000000',
    },
    {
        label: 'Red',
        value: '#F98181',
    },
    {
        label: 'Orange',
        value: '#FBBC88',
    },
    {
        label: 'Yellow',
        value: '#FAF594',
    },
    {
        label: 'Blue',
        value: '#70CFF8',
    },
    {
        label: 'Teal',
        value: '#94FADB',
    },
    {
        label: 'Green',
        value: '#B9F18D',
    },
    {
        label: 'Purple',
        value: '#958DF1',
    },
];

const highlightColors = [
    {
        label: 'Red',
        value: '#ffa8a8',
    },
    {
        label: 'Orange',
        value: '#ffc078',
    },
    {
        label: 'Yellow',
        value: '#FAF594',
    },
    {
        label: 'Blue',
        value: '#74c0fc',
    },
    {
        label: 'Green',
        value: '#8ce99a',
    },
    {
        label: 'Purple',
        value: '#b197fc',
    },
];

watch(content, (value) => {
    const isSame = unref(editor)?.getHTML() === value;
    // JSON
    // const isSame = JSON.stringify(this.editor.getJSON()) === JSON.stringify(value)

    if (isSame) {
        return;
    }

    unref(editor)?.commands.setContent(value, false);
});

watch(
    () => props.disabled,
    () => unref(editor)?.setEditable(!props.disabled),
);

const selectedFont = ref<(typeof fonts)[0]>(fonts[0]!);
watch(selectedFont, () => unref(editor)?.chain().focus().setFontFamily(selectedFont.value.value).run());

const selectedFontColor = ref<(typeof fontColors)[0]>(fontColors[0]!);
watch(selectedFontColor, () => unref(editor)?.chain().focus().setColor(selectedFontColor.value.value).run());

const selectedHighlightColor = ref<(typeof highlightColors)[0]>(highlightColors[0]!);
watch(selectedHighlightColor, () =>
    unref(editor)?.chain().focus().toggleHighlight({ color: selectedHighlightColor.value.value }).run(),
);

const searchAndReplace = reactive<{
    search: string;
    replace: string;
    caseSensitive: boolean;
}>({
    search: '',
    replace: '',
    caseSensitive: false,
});

const updateSearchReplace = (clearIndex: boolean = false) => {
    if (!editor.value) {
        return;
    }

    if (clearIndex) {
        editor.value.commands.resetIndex();
    }

    unref(editor)?.commands.setSearchTerm(searchAndReplace.search);
    unref(editor)?.commands.setReplaceTerm(searchAndReplace.replace ?? '');
    unref(editor)?.commands.setCaseSensitive(searchAndReplace.caseSensitive);
};

const goToSelection = () => {
    if (!editor.value) {
        return;
    }

    const { results, resultIndex } = editor.value.storage.searchAndReplace;
    const position: Range = results[resultIndex];

    if (!position) {
        return;
    }

    unref(editor)?.commands.setTextSelection(position);

    const { node } = editor.value.view.domAtPos(editor.value.state.selection.anchor);
    node instanceof HTMLElement && node.scrollIntoView({ behavior: 'smooth', block: 'center' });
};

watch(
    () => searchAndReplace.search.trim(),
    (val, oldVal) => {
        if (!val) clear();
        if (val !== oldVal) updateSearchReplace(true);
    },
);

watch(
    () => searchAndReplace.replace.trim(),
    (val, oldVal) => (val === oldVal ? null : updateSearchReplace()),
);

watch(
    () => searchAndReplace.caseSensitive,
    (val, oldVal) => (val === oldVal ? null : updateSearchReplace(true)),
);

const replace = () => {
    editor.value?.commands.replace();
    goToSelection();
};

const next = () => {
    editor.value?.commands.nextSearchResult();
    goToSelection();
};

const previous = () => {
    editor.value?.commands.previousSearchResult();
    goToSelection();
};

const clear = () => {
    searchAndReplace.search = searchAndReplace.replace = '';
    editor.value?.commands.resetIndex();
};

const replaceAll = () => editor.value?.commands.replaceAll();

onMounted(() => {
    if (unref(editor)) {
        unref(editor)?.commands.setContent(props.modelValue);
    }
});

onBeforeUnmount(() => {
    unref(editor)?.destroy();
});
</script>

<template>
    <div class="rounded border border-gray-100 dark:border-gray-800">
        <div v-if="editor && !hideToolbar" class="flex snap-x flex-wrap gap-1 bg-gray-100 p-0.5 dark:bg-gray-800">
            <UButtonGroup>
                <UButton
                    :disabled="!editor.can().chain().focus().toggleBold().run()"
                    :class="{ 'is-active': editor.isActive('bold') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-bold"
                    @click="editor.chain().focus().toggleBold().run()"
                />
                <UButton
                    :disabled="!editor.can().chain().focus().toggleItalic().run()"
                    :class="{ 'is-active': editor.isActive('italic') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-italic"
                    @click="editor.chain().focus().toggleItalic().run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('underline') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-underline"
                    @click="editor.chain().focus().toggleUnderline().run()"
                />
                <UButton
                    :disabled="!editor.can().chain().focus().toggleStrike().run()"
                    :class="{ 'is-active': editor.isActive('strike') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-strikethrough"
                    @click="editor.chain().focus().toggleStrike().run()"
                />
                <UButton
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-clear"
                    @click="editor.chain().focus().unsetAllMarks().run()"
                />
                <!-- <UButton color="white" variant="ghost" icon="i-mdi-backspace" @click="editor.chain().focus().clearNodes().run()">clear nodes</UButton> -->
            </UButtonGroup>

            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive('superscript') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-superscript"
                    @click="editor.chain().focus().toggleSuperscript().run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('subscript') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-subscript"
                    @click="editor.chain().focus().toggleSubscript().run()"
                />
                <UButton
                    :disabled="!editor.can().chain().focus().toggleCode().run()"
                    :class="{ 'is-active': editor.isActive('code') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-code-braces"
                    @click="editor.chain().focus().toggleCode().run()"
                />
            </UButtonGroup>

            <!-- Text Align -->
            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-align-left"
                    @click="editor.chain().focus().setTextAlign('left').run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-align-center"
                    @click="editor.chain().focus().setTextAlign('center').run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-align-right"
                    @click="editor.chain().focus().setTextAlign('right').run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive({ textAlign: 'justify' }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-align-justify"
                    @click="editor.chain().focus().setTextAlign('justify').run()"
                />
            </UButtonGroup>

            <UButtonGroup>
                <!-- Font Family -->
                <UInputMenu
                    v-model="selectedFont"
                    option-attribute="label"
                    :search-attributes="['label']"
                    :options="fonts"
                    :placeholder="$t('common.font', 1)"
                    search-lazy
                    :search-placeholder="$t('common.search_field')"
                >
                    <template #label>
                        <span class="truncate" :style="{ fontFamily: selectedFont.value }">{{ selectedFont.label }}</span>
                    </template>

                    <template #option="{ option }">
                        <span class="truncate" :style="{ fontFamily: option.value }">{{ option.label }}</span>
                    </template>

                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>>
                </UInputMenu>

                <UPopover>
                    <UButton
                        :class="{ 'is-active': editor.isActive('color', { color: selectedFontColor.value }) }"
                        color="white"
                        variant="ghost"
                        :style="{ color: selectedFontColor.value }"
                        icon="i-mdi-format-color-text"
                    />

                    <template #panel>
                        <div class="grid grid-cols-6 gap-0.5 p-4">
                            <UButton
                                v-for="(color, idx) in fontColors"
                                :key="idx"
                                class="size-6 rounded-none border-0"
                                :style="{ backgroundColor: color.value }"
                                @click="selectedFontColor = color"
                            />
                        </div>
                    </template>
                </UPopover>

                <UButton
                    :class="{ 'is-active': editor.isActive('paragraph') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-paragraph"
                    @click="editor.chain().focus().setParagraph().run()"
                />

                <!-- Headers -->
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-1"
                    @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-2"
                    @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-3"
                    @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 4 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-4"
                    @click="editor.chain().focus().toggleHeading({ level: 4 }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 5 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-5"
                    @click="editor.chain().focus().toggleHeading({ level: 5 }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('heading', { level: 6 }) }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-header-6"
                    @click="editor.chain().focus().toggleHeading({ level: 6 }).run()"
                />
            </UButtonGroup>

            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive('highlight') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-color-highlight"
                    @click="editor.chain().focus().toggleHighlight().run()"
                />

                <UPopover>
                    <UButton
                        :class="{ 'is-active': editor.isActive('highlight', { color: selectedHighlightColor.value }) }"
                        color="white"
                        variant="ghost"
                        :style="{ color: selectedHighlightColor.value }"
                        icon="i-mdi-format-color-fill"
                    />

                    <template #panel>
                        <div class="grid grid-cols-6 gap-0.5 p-4">
                            <UButton
                                v-for="(color, idx) in highlightColors"
                                :key="idx"
                                class="size-6 rounded-none border-0"
                                :style="{ backgroundColor: color.value }"
                                @click="selectedHighlightColor = color"
                            />
                        </div>
                    </template>
                </UPopover>
            </UButtonGroup>

            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive('bulletList') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-list-bulleted"
                    @click="editor.chain().focus().toggleBulletList().run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('orderedList') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-list-numbered"
                    @click="editor.chain().focus().toggleOrderedList().run()"
                />
                <UButton
                    icon="i-mdi-format-list-checks"
                    color="white"
                    variant="ghost"
                    :class="{ 'is-active': editor.isActive('taskList') }"
                    @click="editor.chain().focus().toggleTaskList().run()"
                />
            </UButtonGroup>

            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive('table') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-table"
                    @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('codeBlock') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-code-block-braces"
                    @click="editor.chain().focus().toggleCodeBlock().run()"
                />
                <UButton
                    :class="{ 'is-active': editor.isActive('blockquote') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-quote-open"
                    @click="editor.chain().focus().toggleBlockquote().run()"
                />
                <UButton
                    color="white"
                    variant="ghost"
                    icon="i-mdi-minus"
                    @click="editor.chain().focus().setHorizontalRule().run()"
                />
                <UButton
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-page-break"
                    @click="editor.chain().focus().setHardBreak().run()"
                />
            </UButtonGroup>

            <UPopover>
                <UButton color="white" variant="ghost" icon="i-mdi-text-search" />

                <template #panel>
                    <div class="flex flex-1 gap-0.5 p-4">
                        <UForm :state="searchAndReplace">
                            <UFormGroup name="search" :label="$t('common.search')">
                                <UInput v-model="searchAndReplace.search" />
                            </UFormGroup>

                            <UFormGroup name="replace" :label="$t('components.partials.TipTapEditor.replace')">
                                <UInput v-model="searchAndReplace.replace" />
                            </UFormGroup>

                            <UFormGroup name="caseSensitive" :label="$t('common.case_sensitive')">
                                <UToggle v-model="searchAndReplace.caseSensitive" />
                            </UFormGroup>

                            <UFormGroup class="flex flex-col lg:flex-row">
                                <UButtonGroup>
                                    <UButton
                                        color="white"
                                        variant="outline"
                                        :label="$t('components.partials.TipTapEditor.clear')"
                                        @click="clear"
                                    />
                                    <UButton
                                        color="white"
                                        variant="outline"
                                        :label="$t('components.partials.TipTapEditor.previous')"
                                        @click="previous"
                                    />
                                    <UButton
                                        color="white"
                                        variant="outline"
                                        :label="$t('components.partials.TipTapEditor.next')"
                                        @click="next"
                                    />
                                    <UButton
                                        color="white"
                                        variant="outline"
                                        :label="$t('components.partials.TipTapEditor.replace')"
                                        @click="replace"
                                    />
                                    <UButton
                                        color="white"
                                        variant="outline"
                                        :label="$t('components.partials.TipTapEditor.replace_all')"
                                        @click="replaceAll"
                                    />
                                </UButtonGroup>

                                <div class="mt-1 block text-sm font-medium">
                                    {{ $t('common.result', 2) }}:
                                    {{
                                        editor?.storage?.searchAndReplace?.resultIndex > 0
                                            ? editor?.storage?.searchAndReplace?.resultIndex + 1
                                            : 0
                                    }}
                                    /
                                    {{ editor?.storage?.searchAndReplace?.results.length }}
                                </div>
                            </UFormGroup>
                        </UForm>
                    </div>
                </template>
            </UPopover>

            <UButtonGroup>
                <UButton
                    :disabled="!editor.can().chain().focus().undo().run()"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-undo"
                    @click="editor.chain().focus().undo().run()"
                />
                <UButton
                    :disabled="!editor.can().chain().focus().redo().run()"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-redo"
                    @click="editor.chain().focus().redo().run()"
                />
            </UButtonGroup>
        </div>

        <TiptapEditorContent :editor="editor" class="w-full min-w-0 max-w-full" :class="wrapperClass" />
        <UTextarea
            v-if="splitScreen"
            v-model="content"
            :rows="24"
            autoresize
            class="border-t-2 border-gray-500"
            :ui="{ rounded: '' }"
        />

        <div v-if="editor" class="flex w-full justify-between bg-gray-100 px-1 text-center dark:bg-gray-800">
            <div class="flex-1">
                <slot name="footer" />
            </div>

            <div>
                {{ editor.storage.characterCount.characters() }}<template v-if="limit && limit > 0"> / {{ limit }}</template>
                {{ $t('common.chars', editor.storage.characterCount.characters()) }}
                |
                {{ editor.storage.characterCount.words() }} {{ $t('common.word', editor.storage.characterCount.words()) }}
            </div>
        </div>
    </div>
</template>

<style lang="scss">
/* Basic editor styles */
.tiptap {
    :first-child {
        margin-top: 0;
    }

    /* Task list specific styles */
    ul[data-type='taskList'] {
        list-style: none;
        margin-left: 0;
        padding: 0;

        li {
            align-items: flex-start;
            display: flex;

            > label {
                flex: 0 0 auto;
                margin-right: 0.5rem;
                user-select: none;
            }

            > div {
                flex: 1 1 auto;
            }
        }

        input[type='checkbox'] {
            cursor: pointer;
        }

        ul[data-type='taskList'] {
            margin: 0;
        }
    }

    /* Search And Replace */
    .search-result {
        background-color: rgba(255, 217, 0, 0.5);

        &-current {
            background-color: rgba(13, 255, 0, 0.5);
        }
    }

    /* Table-specific styling */
    table {
        display: table;
        border-collapse: collapse;
        margin: 0;
        overflow: hidden;
        table-layout: fixed;
        width: 100%;

        td,
        th {
            box-sizing: border-box;
            min-width: 1em;
            padding: 6px 8px;
            position: relative;
            vertical-align: top;

            > * {
                margin-bottom: 0;
            }
        }

        th {
            font-weight: bold;
            text-align: left;
        }

        .selectedCell:after {
            content: '';
            left: 0;
            right: 0;
            top: 0;
            bottom: 0;
            pointer-events: none;
            position: absolute;
            z-index: 2;
        }

        .column-resize-handle {
            background-color: var(--color-primary-500);
            bottom: -2px;
            pointer-events: none;
            position: absolute;
            right: -2px;
            top: 0;
            width: 4px;
        }
    }

    .tableWrapper {
        margin: 1.5rem 0;
        overflow-x: auto;
    }

    &.resize-cursor {
        cursor: ew-resize;
        cursor: col-resize;
    }
}
</style>
