import { Editor } from '@tiptap/core';
import Document from '@tiptap/extension-document';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { generateHTML, generateJSON } from '@tiptap/html';
import { describe, expect, it } from 'vitest';
import { TemplateBlock, TemplateBlockEnd } from './TemplateBlock';
import { TemplateVar } from './TemplateVar';

describe('TemplateVar', () => {
    it('normalizes raw Go template actions into template variable nodes', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateVar],
            content: '',
        });

        editor.commands.setContent('<p>Hello {{ .Firstname }} {{- .Lastname -}}</p>', { emitUpdate: false });

        expect(editor.getJSON()).toEqual({
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        { type: 'text', text: 'Hello ' },
                        {
                            type: 'templateVar',
                            attrs: {
                                'data-template-var': '.Firstname',
                                'data-left-trim': false,
                                'data-right-trim': false,
                            },
                        },
                        { type: 'text', text: ' ' },
                        {
                            type: 'templateVar',
                            attrs: {
                                'data-template-var': '.Lastname',
                                'data-left-trim': true,
                                'data-right-trim': true,
                            },
                        },
                    ],
                },
            ],
        });

        editor.destroy();
    });

    it('does not normalize Go control actions as variables', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateVar],
            content: '<p>{{ range .Users }}Hello{{ end }}</p>',
        });

        expect(editor.getText()).toBe('{{ range .Users }}Hello{{ end }}');
        expect(editor.getJSON().content?.[0]?.content).toEqual([{ type: 'text', text: '{{ range .Users }}Hello{{ end }}' }]);

        editor.destroy();
    });

    it('leaves raw control actions and comments available for advanced users', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateVar, TemplateBlock, TemplateBlockEnd],
            content: '',
        });

        const content = '{{ define "item" }}{{ template "item" . }}{{ block "item" . }}{{/* comment */}}';
        editor.commands.setContent(`<p>${content}</p>`, { emitUpdate: false });

        expect(editor.getText()).toBe(content);
        expect(editor.getJSON().content?.[0]?.content).toEqual([{ type: 'text', text: content }]);

        editor.destroy();
    });

    it('round-trips template variables through HTML with their attributes intact', () => {
        const extensions = [Document, Paragraph, Text, TemplateVar];
        const document = {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        {
                            type: 'templateVar',
                            attrs: {
                                'data-template-var': '.Firstname',
                                'data-left-trim': true,
                                'data-right-trim': true,
                            },
                        },
                    ],
                },
            ],
        };

        const html = generateHTML(document, extensions);
        const parsed = generateJSON(html, extensions);

        expect(parsed.content?.[0]?.content).toEqual(document.content?.[0]?.content);
    });
});
