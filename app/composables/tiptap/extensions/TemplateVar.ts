import { Node, mergeAttributes } from '@tiptap/core';

export interface TemplateVarOptions {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    HTMLAttributes: Record<string, any>;
}

declare module '@tiptap/core' {
    /**
     * Represents a set of commands related to template variables.
     *
     * @template ReturnType - The type returned by the commands.
     */
    interface Commands<ReturnType> {
        /**
         * Commands for managing template variables.
         */
        templateVar: {
            /**
             * Inserts a template variable into the editor.
             *
             * @param payload - The payload containing the template variable details.
             * @param payload.value - The value of the template variable to insert.
             * @param payload.leftTrim - Optional. Whether to trim whitespace from the left side of the variable. Defaults to `false`.
             * @param payload.rightTrim - Optional. Whether to trim whitespace from the right side of the variable. Defaults to `false`.
             * @returns ReturnType - The result of the command execution.
             */
            insertTemplateVar: (payload: { value: string; leftTrim?: boolean; rightTrim?: boolean }) => ReturnType;
        };
    }
}

export const TemplateVar = Node.create({
    name: 'templateVar',

    inline: true,

    group: 'inline',

    atom: true,

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            'data-template-var': {
                default: '',
            },
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
                tag: 'span[data-template-var]',
            },
        ];
    },

    renderHTML({ HTMLAttributes }) {
        const {
            'data-template-var': dataTemplateVar,
            'data-left-trim': dataLeftTrim,
            'data-right-trim': dataRightTrim,
        } = HTMLAttributes;
        const opening = dataLeftTrim ? '{{-' : '{{';
        const closing = dataRightTrim ? '-}}' : '}}';
        return [
            'span',
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes, {
                'data-template-var': dataTemplateVar,
                'data-left-trim': dataLeftTrim,
                'data-right-trim': dataRightTrim,
                class: 'template-var',
            }),
            `${opening} ${dataTemplateVar} ${closing}`,
        ];
    },

    addCommands() {
        return {
            insertTemplateVar:
                ({ value, leftTrim = false, rightTrim = false }) =>
                ({ commands }) => {
                    return commands.insertContent({
                        type: this.name,
                        attrs: { 'data-template-var': value, 'data-left-trim': leftTrim, 'data-right-trim': rightTrim },
                    });
                },
        };
    },
});
