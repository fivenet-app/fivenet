import type { TocLink } from '@nuxt/content';
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

export function jsonNodeToTocLinks(n: JSONNode): TocLink[] {
    const headers: TocLink[] = [];
    if (/h[1-6]/i.test(n.tag)) {
        headers.push({
            id: n.id,
            depth: parseInt(n.tag.replace('h', '')),
            text: n.content[0]?.text ?? n.text,
        });
    }

    n.content.forEach((c) => headers.push(...jsonNodeToTocLinks(c)));

    return headers;
}
