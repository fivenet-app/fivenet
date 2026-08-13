import { DECORATION_MANAGER_PLUGIN_KEY, Editor } from '@tiptap/core';
import Document from '@tiptap/extension-document';
import Paragraph from '@tiptap/extension-paragraph';
import Text from '@tiptap/extension-text';
import { describe, expect, it } from 'vitest';
import SearchAndReplace from './SearchAndReplace';

function createEditor(content: string) {
    return new Editor({
        extensions: [Document, Paragraph, Text, SearchAndReplace],
        content,
    });
}

function getDecorationClasses(editor: Editor): string[] {
    const state = DECORATION_MANAGER_PLUGIN_KEY.getState(editor.state);
    const decorations = state?.mergedDecorationSet.find() ?? [];

    return decorations.map((decoration) => (decoration.type.attrs.class as string | undefined) ?? '');
}

describe('SearchAndReplace decorations', () => {
    it('recomputes highlights when the search term changes', () => {
        const editor = createEditor('<p>hello world hello</p>');

        expect(editor.storage.searchAndReplace.results).toHaveLength(0);
        expect(getDecorationClasses(editor)).toHaveLength(0);

        editor.commands.setSearchTerm('hello');

        expect(editor.storage.searchAndReplace.results).toHaveLength(2);
        expect(getDecorationClasses(editor)).toEqual(['search-result search-result-current', 'search-result']);

        editor.destroy();
    });

    it('updates the current result decoration when navigating results', () => {
        const editor = createEditor('<p>hello world hello</p>');

        editor.commands.setSearchTerm('hello');
        editor.commands.nextSearchResult();

        expect(editor.storage.searchAndReplace.resultIndex).toBe(1);
        expect(getDecorationClasses(editor)).toEqual(['search-result', 'search-result search-result-current']);

        editor.destroy();
    });

    it('clears highlights when the search term is cleared', () => {
        const editor = createEditor('<p>hello world hello</p>');

        editor.commands.setSearchTerm('hello');
        expect(getDecorationClasses(editor)).toHaveLength(2);

        editor.commands.resetIndex();
        editor.commands.setSearchTerm('');

        expect(editor.storage.searchAndReplace.results).toHaveLength(0);
        expect(getDecorationClasses(editor)).toHaveLength(0);

        editor.destroy();
    });
});
