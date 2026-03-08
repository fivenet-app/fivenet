import { mergeAttributes, Node } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import PenaltyCaculatorEditorView from '~/components/quickbuttons/penaltycalculator/PenaltyCaculatorEditorView.vue';

declare module '@tiptap/core' {
    /**
     * Represents a set of commands related to template variables.
     *
     * @template ReturnType - The type returned by the commands.
     */
    interface Commands<ReturnType> {
        penaltyCalculator: {
            insertPenaltyCalculator: () => ReturnType;
        };
    }
}

export const PenaltyCalculator = Node.create({
    name: 'penaltyCalculator',
    inline: false,
    group: 'block',
    atom: true,

    addOptions() {
        return {};
    },

    addAttributes() {
        return {};
    },

    parseHTML() {
        return [
            {
                tag: 'div[data-type="penalty-calculator"]',
            },
            {
                tag: 'div[data-type="penaltyCalculator"]',
            },
            {
                tag: 'span[data-type="penalty-calculator"]',
            },
        ];
    },

    renderHTML({ HTMLAttributes }) {
        return [
            'div',
            mergeAttributes(this.options.HTMLAttributes, HTMLAttributes, {
                'data-type': 'penalty-calculator',
                'data-embed': 'penalty-calculator',
            }),
        ];
    },

    addCommands() {
        return {
            insertPenaltyCalculator:
                () =>
                ({ commands, state }) => {
                    let existingPos: number | null = null;
                    state.doc.descendants((node, pos) => {
                        if (node.type.name === this.name) {
                            existingPos = pos;
                            return false;
                        }
                        return true;
                    });

                    if (existingPos !== null) {
                        commands.focus();
                        return commands.setNodeSelection(existingPos);
                    }

                    return commands.insertContent({
                        type: this.name,
                    });
                },
        };
    },

    addNodeView() {
        return VueNodeViewRenderer(PenaltyCaculatorEditorView);
    },
});
