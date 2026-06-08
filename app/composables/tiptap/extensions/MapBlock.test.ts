import Document from '@tiptap/extension-document';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { generateHTML, generateJSON } from '@tiptap/html';
import { describe, expect, it } from 'vitest';
import { MapBlock } from './MapBlock';

describe('MapBlock', () => {
    const extensions = [Document, Paragraph, Text, MapBlock];

    it('round-trips html with data attributes', () => {
        const doc = {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        {
                            type: 'mapBlock',
                            attrs: {
                                x: 12.34,
                                y: 56.78,
                                zoom: 3,
                                postal: '12345',
                                layer: 'postal',
                            },
                        },
                    ],
                },
            ],
        };

        const html = generateHTML(doc, extensions);
        expect(html).toContain('<span');
        expect(html).toContain('data-embed="map"');
        expect(html).toContain('data-map-x="12.34"');
        expect(html).toContain('data-map-y="56.78"');
        expect(html).toContain('data-map-zoom="3"');
        expect(html).toContain('data-map-postal="12345"');
        expect(html).toContain('data-map-layer="postal"');

        const parsed = generateJSON(html, extensions);
        expect(parsed).toEqual(doc);
    });

    it('parses a minimal map wrapper', () => {
        const html = '<span data-embed="map" data-map-x="1.5" data-map-y="2.5" data-map-zoom="4"></span>';
        const parsed = generateJSON(html, extensions);

        expect(parsed).toEqual({
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        {
                            type: 'mapBlock',
                            attrs: {
                                x: 1.5,
                                y: 2.5,
                                zoom: 4,
                                postal: '',
                                layer: 'postal',
                            },
                        },
                    ],
                },
            ],
        });
    });
});
