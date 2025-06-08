import type { Editor } from '@tiptap/vue-3';
import { Plugin } from 'prosemirror-state';

export function imageUploadPlugin(_: Editor, onFiles: (files: File[]) => Promise<void>) {
    return new Plugin({
        props: {
            handlePaste: (_, e) => {
                const files = [...(e.clipboardData?.files ?? [])];
                if (files.length === 0) return false; // Fallback to default paste
                onFiles(files);
                return true; // handled
            },
            handleDrop: (_, e) => {
                const files = [...(e.dataTransfer?.files ?? [])];
                if (files.length === 0) return false; // Fallback to default drop
                onFiles(files);
                return true; // handled
            },
        },
    });
}
