// extensions/EnhancedImage.ts
import { Node, mergeAttributes } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import EnhancedImageView from '~/components/partials/editor/EnhancedImageView.vue';

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
            style: {
                default: 'width: 100%; height: auto; cursor: pointer; margin: 0 auto 0 0;',
                parseHTML: (element) => {
                    // Get style string from element
                    let style = element.getAttribute('style') || '';
                    // If width attribute is present, ensure it's in the style
                    const width = element.getAttribute('width');
                    if (width && !/width:\s*\d+px;/.test(style)) {
                        style = `width: ${width}px; ${style}`;
                    }
                    // If margin is present as attribute or in style, preserve it
                    const margin = element.style.margin || element.getAttribute('margin');
                    if (margin && !/margin:/.test(style)) {
                        style += ` margin: ${margin};`;
                    }
                    // Always ensure cursor and height are present
                    if (!/height:/.test(style)) {
                        style += ' height: auto;';
                    }
                    if (!/cursor:/.test(style)) {
                        style += ' cursor: pointer;';
                    }
                    return style.trim();
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

    addNodeView() {
        return VueNodeViewRenderer(EnhancedImageView);
    },
});
