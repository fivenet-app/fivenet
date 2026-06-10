/* eslint-disable @typescript-eslint/no-explicit-any */
import { Node, mergeAttributes } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import EnhancedImageView from '~/components/partials/editor/EnhancedImageView.vue';
import { deleteImageTrackerForceRemovedMetaKey } from './DeleteImageTracker';

export interface EnhancedImageOptions {
    /**
     * HTML attributes to add to the task item element.
     * @default {}
     * @example { class: 'foo' }
     */
    HTMLAttributes: Record<string, any>;
}

declare module '@tiptap/core' {
    interface Commands<ReturnType> {
        enhancedImage: {
            /**
             * Insert an image (= writer uploads or paste)
             */
            setEnhancedImage: (options: { src: string; alt?: string; title?: string; fileId?: number }) => ReturnType;
            /**
             * Remove all image nodes for a tracked file.
             */
            removeEnhancedImageByFileId: (fileId: number) => ReturnType;
        };
    }
}

export type ImageAlign = 'left' | 'center' | 'right';

const DEFAULT_IMAGE_ALIGN: ImageAlign = 'left';

function extractStyleValue(style: unknown, property: 'width' | 'height'): string | null {
    if (typeof style !== 'string') return null;

    const match = style.match(new RegExp(`(?:^|;)\\s*${property}\\s*:\\s*([^;]+)`, 'i'));
    return match?.[1]?.trim() ?? null;
}

function parsePixelValue(value: string | null): number | null {
    if (!value) return null;

    const match = value.match(/^(\d+)px$/i);
    return match ? Number(match[1]) : null;
}

function getAlignFromStyle(style: unknown): ImageAlign | null {
    if (typeof style !== 'string') return null;

    const normalized = style.replace(/\s+/g, ' ').toLowerCase();

    if (/margin:\s*0\s+auto/.test(normalized)) {
        return 'center';
    }

    if (/margin:\s*0\s+0\s+0\s+auto/.test(normalized)) {
        return 'right';
    }

    if (/margin:\s*0\s+auto\s+0\s+0/.test(normalized)) {
        return 'left';
    }

    return null;
}

export function removeStyleProperties(style: unknown, properties: string[]): string {
    if (typeof style !== 'string') return '';

    const blocked = new Set(properties.map((p) => p.toLowerCase()));

    return style
        .split(';')
        .map((part) => part.trim())
        .filter(Boolean)
        .filter((part) => {
            const [key] = part.split(':', 1);
            return !blocked.has(key!.trim().toLowerCase());
        })
        .join('; ');
}

function normalizeImageAttrs(attrs: Record<string, any>): Record<string, any> {
    const next = { ...attrs };

    if (!next.align) {
        next.align = getAlignFromStyle(next.style) ?? attrs['data-align'] ?? DEFAULT_IMAGE_ALIGN;
    }

    if (next.width == null) {
        next.width = parsePixelValue(extractStyleValue(next.style, 'width'));
    }
    if (next.height == null) {
        next.height = parsePixelValue(extractStyleValue(next.style, 'height'));
    }

    next.style = removeStyleProperties(next.style, ['width', 'height', 'margin']);

    return next;
}

export function getAlignStyle(align: unknown): string {
    switch (align) {
        case 'center':
            return 'margin: 0 auto;';
        case 'right':
            return 'margin: 0 0 0 auto;';
        case 'left':
        default:
            return 'margin: 0 auto 0 0;';
    }
}

const mergeStyle = (...styles: Array<string | null | undefined>) => {
    return styles
        .map((style) => style?.trim())
        .filter(Boolean)
        .join(' ');
};

export const EnhancedImage = Node.create<EnhancedImageOptions>({
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
            align: {
                default: DEFAULT_IMAGE_ALIGN,
                parseHTML: (element) => {
                    console.log('parsed from html', element);
                    return (
                        element.getAttribute('data-align') ||
                        getAlignFromStyle(element.getAttribute('style')) ||
                        DEFAULT_IMAGE_ALIGN
                    );
                },
                renderHTML: (attributes) => {
                    return { 'data-align': attributes.align ?? DEFAULT_IMAGE_ALIGN };
                },
            },
            width: {
                default: null,
                parseHTML: (element) => {
                    const width = element.getAttribute('width');
                    if (width && Number(width) > 0) return Number(width);

                    const match = element.getAttribute('style')?.match(/width:\s*(\d+)px/i);
                    return match ? Number(match[1]) : null;
                },
                renderHTML: (attributes) => {
                    if (attributes.width == null) return {};
                    return { width: attributes.width };
                },
            },
            height: {
                default: null,
                parseHTML: (element) => {
                    const height = element.getAttribute('height');
                    if (height && Number(height) > 0) return Number(height);

                    const match = element.getAttribute('style')?.match(/height:\s*(\d+)px/i);
                    return match ? Number(match[1]) : null;
                },
                renderHTML: (attributes) => {
                    if (attributes.height == null) return {};
                    return { height: attributes.height };
                },
            },
            style: {
                default: 'width: 100%; height: auto; cursor: pointer;',
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
        const attrs = normalizeImageAttrs(HTMLAttributes);

        const { align, style, ...rest } = attrs;

        return [
            'img',
            mergeAttributes(this.options.HTMLAttributes, rest, {
                src: cleanupImageURL(attrs.src),
                'data-align': align,
                style: mergeStyle(style, 'display: block; ' + getAlignStyle(align)),
            }),
        ];
    },

    parseMarkdown: (token, helpers) => {
        return helpers.createNode('image', {
            src: token.href,
            title: token.title,
            alt: token.text,
        });
    },

    renderMarkdown: (node) => {
        const src = cleanupImageURL(node.attrs?.src ?? '');
        const alt = node.attrs?.alt ?? '';
        const title = node.attrs?.title ?? '';

        return title ? `![${alt}](${src} "${title}")` : `![${alt}](${src})`;
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

            removeEnhancedImageByFileId:
                (fileId) =>
                ({ state, dispatch }) => {
                    if (!fileId) return false;

                    const positions: number[] = [];
                    state.doc.descendants((node, pos) => {
                        if (node.type.name === this.name && Number(node.attrs.fileId) === fileId) {
                            positions.push(pos);
                        }
                    });

                    if (!positions.length) {
                        if (dispatch) {
                            dispatch(state.tr.setMeta(deleteImageTrackerForceRemovedMetaKey, [fileId]));
                        }
                        return true;
                    }

                    if (dispatch) {
                        let tr = state.tr;

                        // Delete from back to front to keep positions valid.
                        for (let i = positions.length - 1; i >= 0; i--) {
                            const pos = positions[i]!;
                            const node = tr.doc.nodeAt(pos);
                            if (!node || node.type.name !== this.name) continue;

                            tr = tr.delete(pos, pos + node.nodeSize);
                        }

                        dispatch(tr);
                    }

                    return true;
                },
        };
    },

    addNodeView() {
        return VueNodeViewRenderer(EnhancedImageView);
    },
});
