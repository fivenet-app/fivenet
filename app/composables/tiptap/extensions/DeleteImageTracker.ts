import { Extension } from '@tiptap/core';
import type { Node } from 'prosemirror-model';
import { Plugin, PluginKey } from 'prosemirror-state';

export const deleteImageTrackerKey = new PluginKey('delete-image-tracker');
export const deleteImageTrackerForceRemovedMetaKey = 'delete-image-tracker:force-removed';

// Helper - collect current <image> fileIds in the document
function collectImageIds(doc: Node): Set<number> {
    const out = new Set<number>();
    doc.descendants((node) => {
        if (node.type.name === 'image' && node.attrs.fileId) {
            out.add(parseInt(node.attrs.fileId as string));
        }
    });
    return out;
}

// Factory that triggers `onRemoved(ids)` whenever an image vanishes
function deleteImageTracker(onRemoved: (ids: number[]) => void) {
    return new Plugin({
        key: deleteImageTrackerKey,

        state: {
            init: (_, state) => collectImageIds(state.doc),
            apply: (tr, prev: Set<number>, _old, newState) => {
                const forcedRemoved = tr.getMeta(deleteImageTrackerForceRemovedMetaKey) as number[] | undefined;
                if (forcedRemoved?.length) onRemoved([...new Set(forcedRemoved.filter((id) => !!id))]);

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

// Exported Tiptap extension
export const DeleteImageTracker = Extension.create<{
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
        return [deleteImageTracker(this.options.onRemoved)];
    },
});
