import { Editor } from '@tiptap/core';
import Document from '@tiptap/extension-document';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { describe, expect, it } from 'vitest';
import { EmbeddedImageGuard, isEmbeddedImageSource } from './EmbeddedImageGuard';
import { EnhancedImage } from './EnhancedImage';

describe('isEmbeddedImageSource', () => {
    it('detects data image URLs regardless of case or leading whitespace', () => {
        expect(isEmbeddedImageSource(' data:image/png;base64,aGVsbG8=')).toBe(true);
        expect(isEmbeddedImageSource('DATA:IMAGE/WEBP;base64,aGVsbG8=')).toBe(true);
    });

    it('allows normal image URLs', () => {
        expect(isEmbeddedImageSource('/api/filestore/images/example.webp')).toBe(false);
        expect(isEmbeddedImageSource('https://example.com/image.webp')).toBe(false);
    });

    it('rejects transactions that introduce an embedded image', () => {
        const editor = new Editor({
            content: '<p>Safe content</p>',
            extensions: [Document, Paragraph, Text, EnhancedImage, EmbeddedImageGuard],
        });

        editor.commands.insertContent({
            type: 'image',
            attrs: { src: 'data:image/png;base64,aGVsbG8=' },
        });

        expect(editor.getHTML()).toBe('<p>Safe content</p>');
        editor.destroy();
    });
});
