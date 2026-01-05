import type { TocLink } from '@nuxt/content';
import type { RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';

export function jsonNodeToTocLinks(n: RichTextHtmlNode): TocLink[] {
    const headers: TocLink[] = [];
    if (/h[1-6]/i.test(n.tag)) {
        const text = getTextFromContent(n);
        headers.push({
            id: n.id ?? n.tag,
            depth: parseInt(n.tag.replace('h', '')),
            text: text,
        });
    }

    n.content.forEach((c) => headers.push(...jsonNodeToTocLinks(c)));

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
