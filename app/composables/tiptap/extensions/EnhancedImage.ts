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

export interface EnhancedImageAttrs {
    src?: string | null;
    alt?: string | null;
    title?: string | null;
    fileId?: number | string | null;
    align?: ImageAlign | null;
    width?: number | string | null;
    height?: number | string | null;
    style?: string | null;
    'data-align'?: string | null;
}

export type EnhancedImageAttrsInput = Partial<EnhancedImageAttrs> | Record<string, unknown> | undefined;

export interface NormalizedEnhancedImageAttrs extends Omit<EnhancedImageAttrs, 'align' | 'width' | 'height' | 'style'> {
    align: ImageAlign;
    width: number | null;
    height: number | null;
    style: string;
}

export const DEFAULT_IMAGE_ALIGN: ImageAlign = 'left';

export function extractStyleValue(style: unknown, property: 'width' | 'height'): string | null {
    if (typeof style !== 'string') return null;

    const match = style.match(new RegExp(`(?:^|;)\\s*${property}\\s*:\\s*([^;]+)`, 'i'));
    return match?.[1]?.trim() ?? null;
}

export function parseImageDimension(value: unknown): number | null {
    if (typeof value === 'number') {
        return Number.isFinite(value) && value > 0 ? value : null;
    }

    if (typeof value !== 'string') return null;

    const normalized = value.trim();
    if (!normalized) return null;

    const pixelMatch = normalized.match(/^(\d+)px$/i);
    if (pixelMatch) return Number(pixelMatch[1]);

    const parsed = Number.parseFloat(normalized);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : null;
}

export function getAlignFromStyle(style: unknown): ImageAlign | null {
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

export function normalizeImageAlign(value: unknown, style?: unknown): ImageAlign {
    if (value === 'center' || value === 'right' || value === 'left') {
        return value;
    }

    return getAlignFromStyle(style) ?? DEFAULT_IMAGE_ALIGN;
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

export function normalizeImageAttrs(
    attrs: EnhancedImageAttrsInput = {},
): Record<string, unknown> & NormalizedEnhancedImageAttrs {
    const next = { ...(attrs ?? {}) } as Record<string, unknown> & EnhancedImageAttrs;

    next.align = normalizeImageAlign(next.align ?? next['data-align'], next.style);

    next.width = parseImageDimension(next.width) ?? parseImageDimension(extractStyleValue(next.style, 'width'));
    next.height = parseImageDimension(next.height) ?? parseImageDimension(extractStyleValue(next.style, 'height'));

    next.style = removeStyleProperties(next.style, ['width', 'height', 'margin']);

    return next as Record<string, unknown> & NormalizedEnhancedImageAttrs;
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
                src: cleanupImageURL(String(attrs.src ?? '')),
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
