// extensions/EnhancedImage.ts
import { Node, mergeAttributes } from '@tiptap/core';

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        enhancedImage: {
            /**
             * Insert an image (= writer uploads or paste)
             */
            setEnhancedImage: (options: { src: string; alt?: string; title?: string; fileId?: number }) => ReturnType;
        };
    }
}

export const EnhancedImage = Node.create({
    name: 'image',
    inline: true,
    group: 'inline',
    draggable: true,
    selectable: true,
    atom: true,

    addOptions() {
        return {
            HTMLAttributes: {},
        };
    },

    addAttributes() {
        return {
            src: {
                default: null,
            },
            alt: {
                default: null,
            },
            title: {
                default: null,
            },
            fileId: {
                default: null,
                parseHTML: (element) => element.getAttribute('data-file-id') || null,
                renderHTML: (attributes) => {
                    if (!attributes.fileId) return {};
                    return { 'data-file-id': attributes.fileId };
                },
            },
        };
    },

    parseHTML() {
        return [
            {
                tag: 'img[src]',
            },
        ];
    },

    renderHTML({ HTMLAttributes }) {
        return ['img', mergeAttributes(this.options.HTMLAttributes, HTMLAttributes)];
    },

    addCommands() {
        return {
            setEnhancedImage:
                (attrs) =>
                ({ chain }) =>
                    chain()
                        .insertContent({
                            type: this.name,
                            attrs,
                        })
                        .run(),
        };
    },
});
