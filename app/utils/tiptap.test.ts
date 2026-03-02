import type { Extensions } from '@tiptap/core';
import { Blockquote } from '@tiptap/extension-blockquote';
import Document from '@tiptap/extension-document';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { describe, expect, it } from 'vitest';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { NodeType, type RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';
import { isEmptyDoc, isEmptyRichContentDoc, isSameDoc, tiptapTextPreview } from './tiptap';

describe('isSameDoc', () => {
    it('should return true for identical documents', () => {
        const extensions: Extensions = [Document, Text, Blockquote, Paragraph];
        const docA = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Hello' }] }] };
        const docB = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Hello' }] }] };

        expect(isSameDoc(docA, docB, extensions)).toBe(true);
    });

    it('should return false for different documents', () => {
        const extensions: Extensions = [Document, Text, Blockquote, Paragraph];
        const docA = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Hello' }] }] };
        const docB = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'World' }] }] };

        expect(isSameDoc(docA, docB, extensions)).toBe(false);
    });
});

describe('tiptapTextPreview', () => {
    it('should generate a text preview from a document', () => {
        const doc = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Hello World' }] }] };
        const preview = tiptapTextPreview(doc, 7);

        expect(preview).toBe('Hello…');
    });

    it('should generate a text preview from a document without ellipsis', () => {
        const doc = {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        { type: 'text', text: 'Hello ' },
                        { type: 'text', text: 'World' },
                    ],
                },
            ],
        };
        const preview = tiptapTextPreview(doc, 11);

        expect(preview).toBe('Hello World');
    });

    it('should handle empty documents', () => {
        const doc = { type: 'doc', content: [] };
        const preview = tiptapTextPreview(doc, 5);

        expect(preview).toBe('');
    });
});

describe('isEmptyDoc', () => {
    it('should return true for empty documents', () => {
        const doc = { type: 'doc', content: [] };

        expect(isEmptyDoc(doc)).toBe(true);
    });

    it('should return false for non-empty documents', () => {
        const doc = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Hello' }] }] };

        expect(isEmptyDoc(doc)).toBe(false);
    });

    it('should handle Struct input', () => {
        const struct = Struct.fromJson({ fields: { content: { listValue: { values: [] } } } });

        expect(isEmptyDoc(struct)).toBe(true);
    });
});

describe('isEmptyRichContentDoc', () => {
    it('should return true for empty rich content documents', () => {
        const doc: RichTextHtmlNode = {
            type: NodeType.DOC,
            tag: 'body',
            content: [{ type: NodeType.ELEMENT, tag: 'p', content: [], attrs: {} }],
            attrs: {},
        };

        expect(isEmptyRichContentDoc(doc)).toBe(true);
    });

    it('should return false for non-empty rich content documents', () => {
        const doc: RichTextHtmlNode = {
            type: NodeType.DOC,
            tag: 'body',
            content: [
                {
                    type: NodeType.ELEMENT,
                    tag: 'p',
                    content: [
                        {
                            type: NodeType.ELEMENT,
                            tag: 'p',
                            content: [{ type: NodeType.TEXT, tag: 'text', text: 'Hello', content: [], attrs: {} }],
                            attrs: {},
                        },
                    ],
                    attrs: {},
                },
            ],
            attrs: {},
        };

        expect(isEmptyRichContentDoc(doc)).toBe(false);
    });

    it('should return false for non-paragraph nodes', () => {
        const doc: RichTextHtmlNode = {
            type: NodeType.DOC,
            tag: 'body',
            content: [
                {
                    type: NodeType.ELEMENT,
                    tag: 'div',
                    content: [
                        {
                            type: NodeType.ELEMENT,
                            tag: 'div',
                            content: [{ type: NodeType.TEXT, tag: 'text', text: 'Hello', content: [], attrs: {} }],
                            attrs: {},
                        },
                    ],
                    attrs: {},
                },
            ],
            attrs: {},
        };

        expect(isEmptyRichContentDoc(doc)).toBe(false);
    });

    it('should return false for nested heading text content', () => {
        const doc: RichTextHtmlNode = {
            type: NodeType.DOC,
            tag: 'body',
            content: [
                { type: NodeType.ELEMENT, tag: 'p', content: [], attrs: {} },
                {
                    type: NodeType.ELEMENT,
                    tag: 'h3',
                    content: [
                        {
                            type: NodeType.ELEMENT,
                            tag: 'span',
                            content: [
                                {
                                    type: NodeType.TEXT,
                                    tag: 'strong',
                                    content: [
                                        {
                                            type: NodeType.TEXT,
                                            tag: 'text',
                                            text: 'Test Dokument',
                                            content: [],
                                            attrs: {},
                                        },
                                    ],
                                    attrs: {},
                                },
                            ],
                            attrs: {},
                        },
                    ],
                    attrs: {},
                },
                { type: NodeType.ELEMENT, tag: 'hr', content: [], attrs: {} },
            ],
            attrs: {},
        };

        expect(isEmptyRichContentDoc(doc)).toBe(false);
    });
});
