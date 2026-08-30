import { Node, mergeAttributes } from '@tiptap/core';
import type { DOMOutputSpec, Node as ProseMirrorNode } from '@tiptap/pm/model';
import { Plugin, PluginKey, type EditorState, type Transaction } from '@tiptap/pm/state';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import TemplateBlockNodeView from '~/components/documents/templates/editor/TemplateBlockNodeView.vue';
import { isTemplateBlockActionValid } from './TemplateBlockValidation';

export type TemplateActionKind = 'block-start' | 'block-end' | 'comment' | 'raw-control' | 'variable';

const templateActionPattern = /\{\{(-)?\s*([^{}]*?)\s*(-)?\}\}/g;

export function createTemplateActionMatcher(): RegExp {
    return new RegExp(templateActionPattern.source, templateActionPattern.flags);
}

const blockStartActions = new Set(['range', 'if', 'else', 'with', 'break', 'continue']);
const rawControlActions = new Set(['define', 'template', 'block', 'end']);

export function classifyTemplateAction(value: string): TemplateActionKind {
    const action = value.trim();
    const actionName = action.split(/\s+/, 1)[0] ?? '';

    if (action === 'end') return 'block-end';
    if (action.startsWith('/*') && action.endsWith('*/')) return 'comment';
    if (blockStartActions.has(actionName)) return 'block-start';
    if (rawControlActions.has(actionName)) return 'raw-control';

    return 'variable';
}

function createBlockNormalizationTransaction(state: EditorState): Transaction | null {
    const templateAction = createTemplateActionMatcher();
    const replacements: Array<{
        from: number;
        to: number;
        type: 'start' | 'end';
        value: string;
        leftTrim: boolean;
        rightTrim: boolean;
        marks: ProseMirrorNode['marks'];
    }> = [];

    state.doc.descendants((node, pos) => {
        if (!node.isText || !node.text) return;

        templateAction.lastIndex = 0;
        let match: RegExpExecArray | null;
        while ((match = templateAction.exec(node.text)) !== null) {
            const value = match[2]?.trim() ?? '';
            const actionKind = classifyTemplateAction(value);
            if (actionKind !== 'block-start' && actionKind !== 'block-end') continue;

            replacements.push({
                from: pos + match.index,
                to: pos + match.index + match[0].length,
                type: actionKind === 'block-start' ? 'start' : 'end',
                value: actionKind === 'block-start' ? value : 'end',
                leftTrim: Boolean(match[1]),
                rightTrim: Boolean(match[3]),
                marks: node.marks,
            });
        }
    });

    if (replacements.length === 0) return null;

    const transaction = state.tr;
    replacements.reverse().forEach((replacement) => {
        const nodeType = state.schema.nodes[replacement.type === 'start' ? 'templateBlock' : 'templateBlockEnd'];
        if (!nodeType) return;

        transaction.replaceWith(
            replacement.from,
            replacement.to,
            nodeType.create(
                {
                    ...(replacement.type === 'start' ? { 'data-template-block': replacement.value } : {}),
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

export interface TemplateBlockOptions {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    HTMLAttributes: Record<string, any>;
}

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        templateBlock: {
            insertTemplateBlock: (payload: { value: string; leftTrim?: boolean; rightTrim?: boolean }) => ReturnType;
            insertTemplateBlockEnd: (payload?: { leftTrim?: boolean; rightTrim?: boolean }) => ReturnType;
        };
    }
}

function renderActionHTML(
    attributes: Record<string, unknown>,
    value: string,
    leftTrim: boolean,
    rightTrim: boolean,
    end = false,
): DOMOutputSpec {
    const opening = leftTrim ? '{{-' : '{{';
    const closing = rightTrim ? '-}}' : '}}';
    const actionAttributes = { ...attributes };
    delete actionAttributes['data-template-block'];
    delete actionAttributes['data-template-block-end'];

    return [
        'span',
        mergeAttributes(actionAttributes, {
            ...(end ? { 'data-template-block-end': 'end' } : { 'data-template-block': value }),
            'data-left-trim': leftTrim,
            'data-right-trim': rightTrim,
            class: 'template-block',
        }),
        `${opening} ${value} ${closing}`,
    ];
}

export const TemplateBlock = Node.create<TemplateBlockOptions>({
    name: 'templateBlock',
    inline: true,
    group: 'inline',
    atom: true,

    addOptions() {
        return { HTMLAttributes: {} };
    },

    onCreate({ editor }) {
        const transaction = createBlockNormalizationTransaction(editor.state);
        if (transaction) editor.view.dispatch(transaction);
    },

    addNodeView() {
        return VueNodeViewRenderer(TemplateBlockNodeView);
    },

    addAttributes() {
        return {
            'data-template-block': { default: '' },
            'data-left-trim': {
                default: false,
                parseHTML: (element: Element) => element.getAttribute('data-left-trim') === 'true',
            },
            'data-right-trim': {
                default: false,
                parseHTML: (element: Element) => element.getAttribute('data-right-trim') === 'true',
            },
        };
    },

    parseHTML() {
        return [
            {
                tag: 'span[data-template-block]',
                getAttrs: (element) => (element.hasAttribute('data-template-block-end') ? false : null),
            },
        ];
    },

    renderHTML({ HTMLAttributes }) {
        return renderActionHTML(
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes),
            HTMLAttributes['data-template-block'] as string,
            HTMLAttributes['data-left-trim'] as boolean,
            HTMLAttributes['data-right-trim'] as boolean,
        );
    },

    addCommands() {
        return {
            insertTemplateBlock:
                ({ value, leftTrim = false, rightTrim = false }) =>
                ({ commands }) => {
                    if (!isTemplateBlockActionValid(value)) return false;

                    return commands.insertContent({
                        type: this.name,
                        attrs: { 'data-template-block': value, 'data-left-trim': leftTrim, 'data-right-trim': rightTrim },
                    });
                },
            insertTemplateBlockEnd:
                ({ leftTrim = false, rightTrim = false } = {}) =>
                ({ commands }) =>
                    commands.insertContent({
                        type: 'templateBlockEnd',
                        attrs: { 'data-left-trim': leftTrim, 'data-right-trim': rightTrim },
                    }),
        };
    },

    addProseMirrorPlugins() {
        return [
            new Plugin({
                key: new PluginKey('templateBlockNormalize'),
                appendTransaction: (_transactions, _oldState, newState) => createBlockNormalizationTransaction(newState),
            }),
        ];
    },
});

export const TemplateBlockEnd = Node.create<TemplateBlockOptions>({
    name: 'templateBlockEnd',
    inline: true,
    group: 'inline',
    atom: true,

    addOptions() {
        return { HTMLAttributes: {} };
    },

    addNodeView() {
        return VueNodeViewRenderer(TemplateBlockNodeView);
    },

    addAttributes() {
        return {
            'data-left-trim': {
                default: false,
                parseHTML: (element: Element) => element.getAttribute('data-left-trim') === 'true',
            },
            'data-right-trim': {
                default: false,
                parseHTML: (element: Element) => element.getAttribute('data-right-trim') === 'true',
            },
        };
    },

    parseHTML() {
        return [{ tag: 'span[data-template-block-end]' }];
    },

    renderHTML({ HTMLAttributes }) {
        return renderActionHTML(
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes),
            'end',
            HTMLAttributes['data-left-trim'] as boolean,
            HTMLAttributes['data-right-trim'] as boolean,
            true,
        );
    },
});
