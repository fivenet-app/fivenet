import { Node } from '@tiptap/core';
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

    addCommands() {
        return {
            insertPenaltyCalculator:
                () =>
                ({ commands }) => {
                    return commands.insertContent({
                        type: this.name,
                        attrs: { value: 'Penalty Calculator' },
                    });
                },
        };
    },

    addNodeView() {
        return VueNodeViewRenderer(PenaltyCaculatorEditorView);
    },
});
