import { Node, mergeAttributes } from '@tiptap/core';

export interface TemplateBlockOptions {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    HTMLAttributes: Record<string, any>;
}

declare module '@tiptap/core' {
    /**
     * Commands for managing template blocks.
     */
    interface Commands<ReturnType> {
        templateBlock: {
            /**
             * Represents a set of commands for managing template blocks within the editor.
             *
             *
             * @param payload - The data required to insert the template block.
             * @param payload.value - The string value of the template block to be inserted.
             * @param payload.leftTrim - Optional. Whether to trim whitespace from the left side of the variable. Defaults to `false`.
             * @param payload.rightTrim - Optional. Whether to trim whitespace from the right side of the variable. Defaults to `false`.
             *
             * @returns ReturnType - The result of the command execution.
             */
            insertTemplateBlock: (payload: { value: string; leftTrim?: boolean; rightTrim?: boolean }) => ReturnType;
        };
    }
}

export const TemplateBlock = Node.create({
    name: 'templateBlock',

    group: 'block',
    content: 'block+',
    defining: true,
    isolating: true,

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            'data-template-block': { default: '' }, // e.g. "range .Items"
            'data-left-trim': {
                default: false,
                parseHTML: (element) => element.getAttribute('data-left-trim') === 'true',
            },
            'data-right-trim': {
                default: false,
                parseHTML: (element) => element.getAttribute('data-right-trim') === 'true',
            },
        };
    },

    parseHTML() {
        return [
            {
                tag: 'div[data-template-block]',
            },
        ];
    },

    renderHTML(element) {
        const { HTMLAttributes } = element;
        const { 'data-template-block': value, 'data-left-trim': leftTrim, 'data-right-trim': rightTrim } = HTMLAttributes;
        const opening = leftTrim ? '{{-' : '{{';
        const closing = rightTrim ? '-}}' : '}}';

        return [
            'div',
            mergeAttributes(HTMLAttributes, {
                'data-template-block': value,
                class: 'template-block',
            }),
            ['div', { class: 'template-open' }, `${opening} ${value} ${closing}`],
            ['div', { class: 'template-inner' }, 0],
            ['div', { class: 'template-close' }, `{{ end }}`],
        ];
    },

    addCommands() {
        return {
            insertTemplateBlock:
                ({ value, leftTrim = false, rightTrim = false }) =>
                ({ commands }) =>
                    commands.insertContent({
                        type: this.name,
                        attrs: { 'data-template-block': value, 'data-left-trim': leftTrim, 'data-right-trim': rightTrim },
                        content: [{ type: 'paragraph' }],
                    }),
        };
    },
});
