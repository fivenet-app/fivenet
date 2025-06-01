import { Extension } from '@tiptap/core';
import type { Node } from 'prosemirror-model';
import { Plugin, PluginKey } from 'prosemirror-state';

/* ------------------------------------------------------------------ */
/*  Helper – collect current <image> fileIds in the document          */
/* ------------------------------------------------------------------ */
function collectImageIds(doc: Node): Set<number> {
    const out = new Set<number>();
    doc.descendants((node) => {
        if (node.type.name === 'image' && node.attrs.fileId) {
            out.add(parseInt(node.attrs.fileId as string));
        }
    });
    return out;
}

/* ------------------------------------------------------------------ */
/*  Factory that triggers `onRemoved(ids)` whenever an image vanishes */
/* ------------------------------------------------------------------ */
function DeleteImageTracker(onRemoved: (ids: number[]) => void) {
    return new Plugin({
        key: new PluginKey('delete-image-tracker'),

        state: {
            init: (_, state) => collectImageIds(state.doc),
            apply: (tr, prev: Set<number>, _old, newState) => {
                if (!tr.docChanged) return prev;
                const next = collectImageIds(newState.doc);

                const removed: number[] = [];
                prev.forEach((id) => {
                    if (!next.has(id)) removed.push(id);
                });

                if (removed.length) onRemoved(removed);
                return next;
            },
        },
    });
}

/* ------------------------------------------------------------------ */
/*  Exported Tiptap extension                                         */
/* ------------------------------------------------------------------ */
export const DeleteImageTrackerExt = Extension.create<{
    /** callback fired with *all* fileIds that disappeared in a tx */
    onRemoved: (ids: number[]) => void;
}>({
    name: 'deleteImageObserver',

    addOptions() {
        return {
            onRemoved: () => {}, // no-op default
        };
    },

    addProseMirrorPlugins() {
        return [DeleteImageTracker(this.options.onRemoved)];
    },
});
