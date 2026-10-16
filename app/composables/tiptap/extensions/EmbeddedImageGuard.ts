import { Extension } from '@tiptap/core';
import type { Node as ProseMirrorNode } from '@tiptap/pm/model';
import { Plugin, PluginKey } from '@tiptap/pm/state';
import { isEmbeddedImageSource } from '../imageSources';

const embeddedImageGuardPluginKey = new PluginKey('embeddedImageGuard');

export function hasEmbeddedImage(doc: ProseMirrorNode): boolean {
    let found = false;
    doc.descendants((node) => {
        if (node.type.name === 'image' && isEmbeddedImageSource(node.attrs.src)) {
            found = true;
            return false;
        }
        return true;
    });
    return found;
}

export const EmbeddedImageGuard = Extension.create<{ onBlocked?: () => void }>({
    name: 'embeddedImageGuard',

    addOptions() {
        return {
            onBlocked: undefined,
        };
    },

    addProseMirrorPlugins() {
        const onBlocked = this.options.onBlocked;

        return [
            new Plugin({
                key: embeddedImageGuardPluginKey,
                props: {
                    handlePaste: (_view, event) => {
                        const html = event.clipboardData?.getData('text/html');
                        if (!html || !isEmbeddedImageSource(html.match(/<img[^>]+src=["']([^"']+)/i)?.[1])) return false;

                        onBlocked?.();
                        return true;
                    },
                },
                filterTransaction: (transaction, state) => {
                    if (!transaction.docChanged || !hasEmbeddedImage(transaction.doc) || hasEmbeddedImage(state.doc)) {
                        return true;
                    }

                    onBlocked?.();
                    return false;
                },
            }),
        ];
    },
});
