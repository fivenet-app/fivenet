import type { JsonObject } from '@protobuf-ts/runtime';
import { describe, expect, it } from 'vitest';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { ContentType, NodeType, type Content, type RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';
import { getTextFromContent, jsonNodeToTocLinks } from './content';

function richNode(node: Partial<RichTextHtmlNode> & Pick<RichTextHtmlNode, 'tag'>): RichTextHtmlNode {
    const { tag, ...rest } = node;

    return {
        type: NodeType.ELEMENT,
        tag,
        attrs: {},
        content: [],
        ...rest,
    };
}

function richContent(content: RichTextHtmlNode): Content {
    return {
        version: '1',
        contentType: ContentType.HTML,
        content,
    };
}

function tiptapContent(tiptapJson: JsonObject): Content {
    return {
        version: '1',
        contentType: ContentType.TIPTAP_JSON,
        tiptapJson: Struct.fromJson(tiptapJson),
    };
}

describe('jsonNodeToTocLinks', () => {
    it('should return no toc links without content', () => {
        const content: Content = {
            version: '1',
            contentType: ContentType.UNSPECIFIED,
        };

        expect(jsonNodeToTocLinks(content)).toEqual([]);
    });

    it('should extract toc links from rich html headings', () => {
        const content = richContent(
            richNode({
                type: NodeType.DOC,
                tag: 'body',
                content: [
                    richNode({
                        id: 'intro',
                        tag: 'h2',
                        content: [richNode({ type: NodeType.TEXT, tag: 'text', text: 'Introduction' })],
                    }),
                    richNode({
                        tag: 'p',
                        content: [richNode({ type: NodeType.TEXT, tag: 'text', text: 'Body text' })],
                    }),
                    richNode({
                        tag: 'section',
                        content: [
                            richNode({
                                tag: 'h3',
                                content: [
                                    richNode({
                                        tag: 'span',
                                        content: [richNode({ type: NodeType.TEXT, tag: 'text', text: 'Nested Heading' })],
                                    }),
                                ],
                            }),
                        ],
                    }),
                ],
            }),
        );

        expect(jsonNodeToTocLinks(content)).toEqual([
            {
                id: 'intro',
                depth: 2,
                text: 'Introduction',
            },
            {
                id: 'h3',
                depth: 3,
                text: 'Nested Heading',
            },
        ]);
    });

    it('should prefer rich html content over tiptap json content', () => {
        const content: Content = {
            ...richContent(
                richNode({
                    type: NodeType.DOC,
                    tag: 'body',
                    content: [richNode({ tag: 'h1', content: [richNode({ tag: 'text', text: 'Rich Text' })] })],
                }),
            ),
            tiptapJson: Struct.fromJson({
                type: 'doc',
                content: [{ type: 'heading', attrs: { level: 1 }, content: [{ type: 'text', text: 'Tiptap Text' }] }],
            }),
        };

        expect(jsonNodeToTocLinks(content)).toEqual([
            {
                id: 'h1',
                depth: 1,
                text: 'Rich Text',
            },
        ]);
    });

    it('should extract toc links from tiptap json headings', () => {
        const content = tiptapContent({
            type: 'doc',
            content: [
                {
                    type: 'heading',
                    attrs: {
                        id: 'setup',
                        level: 2,
                    },
                    content: [{ type: 'text', text: 'Project Setup' }],
                },
                {
                    type: 'paragraph',
                    content: [{ type: 'text', text: 'Body text' }],
                },
                {
                    type: 'heading',
                    attrs: {
                        level: 3,
                    },
                    content: [{ type: 'text', text: 'API' }, { type: 'hardBreak' }, { type: 'text', text: ' Reference' }],
                },
            ],
        });

        expect(jsonNodeToTocLinks(content)).toEqual([
            {
                id: 'setup',
                depth: 2,
                text: 'Project Setup',
            },
            {
                id: 'h3-api-reference',
                depth: 3,
                text: 'API Reference',
            },
        ]);
    });
});

describe('getTextFromContent', () => {
    it('should return direct text content', () => {
        expect(getTextFromContent(richNode({ tag: 'text', text: 'Direct text' }))).toBe('Direct text');
    });

    it('should concatenate nested text content', () => {
        const node = richNode({
            tag: 'h2',
            content: [
                richNode({ tag: 'span', content: [richNode({ tag: 'text', text: 'Hello ' })] }),
                richNode({ tag: 'strong', content: [richNode({ tag: 'text', text: 'World' })] }),
            ],
        });

        expect(getTextFromContent(node)).toBe('Hello World');
    });

    it('should fall back to the node id when no text exists', () => {
        expect(getTextFromContent(richNode({ id: 'empty-heading', tag: 'h2' }))).toBe('empty-heading');
    });

    it('should return an empty string without text or id', () => {
        expect(getTextFromContent(richNode({ tag: 'h2' }))).toBe('');
    });
});
