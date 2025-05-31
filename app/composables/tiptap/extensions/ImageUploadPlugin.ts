import type { Editor } from '@tiptap/vue-3';
import { Plugin } from 'prosemirror-state';

export function imageUploadPlugin(_: Editor, onFiles: (files: File[]) => Promise<void>) {
    return new Plugin({
        props: {
            // @ts-expect-error Prosemirror types do not include these methods
            handlePaste: (_, e) => onFiles([...(e.clipboardData?.files ?? [])]),
            // @ts-expect-error Prosemirror types do not include these methods
            handleDrop: (_, e) => onFiles([...(e.dataTransfer?.files ?? [])]),
        },
    });
}
