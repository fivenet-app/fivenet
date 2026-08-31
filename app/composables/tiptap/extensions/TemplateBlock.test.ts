import { Editor } from '@tiptap/core';
import Document from '@tiptap/extension-document';
import { Paragraph } from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { generateHTML, generateJSON } from '@tiptap/html';
import { describe, expect, it } from 'vitest';
import { classifyTemplateAction, TemplateBlock, TemplateBlockEnd } from './TemplateBlock';
import { isTemplateBlockActionValid } from './TemplateBlockValidation';
import { TemplateVar } from './TemplateVar';

describe('TemplateBlock', () => {
    it('normalizes supported Go block actions into template block nodes', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateBlock, TemplateBlockEnd],
            content: '',
        });

        editor.commands.setContent(
            '<p>{{ range .Users }}{{ if .Active }}{{ else if .Pending }}{{ else }}{{ with .Profile }}{{ break }}{{ continue }}</p>',
            { emitUpdate: false },
        );

        expect(editor.getJSON().content?.[0]?.content).toEqual([
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'range .Users',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'if .Active',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'else if .Pending',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'else',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'with .Profile',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'break',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'continue',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
        ]);

        editor.destroy();
    });

    it('preserves trim markers and normalizes block end actions', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateBlock, TemplateBlockEnd],
            content: '',
        });

        editor.commands.setContent('<p>before {{- if   .Active -}}inside{{- end -}} after</p>', { emitUpdate: false });

        expect(editor.getJSON().content?.[0]?.content).toEqual([
            { type: 'text', text: 'before ' },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'if .Active',
                    'data-left-trim': true,
                    'data-right-trim': true,
                },
            },
            { type: 'text', text: 'inside' },
            {
                type: 'templateBlockEnd',
                attrs: {
                    'data-left-trim': true,
                    'data-right-trim': true,
                },
            },
            { type: 'text', text: ' after' },
        ]);

        editor.destroy();
    });

    it('does not normalize unsupported template actions', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateBlock, TemplateBlockEnd],
            content: '<p>{{ define "item" }}{{ .Name }}{{ template "item" . }}</p>',
        });

        expect(editor.getText()).toBe('{{ define "item" }}{{ .Name }}{{ template "item" . }}');
        expect(editor.getJSON().content?.[0]?.content).toEqual([
            { type: 'text', text: '{{ define "item" }}{{ .Name }}{{ template "item" . }}' },
        ]);

        editor.destroy();
    });

    it('preserves braces in block action arguments', () => {
        const value = "if eq .Status '{pending}'";
        const editor = new Editor({ extensions: [Document, Paragraph, Text, TemplateBlock, TemplateBlockEnd], content: '' });
        editor.commands.setContent(`<p>{{ ${value} }}</p>`, { emitUpdate: false });

        expect(editor.getJSON().content?.[0]?.content?.[0]?.attrs?.['data-template-block']).toBe(value);

        editor.destroy();
    });

    it('normalizes break and continue as blocks when TemplateVar is also enabled', () => {
        const editor = new Editor({
            extensions: [Document, Paragraph, Text, TemplateVar, TemplateBlock, TemplateBlockEnd],
            content: '',
        });

        editor.commands.setContent('<p>{{ break }}{{ continue }}{{ .Name }}</p>', { emitUpdate: false });

        expect(editor.getJSON().content?.[0]?.content).toEqual([
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'break',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateBlock',
                attrs: {
                    'data-template-block': 'continue',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
            {
                type: 'templateVar',
                attrs: {
                    'data-template-var': '.Name',
                    'data-left-trim': false,
                    'data-right-trim': false,
                },
            },
        ]);

        editor.destroy();
    });

    it('includes configured HTML attributes when serializing block actions', () => {
        const extensions = [
            Document,
            Paragraph,
            Text,
            TemplateBlock.configure({ HTMLAttributes: { 'data-template-editor': 'true' } }),
            TemplateBlockEnd.configure({ HTMLAttributes: { 'data-template-editor': 'true' } }),
        ];
        const document = {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        { type: 'templateBlock', attrs: { 'data-template-block': 'if .Active' } },
                        { type: 'templateBlockEnd', attrs: {} },
                    ],
                },
            ],
        };

        const html = generateHTML(document, extensions);

        expect(html.match(/data-template-editor="true"/g)).toHaveLength(2);
    });

    it('round-trips start and end actions through HTML without changing node types', () => {
        const extensions = [Document, Paragraph, Text, TemplateBlock, TemplateBlockEnd];
        const document = {
            type: 'doc',
            content: [
                {
                    type: 'paragraph',
                    content: [
                        {
                            type: 'templateBlock',
                            attrs: {
                                'data-template-block': 'if .Active',
                                'data-left-trim': true,
                                'data-right-trim': false,
                            },
                        },
                        { type: 'text', text: 'Hello' },
                        {
                            type: 'templateBlockEnd',
                            attrs: {
                                'data-left-trim': false,
                                'data-right-trim': true,
                            },
                        },
                    ],
                },
            ],
        };

        const html = generateHTML(document, extensions);
        expect(html).toContain('data-template-block-end="end"');
        expect(html).not.toContain('data-template-block="end"');

        const parsed = generateJSON(html, extensions);
        expect(parsed.content?.[0]?.content).toEqual(document.content?.[0]?.content);
    });

    it.each([
        ['range .Users', 'block-start'],
        ['else if .Pending', 'block-start'],
        ['break', 'block-start'],
        ['continue', 'block-start'],
        ['end', 'block-end'],
        ['/* comment */', 'comment'],
        ['define "item"', 'raw-control'],
        ['template "item" .', 'raw-control'],
        ['end anything', 'raw-control'],
        ['.Firstname | upper', 'variable'],
    ] as const)('classifies %s as %s', (action, kind) => {
        expect(classifyTemplateAction(action)).toBe(kind);
    });

    it.each([
        ['else', true],
        ['else if .Pending', true],
        ['range .Users', true],
        ['if and .Active .Ready', true],
        ['with .Profile', true],
        ['break', true],
        ['continue', true],
        ['end', false],
        ['else arbitrary', false],
        ['else if', false],
        ['range', false],
        ['break .Item', false],
        ['end anything', false],
    ] as const)('validates %s as %s', (action, valid) => {
        expect(isTemplateBlockActionValid(action)).toBe(valid);
    });
});
