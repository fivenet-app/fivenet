import { Node, mergeAttributes } from '@tiptap/core';
import type { Node as ProseMirrorNode } from '@tiptap/pm/model';
import { Plugin, PluginKey, type EditorState, type Transaction } from '@tiptap/pm/state';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import TemplateVarNodeView from '~/components/documents/templates/editor/TemplateVarNodeView.vue';
import { classifyTemplateAction, createTemplateActionMatcher } from './TemplateBlock';

function createTemplateVarNormalizationTransaction(state: EditorState, nodeName: string): Transaction | null {
    const templateAction = createTemplateActionMatcher();
    const replacements: Array<{
        from: number;
        to: number;
        expression: string;
        leftTrim: boolean;
        rightTrim: boolean;
        marks: ProseMirrorNode['marks'];
    }> = [];

    state.doc.descendants((node, pos) => {
        if (!node.isText || !node.text) return;

        templateAction.lastIndex = 0;
        let match: RegExpExecArray | null;
        while ((match = templateAction.exec(node.text)) !== null) {
            const expression = match[2]?.trim() ?? '';
            if (!expression || classifyTemplateAction(expression) !== 'variable') continue;

            replacements.push({
                from: pos + match.index,
                to: pos + match.index + match[0].length,
                expression,
                leftTrim: Boolean(match[1]),
                rightTrim: Boolean(match[3]),
                marks: node.marks,
            });
        }
    });

    if (replacements.length === 0) return null;

    const templateVarNode = state.schema.nodes[nodeName];
    if (!templateVarNode) return null;

    const transaction = state.tr;
    replacements.reverse().forEach((replacement) => {
        transaction.replaceWith(
            replacement.from,
            replacement.to,
            templateVarNode.create(
                {
                    'data-template-var': replacement.expression,
                    'data-left-trim': replacement.leftTrim,
                    'data-right-trim': replacement.rightTrim,
                },
                null,
                replacement.marks,
            ),
        );
    });

    return transaction;
}

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

export const TemplateVar = Node.create<TemplateVarOptions>({
    name: 'templateVar',

    inline: true,

    group: 'inline',

    atom: true,

    addNodeView() {
        return VueNodeViewRenderer(TemplateVarNodeView);
    },

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    onCreate({ editor }) {
        const transaction = createTemplateVarNormalizationTransaction(editor.state, this.name);
        if (transaction) editor.view.dispatch(transaction);
    },

    addProseMirrorPlugins() {
        return [
            new Plugin({
                key: new PluginKey('templateVarNormalize'),
                appendTransaction: (_transactions, _oldState, newState) => {
                    return createTemplateVarNormalizationTransaction(newState, this.name);
                },
            }),
        ];
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
