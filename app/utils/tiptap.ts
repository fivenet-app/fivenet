import { getSchema, type Extensions, type JSONContent } from '@tiptap/core';
import type { DOMOutputSpec, Mark as PMMark, Node as PMNode } from 'prosemirror-model';
import { Fragment, h, type VNodeChild } from 'vue';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { NodeType, type RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';

/* eslint-disable @typescript-eslint/no-explicit-any */
export function isSameDoc(a: JSONContent, b: JSONContent, extensions: Extensions) {
    const schema = getSchema(extensions);

    const nodeA = schema.nodeFromJSON(a);
    const nodeB = schema.nodeFromJSON(b);

    return nodeA.eq(nodeB);
}

export type VueStaticRendererOptions = {
    /**
     * Optional overrides (like nodeMapping/markMapping in static-renderer).
     * If provided for a name, it wins over schema.toDOM.
     */
    nodeMapping?: Record<string, (ctx: { node: PMNode; children: VNodeChild[] }) => VNodeChild>;
    markMapping?: Record<string, (ctx: { mark: PMMark; children: VNodeChild }) => VNodeChild>;

    /**
     * Fallbacks if schema has no toDOM.
     */
    unhandledNode?: (ctx: { name: string; node: PMNode; children: VNodeChild[] }) => VNodeChild;
    unhandledMark?: (ctx: { name: string; mark: PMMark; children: VNodeChild }) => VNodeChild;
};

export function renderToVueVNode(args: {
    extensions: Extensions;
    content: JSONContent | PMNode;
    options?: VueStaticRendererOptions;
}): VNodeChild {
    const { extensions, content, options } = args;
    const schema = getSchema(extensions);

    const doc: PMNode =
        typeof (content as any)?.type?.name === 'string' ? (content as PMNode) : schema.nodeFromJSON(content as JSONContent);

    const nodeMapping = options?.nodeMapping ?? {};
    const markMapping = options?.markMapping ?? {};

    const renderNode = (node: PMNode): VNodeChild => {
        // Text nodes: render text and wrap with marks using mark.toDOM
        if (node.isText) {
            let out: VNodeChild = node.text ?? '';

            for (const mark of node.marks) {
                const name = mark.type.name;

                const mapped = markMapping[name];
                if (mapped) {
                    out = mapped({ mark, children: out });
                    continue;
                }

                const toDOM = mark.type.spec.toDOM;
                if (toDOM) {
                    const spec = toDOM(mark, true);
                    out = domSpecToVNode(spec, out);
                    continue;
                }

                out = options?.unhandledMark ? options.unhandledMark({ name, mark, children: out }) : out;
            }

            return out;
        }

        // Non-text: render children, then build element from node.toDOM
        const children: VNodeChild[] = [];
        node.forEach((child) => children.push(renderNode(child)));

        const name = node.type.name;

        const mapped = nodeMapping[name];
        if (mapped) return mapped({ node, children });

        const toDOM = node.type.spec.toDOM;
        if (toDOM) {
            const spec = toDOM(node) as DOMOutputSpec;
            return domSpecToVNode(spec, children);
        }

        return options?.unhandledNode ? options.unhandledNode({ name, node, children }) : h(Fragment, null, children);
    };

    return renderNode(doc);
}

/**
 * Convert ProseMirror DOMOutputSpec into Vue VNodes.
 * Handles:
 *  - ["tag", {attrs?}, ...children]
 *  - 0 placeholder = insert provided children
 *  - nested specs
 *  - string children
 */
function domSpecToVNode(spec: DOMOutputSpec, children: VNodeChild | VNodeChild[]): VNodeChild {
    // If spec is a string, it’s a tag name with no attrs/children.
    if (typeof spec === 'string') {
        return h(spec, null, Array.isArray(children) ? children : [children]);
    }

    // ProseMirror also allows actual DOM Nodes in DOMOutputSpec,
    // but on the server we can’t create them. Treat as empty.
    if (!Array.isArray(spec)) {
        return Array.isArray(children) ? h(Fragment, null, children) : children;
    }

    const [tag, maybeAttrs, ...rest] = spec;

    const hasAttrs =
        maybeAttrs && typeof maybeAttrs === 'object' && !Array.isArray(maybeAttrs) && !(maybeAttrs as any).nodeType;

    const attrs = (hasAttrs ? maybeAttrs : null) as Record<string, any> | null;
    const content = (hasAttrs ? rest : [maybeAttrs, ...rest]) as any[];

    const renderedChildren: VNodeChild[] = [];

    for (const part of content) {
        if (part === 0) {
            if (Array.isArray(children)) renderedChildren.push(...children);
            else renderedChildren.push(children);
            continue;
        }

        if (part == null) continue;

        if (typeof part === 'string') {
            renderedChildren.push(part);
            continue;
        }

        // nested DOMOutputSpec
        renderedChildren.push(domSpecToVNode(part as DOMOutputSpec, children));
    }

    return h(tag as any, attrs, renderedChildren);
}

type ExtractTextOptions = {
    /**
     * When true, inserts newlines between block-ish nodes.
     * Good for previews and diffs.
     */
    blockSeparators?: boolean;

    /**
     * Replace repeated whitespace/newlines with single spaces in the final output.
     */
    collapseWhitespace?: boolean;

    /**
     * Suffix appended when output is truncated.
     */
    ellipsis?: string;

    /**
     * Limit elements processed.
     */
    limit?: number;
};

const ExtractTextOptionsDEFAULTS: Required<ExtractTextOptions> = {
    blockSeparators: true,
    collapseWhitespace: true,
    ellipsis: '…',
    limit: 150,
};

export function tiptapTextPreview(
    doc: JSONContent | null | undefined,
    maxChars: number,
    options: ExtractTextOptions = {},
): string {
    const opt = { ...ExtractTextOptionsDEFAULTS, ...options };
    if (!doc || maxChars <= 0) return '';

    let out = '';
    let remaining = maxChars;
    let truncated = false;

    const push = (s: string) => {
        if (!s || remaining <= 0) return;

        // Avoid expensive rune splitting. JS strings are UTF-16, so "chars" here are code units.
        // For UI previews this is usually fine.
        if (s.length <= remaining) {
            out += s;
            remaining -= s.length;
        } else {
            out += s.slice(0, remaining);
            remaining = 0;
            truncated = true;
        }
    };

    const blockBefore = (type?: string) => {
        if (!opt.blockSeparators) return;
        if (!type) return;

        // Separate blocks so words don't fuse
        switch (type) {
            case 'paragraph':
            case 'heading':
            case 'image':
            case 'blockquote':
            case 'codeBlock':
            case 'bulletList':
            case 'orderedList':
            case 'listItem':
            case 'taskList':
            case 'taskItem':
            case 'table':
            case 'tableRow':
                // If there's already content and we don't end with whitespace, add a newline
                if (out && !/\s$/.test(out)) push('\n');
                return;
        }
    };

    const blockAfter = (type?: string) => {
        if (!opt.blockSeparators) return;
        if (!type) return;

        switch (type) {
            case 'paragraph':
            case 'heading':
            case 'image':
            case 'blockquote':
            case 'codeBlock':
            case 'listItem':
            case 'taskItem':
            case 'tableRow':
                if (out && !/\s$/.test(out)) push('\n');
                return;
        }
    };

    const walk = (node: any) => {
        if (!node || remaining <= 0 || --opt.limit < 1) return;

        if (Array.isArray(node)) {
            for (const n of node) {
                if (remaining <= 0) break;
                walk(n);
            }
            return;
        }

        if (typeof node !== 'object') return;

        const type: string | undefined = node.type;

        // Ignore tables for text preview
        if (type === 'table' || type === 'tableRow' || type === 'tableCell' || type === 'tableHeader') {
            return;
        }

        // Text
        if (type === 'text') {
            push(String(node.text ?? ''));
            return;
        }

        if (type === 'hardBreak') {
            push(opt.blockSeparators ? '\n' : ' ');
            return;
        }

        if (type === 'mention') {
            const label = node.attrs?.label ?? node.attrs?.name ?? node.attrs?.id;
            push(label ? String(label) : '@mention');
            return;
        }

        if (type === 'image') {
            const alt = node.attrs?.alt;
            push(alt ? `[Image: ${String(alt)}]` : '[Image]');
            return;
        }

        if (type === 'taskItem') {
            const checked = !!node.attrs?.checked;
            push(checked ? '[x] ' : '[ ] ');
        }

        blockBefore(type);

        if (node.content) {
            walk(node.content);
        }

        blockAfter(type);
    };

    walk(doc);

    let result = out;

    if (opt.collapseWhitespace) {
        // collapse internal whitespace but keep it readable
        result = result.replace(/\s+/g, ' ').trim();
    } else {
        result = result.trimEnd();
    }

    if (truncated && opt.ellipsis) {
        // Ensure we don't exceed maxChars by much (ellipsis might add 1 char)
        if (result.length >= maxChars) {
            result = result.slice(0, Math.max(0, maxChars - opt.ellipsis.length)).trimEnd();
        }
        result += opt.ellipsis;
    }

    return result;
}

export function isEmptyDoc(content: Struct | JSONContent | null | undefined): boolean {
    if (!content) return true;

    let c: JSONContent;
    if ('fields' in content) {
        c = Struct.toJson(content as Struct) as JSONContent;
    } else {
        c = content as JSONContent;
    }

    return !c.content || c.content.length === 0 || tiptapTextPreview(c, 1, { limit: 25 }).length === 0;
}

export function isEmptyRichContentDoc(content: RichTextHtmlNode | null | undefined): boolean {
    if (!content || !content.content || content.content.length === 0) return true;

    // Check if all any top-level nodes are not empty
    let nodes: RichTextHtmlNode[] = [];
    for (let i = 0; i < 3; i++) {
        nodes = content.content[i]?.content ?? [];

        for (const node of nodes) {
            if (node.type === NodeType.ELEMENT) {
                if (node.tag === 'p') {
                    if ((node.content && node.content.length > 0) || (node.text && node.text.trim().length > 0)) {
                        // Paragraph has content
                        return false;
                    }
                } else if (node.tag === 'img') {
                    return false;
                }
            } else {
                // Found a non-paragraph node
                return false;
            }
        }
    }

    return true;
}
