import { Editor } from '@tiptap/core';
import Document from '@tiptap/extension-document';
import Link from '@tiptap/extension-link';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { describe, expect, it } from 'vitest';
import AutoLinkIds from './AutoLinkIds';

function createEditor(content: string): Editor {
    const editor = new Editor({
        extensions: [Document, Paragraph, Text, Link, AutoLinkIds],
    });

    editor.commands.setContent(content, { emitUpdate: false });
    return editor;
}

describe('AutoLinkIds', () => {
    it('links valid IDs and leaves punctuation outside the link', () => {
        const editor = createEditor('<p>DOC-1234, CIT-456.</p>');

        expect(editor.getJSON()).toMatchObject({
            content: [
                {
                    content: [
                        { text: 'DOC-1234', marks: [{ type: 'link', attrs: { href: '/documents/1234' } }] },
                        { text: ', ' },
                        { text: 'CIT-456', marks: [{ type: 'link', attrs: { href: '/citizens/456' } }] },
                        { text: '.' },
                    ],
                },
            ],
        });

        editor.destroy();
    });

    it('rejects malformed IDs', () => {
        const editor = createEditor('<p>doc-1234 DOC1234 DOC-12A DOC-1234</p>');

        expect(editor.getHTML()).toBe(
            '<p>doc-1234 DOC1234 DOC-12A <a target="_blank" rel="noopener noreferrer nofollow" href="/documents/1234">DOC-1234</a></p>',
        );

        editor.destroy();
    });

    it('updates and removes generated links while editing', () => {
        const editor = createEditor('<p>DOC-1234</p>');

        editor.commands.setTextSelection({ from: 1, to: 9 });
        editor.commands.insertContent('CIT-456');
        expect(editor.getHTML()).toBe(
            '<p><a target="_blank" rel="noopener noreferrer nofollow" href="/citizens/456">CIT-456</a></p>',
        );

        editor.commands.setTextSelection({ from: 1, to: 8 });
        editor.commands.insertContent('DOC-12A');
        expect(editor.getHTML()).toBe('<p>DOC-12A</p>');

        editor.destroy();
    });
});
