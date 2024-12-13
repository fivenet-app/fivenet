<script lang="ts" setup>
import { Blockquote } from '@tiptap/extension-blockquote';
import { Bold } from '@tiptap/extension-bold';
import { BulletList } from '@tiptap/extension-bullet-list';
import CharacterCount from '@tiptap/extension-character-count';
import { Code } from '@tiptap/extension-code';
import { CodeBlock } from '@tiptap/extension-code-block';
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
import ImageResize from 'tiptap-extension-resize-image';
// @ts-expect-error project doesn't have types
import UniqueId from 'tiptap-unique-id';

const props = defineProps<{
    modelValue: string;
    limit?: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const editor = useEditor({
    content: '',
    editorProps: {
        attributes: {
            class: 'prose prose-sm sm:prose-base lg:prose-lg xl:prose-2xl m-5 focus:outline-none dark:prose-invert',
        },
    },
    extensions: [
        UniqueId.configure({
            attributeName: 'id',
            types: ['heading'],
            createId: () => window.crypto.randomUUID(),
        }),
        // Starter Kit
        Blockquote,
        Bold,
        BulletList,
        Code,
        CodeBlock,
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
        TaskList,
        TaskItem.configure({
            nested: true,
        }),
        CharacterCount.configure({
            limit: props.limit,
        }),
        Placeholder.configure({
            placeholder: 'Write something …',
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
        label: '',
        value: '',
    },
];

const highlightColors = [
    {
        label: '',
        value: '',
    },
];

watch(
    () => props.modelValue,
    (value) => {
        const isSame = unref(editor)?.getHTML() === value;

        // JSON
        // const isSame = JSON.stringify(this.editor.getJSON()) === JSON.stringify(value)

        if (isSame) {
            return;
        }

        unref(editor)?.commands.setContent(value, false);
    },
);

const selectedFont = ref<(typeof fonts)[0]>(fonts[0]!);
watch(selectedFont, () => unref(editor)?.chain().focus().setFontFamily(selectedFont.value.value).run());

const selectedHighlightColor = ref<string>('yellow');
watch(selectedHighlightColor, () =>
    unref(editor)?.chain().focus().toggleHighlight({ color: selectedHighlightColor.value }).run(),
);

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
    <div>
        <div v-if="editor" class="flex snap-x flex-wrap gap-1 bg-gray-100 dark:bg-gray-800">
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
            <UButton
                color="white"
                variant="ghost"
                icon="i-mdi-format-clear"
                @click="editor.chain().focus().unsetAllMarks().run()"
            />
            <!-- <UButton color="white" variant="ghost" icon="i-mdi-backspace" @click="editor.chain().focus().clearNodes().run()">clear nodes</UButton> -->

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

            <UButtonGroup>
                <UButton
                    :class="{ 'is-active': editor.isActive('highlight') }"
                    color="white"
                    variant="ghost"
                    icon="i-mdi-format-color-highlight"
                    @click="editor.chain().focus().toggleHighlight().run()"
                />

                <UPopover>
                    <UButton :class="{ 'is-active': editor.isActive('highlight', { color: '#ffc078' }) }"> Orange </UButton>

                    <template #panel>
                        <div class="p-4">PLACEHOLDER</div>
                    </template>
                </UPopover>
            </UButtonGroup>

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
            <UButton
                icon="i-mdi-format-list-checks"
                color="white"
                variant="ghost"
                :class="{ 'is-active': editor.isActive('taskList') }"
                @click="editor.chain().focus().toggleTaskList().run()"
            />
        </div>

        <TiptapEditorContent :editor="editor" />

        <div v-if="editor" class="w-full text-center">
            {{ editor.storage.characterCount.characters() }}<template v-if="limit && limit > 0"> / {{ limit }}</template>
            {{ $t('common.chars', editor.storage.characterCount.characters()) }}
            |
            {{ editor.storage.characterCount.words() }} {{ $t('common.word', editor.storage.characterCount.words()) }}
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
