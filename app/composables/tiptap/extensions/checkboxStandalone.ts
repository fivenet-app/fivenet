/* eslint-disable @typescript-eslint/no-explicit-any */
// Based on the TaskItem extension from tiptap editor
//
// https://github.com/ueberdosis/tiptap/blob/main/packages/extension-task-item/src/task-item.ts

import { mergeAttributes, Node } from '@tiptap/core';
import type { Node as ProseMirrorNode } from '@tiptap/pm/model';

export interface CheckboxStandaloneOptions {
    /**
     * A callback function that is called when the checkbox is clicked while the editor is in readonly mode.
     * @param node The prosemirror node of the task item
     * @param checked The new checked state
     * @returns boolean
     */
    onReadOnlyChecked?: (node: ProseMirrorNode, checked: boolean) => boolean;

    /**
     * HTML attributes to add to the task item element.
     * @default {}
     * @example { class: 'foo' }
     */
    HTMLAttributes: Record<string, any>;
}

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        checkboxStandalone: {
            /**
             * Toggle a checkbox standalone
             * @example editor.commands.addCheckboxStandalone()
             */
            addCheckboxStandalone: () => ReturnType;
        };
    }
}

/**
 * This extension allows you to create task items.
 * @see https://www.tiptap.dev/api/nodes/task-item
 */
export const CheckboxStandalone = Node.create<CheckboxStandaloneOptions>({
    name: 'checkboxStandalone',

    atom: true,
    inline: true,
    group: 'inline',

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            checked: {
                default: false,
                keepOnSplit: false,
                parseHTML: (element) => {
                    const dataChecked = element.getAttribute('data-checked');

                    return dataChecked === '' || dataChecked === 'true';
                },
                renderHTML: (attributes) => ({
                    'data-checked': attributes.checked,
                }),
            },
        };
    },

    parseHTML() {
        return [
            {
                tag: `span[data-type="${this.name}"]`,
                priority: 51,
            },
        ];
    },

    renderHTML({ node, HTMLAttributes }) {
        return [
            'span',
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes, {
                'data-type': this.name,
            }),
            [
                'label',
                [
                    'input',
                    {
                        type: 'checkbox',
                        checked: node.attrs.checked ? 'checked' : null,
                    },
                ],
                ' ',
            ],
        ];
    },

    addNodeView() {
        return ({ node, HTMLAttributes, getPos, editor }) => {
            const listItem = document.createElement('span');
            const checkboxWrapper = document.createElement('label');
            const checkboxStyler = document.createElement('span');
            const checkbox = document.createElement('input');

            checkboxWrapper.contentEditable = 'false';
            checkbox.type = 'checkbox';
            checkbox.addEventListener('mousedown', (event) => event.preventDefault());
            checkbox.addEventListener('change', (event) => {
                // if the editor isn’t editable and we don't have a handler for
                // readonly checks we have to undo the latest change
                if (!editor.isEditable && !this.options.onReadOnlyChecked) {
                    checkbox.checked = !checkbox.checked;

                    return;
                }

                const { checked } = event.target as any;

                if (editor.isEditable && typeof getPos === 'function') {
                    editor
                        .chain()
                        .focus(undefined, { scrollIntoView: false })
                        .command(({ tr }) => {
                            const position = getPos();

                            if (typeof position !== 'number') {
                                return false;
                            }
                            const currentNode = tr.doc.nodeAt(position);

                            tr.setNodeMarkup(position, undefined, {
                                ...currentNode?.attrs,
                                checked,
                            });

                            return true;
                        })
                        .run();
                }
                if (!editor.isEditable && this.options.onReadOnlyChecked) {
                    // Reset state if onReadOnlyChecked returns false
                    if (!this.options.onReadOnlyChecked(node, checked)) {
                        checkbox.checked = !checkbox.checked;
                    }
                }
            });

            checkboxStyler.innerHTML = '&nbsp;';

            Object.entries(this.options.HTMLAttributes).forEach(([key, value]) => {
                listItem.setAttribute(key, value);
            });

            listItem.dataset.checked = node.attrs.checked;
            checkbox.checked = node.attrs.checked;

            checkboxWrapper.append(checkbox, checkboxStyler);
            listItem.append(checkboxWrapper);

            Object.entries(HTMLAttributes).forEach(([key, value]) => {
                listItem.setAttribute(key, value);
            });

            return {
                dom: listItem,
                update: (updatedNode) => {
                    if (updatedNode.type !== this.type) {
                        return false;
                    }

                    listItem.dataset.checked = updatedNode.attrs.checked;
                    checkbox.checked = updatedNode.attrs.checked;

                    return true;
                },
            };
        };
    },

    addCommands() {
        return {
            addCheckboxStandalone:
                () =>
                ({ commands }) => {
                    return commands.insertContent({
                        type: this.name,
                    });
                },
        };
    },
});
