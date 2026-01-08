import { getSchema, type Extensions, type JSONContent } from '@tiptap/core';
import type { DOMOutputSpec, Mark as PMMark, Node as PMNode } from 'prosemirror-model';
import { Fragment, h, type VNodeChild } from 'vue';

/* eslint-disable @typescript-eslint/no-explicit-any */
export function isSameDoc(a: any, b: any, extensions: Extensions) {
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
