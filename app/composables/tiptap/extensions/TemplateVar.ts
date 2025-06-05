import { Mark, mergeAttributes } from '@tiptap/core';

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

export const TemplateVar = Mark.create<TemplateVarOptions>({
    name: 'templateVar',

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            value: {
                default: '',
            },
            leftTrim: {
                default: false,
                parseHTML: (element) => element.getAttribute('data-left-trim') === 'true',
            },
            rightTrim: {
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
        const { value, leftTrim, rightTrim } = HTMLAttributes;
        const opening = leftTrim ? '{{-' : '{{';
        const closing = rightTrim ? '-}}' : '}}';
        return [
            'span',
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes, {
                'data-template-var': value,
                'data-left-trim': leftTrim,
                'data-right-trim': rightTrim,
                class: 'template-var',
                contenteditable: 'false',
            }),
            `${opening} ${value} ${closing}`,
        ];
    },

    addCommands() {
        return {
            insertTemplateVar:
                ({ value, leftTrim = false, rightTrim = false }) =>
                ({ commands }) => {
                    const opening = leftTrim ? '{{-' : '{{';
                    const closing = rightTrim ? '-}}' : '}}';
                    return commands.insertContent({
                        type: 'text',
                        text: `${opening} ${value} ${closing}`,
                        marks: [
                            {
                                type: 'templateVar',
                                attrs: { value, leftTrim, rightTrim },
                            },
                        ],
                    });
                },
        };
    },
});
