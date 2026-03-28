import { mergeAttributes, type Extensions } from '@tiptap/core';
import { Blockquote } from '@tiptap/extension-blockquote';
import { Bold } from '@tiptap/extension-bold';
import { Code } from '@tiptap/extension-code';
import { CodeBlock } from '@tiptap/extension-code-block';
import { isChangeOrigin } from '@tiptap/extension-collaboration';
import { Details, DetailsContent, DetailsSummary } from '@tiptap/extension-details';
import Document from '@tiptap/extension-document';
import Emoji from '@tiptap/extension-emoji';
import { HardBreak } from '@tiptap/extension-hard-break';
import { Heading } from '@tiptap/extension-heading';
import Highlight from '@tiptap/extension-highlight';
import { HorizontalRule } from '@tiptap/extension-horizontal-rule';
import InvisibleCharacters from '@tiptap/extension-invisible-characters';
import Italic from '@tiptap/extension-italic';
import Link from '@tiptap/extension-link';
import { BulletList, ListItem, ListKeymap, OrderedList, TaskItem, TaskList } from '@tiptap/extension-list';
import NodeRange from '@tiptap/extension-node-range';
import { Paragraph } from '@tiptap/extension-paragraph';
import { Strike } from '@tiptap/extension-strike';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import { Table, TableCell, TableHeader, TableRow } from '@tiptap/extension-table';
import Text from '@tiptap/extension-text';
import TextAlign from '@tiptap/extension-text-align';
import { TextStyleKit } from '@tiptap/extension-text-style';
import Underline from '@tiptap/extension-underline';
import UniqueID from '@tiptap/extension-unique-id';
import { CharacterCount, Dropcursor, Focus, Gapcursor, Placeholder } from '@tiptap/extensions';
import { v4 as uuidv4 } from 'uuid';
import { CheckboxStandalone } from '~/composables/tiptap/extensions/CheckboxStandalone';
import SearchAndReplace from '~/composables/tiptap/extensions/SearchAndReplace';
import { EnhancedImage } from './tiptap/extensions/EnhancedImage';
import { PenaltyCalculator } from './tiptap/extensions/PenaltyCalculator';

export function useTiptapEditor(charLimit?: Ref<number>, placeholder?: Ref<string>): Extensions {
    const settingsStore = useSettingsStore();
    const { editor: editorSettings } = storeToRefs(settingsStore);

    const extensions: Extensions = [
        Blockquote,
        Bold,
        BulletList,
        CharacterCount.configure({
            limit: charLimit?.value ?? 0,
        }),
        CheckboxStandalone,
        Code.extend({
            excludes: 'code',
        }),
        CodeBlock,
        Details.configure({
            persist: true,
            HTMLAttributes: {
                class: 'details',
            },
        }),
        DetailsSummary,
        DetailsContent,
        Document,
        Dropcursor,
        Emoji,
        Focus.configure({
            className: 'has-focus',
            mode: 'all',
        }),
        Gapcursor,
        HardBreak,
        Heading,
        Highlight.configure({
            multicolor: true,
        }),
        HorizontalRule.extend({
            renderHTML() {
                return ['div', mergeAttributes(this.options.HTMLAttributes, { 'data-type': this.name }), ['hr']];
            },
        }),
        InvisibleCharacters.configure({
            visible: editorSettings.value.showInvisibleCharacters,
        }),
        Italic,
        Link.configure({
            openOnClick: false,
            defaultProtocol: 'https',
            HTMLAttributes: {
                target: null,
            },
        }),
        ListItem,
        ListKeymap,
        NodeRange.configure({
            depth: 0,
            key: null,
        }),
        OrderedList,
        Paragraph,
        Placeholder.configure({
            placeholder: () => placeholder?.value ?? '',
        }),
        SearchAndReplace,
        Strike,
        Subscript,
        Superscript,
        Table.configure({
            resizable: true,
            allowTableNodeSelection: true,
            HTMLAttributes: {
                class: 'border border-collapse border-solid border-neutral-500',
            },
        }),
        TableRow,
        TableHeader.configure({
            HTMLAttributes: {
                class: 'border border-solid border-neutral-600 bg-neutral-100 dark:bg-neutral-800',
            },
        }),
        TableCell.configure({
            HTMLAttributes: {
                class: 'border border-solid border-neutral-500',
            },
        }),
        TaskList,
        TaskItem.configure({
            nested: true,
        }),
        Text,
        TextAlign.configure({
            types: ['heading', 'paragraph', 'image'],
        }),
        TextStyleKit.configure({
            backgroundColor: {
                types: ['textStyle'],
            },
            color: {
                types: ['textStyle'],
            },
            fontFamily: {
                types: ['textStyle'],
            },
            fontSize: {
                types: ['textStyle'],
            },
            lineHeight: {
                types: ['textStyle'],
            },
        }),
        Underline,
        UniqueID.configure({
            attributeName: 'id',
            types: ['heading'],
            generateID: ({ node }) => `${node.type.name}-${uuidv4()}`,
            filterTransaction: (transaction) => !isChangeOrigin(transaction),
            updateDocument: true,
        }),
        // Custom
        EnhancedImage.configure({}),
        PenaltyCalculator.configure({}),
    ];

    return extensions;
}
