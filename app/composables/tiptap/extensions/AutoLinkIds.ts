import { Extension } from '@tiptap/core';
import type { Node as PMNode } from '@tiptap/pm/model';
import { Plugin, PluginKey, type Transaction } from '@tiptap/pm/state';

export interface AutoLinkIdType {
    pattern: RegExp;
    href: (id: string) => string;
    hrefPattern: RegExp;
}

/** The supported public IDs and their deterministic application routes. */
export const autoLinkIdTypes = {
    DOC: {
        pattern: /(?<![A-Za-z0-9_])DOC-(\d+)(?![A-Za-z0-9_])/g,
        href: (id: string) => `/documents/${id}`,
        hrefPattern: /^\/documents\/(\d+)$/,
    },
    CIT: {
        pattern: /(?<![A-Za-z0-9_])CIT-(\d+)(?![A-Za-z0-9_])/g,
        href: (id: string) => `/citizens/${id}`,
        hrefPattern: /^\/citizens\/(\d+)$/,
    },
} satisfies Record<string, AutoLinkIdType>;

const autoLinkIdsPluginKey = new PluginKey('autoLinkIds');

function getIdMatches(text: string): Array<{ from: number; to: number; href: string }> {
    const matches: Array<{ from: number; to: number; href: string }> = [];

    for (const type of Object.values(autoLinkIdTypes)) {
        type.pattern.lastIndex = 0;
        for (const match of text.matchAll(type.pattern)) {
            const id = match[1];
            const value = match[0];
            if (!id || !value || match.index === undefined) continue;

            matches.push({
                from: match.index,
                to: match.index + value.length,
                href: type.href(id),
            });
        }
    }

    return matches.sort((a, b) => a.from - b.from);
}

function isGeneratedHref(href: unknown): boolean {
    return typeof href === 'string' && Object.values(autoLinkIdTypes).some((type) => type.hrefPattern.test(href));
}

function synchronizeLinks(state: { doc: PMNode; tr: Transaction }): Transaction | undefined {
    const linkType = state.doc.type.schema.marks.link;
    if (!linkType) return undefined;

    const tr = state.tr;
    let changed = false;

    state.doc.descendants((node, pos) => {
        if (!node.isText || !node.text) return true;

        const generatedLinks = node.marks.filter((mark) => mark.type === linkType && isGeneratedHref(mark.attrs.href));
        generatedLinks.forEach((mark) => {
            tr.removeMark(pos, pos + node.nodeSize, mark);
            changed = true;
        });

        getIdMatches(node.text).forEach(({ from, to, href }) => {
            tr.addMark(pos + from, pos + to, linkType.create({ href }));
            changed = true;
        });

        return true;
    });

    return changed ? tr.setMeta(autoLinkIdsPluginKey, true) : undefined;
}

const AutoLinkIds = Extension.create({
    name: 'autoLinkIds',

    onCreate({ editor }) {
        const transaction = synchronizeLinks({ doc: editor.state.doc, tr: editor.state.tr });
        if (transaction?.docChanged) editor.view.dispatch(transaction);
    },

    addProseMirrorPlugins() {
        return [
            new Plugin({
                key: autoLinkIdsPluginKey,
                appendTransaction: (transactions, _oldState, newState) => {
                    if (!transactions.some((transaction) => transaction.docChanged)) return undefined;
                    if (transactions.some((transaction) => transaction.getMeta(autoLinkIdsPluginKey))) return undefined;

                    return synchronizeLinks({ doc: newState.doc, tr: newState.tr });
                },
            }),
        ];
    },
});

export default AutoLinkIds;
