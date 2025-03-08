import type { TocLink } from '@nuxt/content';
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

export function jsonNodeToTocLinks(n: JSONNode): TocLink[] {
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

export function getTextFromContent(n: JSONNode): string {
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

function walkContentForText(ns: JSONNode[]): string {
    let text = '';
    for (let i = 0; i < ns.length; i++) {
        const element = ns[i]!;
        if (element.text === '') {
            text += walkContentForText(element.content);
        } else {
            text += element.text;
            break;
        }
    }

    return text;
}
