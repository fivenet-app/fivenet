import type { TocLink } from '@nuxt/content';
import type { JSONContent } from '@tiptap/core';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import type { Content, RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';

export function jsonNodeToTocLinks(n: Content): TocLink[] {
    if (n.content) {
        return richTextHTMLNodeToTocLinks(n.content);
    }

    if (!n.tiptapJson) return [];

    return tiptapJSONNodeToTocLinks(Struct.toJson(n.tiptapJson) as JSONContent);
}

function richTextHTMLNodeToTocLinks(n: RichTextHtmlNode): TocLink[] {
    const headers: TocLink[] = [];
    if (/h[1-6]/i.test(n.tag)) {
        const text = getTextFromContent(n);
        headers.push({
            id: n.id ?? n.tag,
            depth: parseInt(n.tag.replace('h', '')),
            text: text,
        });
    }

    n.content.forEach((c) => headers.push(...richTextHTMLNodeToTocLinks(c)));

    return headers;
}

function tiptapJSONNodeToTocLinks(n: JSONContent): TocLink[] {
    const headers: TocLink[] = [];
    const type = n.type;
    if (type === 'heading') {
        const level = parseInt(n.attrs?.level ?? '1');
        let text = '';
        n.content?.forEach((c) => {
            if (c.type === 'text') {
                if (c.text) {
                    text += c.text;
                }
            }
        });

        headers.push({
            id: n.attrs?.id ?? `h${level}-${text.toLowerCase().replace(/\s+/g, '-')}`,
            depth: level,
            text: text,
        });
    }

    if (n.content) {
        n.content.forEach((c) => {
            headers.push(...tiptapJSONNodeToTocLinks(c));
        });
    }

    return headers;
}

export function getTextFromContent(n: RichTextHtmlNode): string {
    if (n.text && n.text !== '') {
        return n.text;
    }
    if (n.content.length > 0) {
        const text = walkContentForText(n.content);

        if (text !== '') {
            return text;
        }
    }

    return n.id ?? '';
}

function walkContentForText(ns: RichTextHtmlNode[]): string {
    let text = '';
    for (let i = 0; i < ns.length; i++) {
        const element = ns[i]!;
        if (element.text !== undefined && element.text !== '') {
            text += element.text;
        }

        if (element.content.length > 0) {
            text += walkContentForText(element.content);
        }
    }

    return text;
}
