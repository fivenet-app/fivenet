// Based upon https://github.com/harlan-zw/nuxt-seo-kit/blob/8517ccf0fedefdb93af869f813e44a58e944d848/layer/composables/breacrumbs.ts
// Licensed under MIT License © 2022-PRESENT Harlan Wilton, see https://github.com/harlan-zw/nuxt-seo-kit/blob/8517ccf0fedefdb93af869f813e44a58e944d848/README.md#license

import type { ParsedURL } from 'ufo';
import { hasTrailingSlash, parseURL, stringifyParsedURL, withBase, withTrailingSlash, withoutTrailingSlash } from 'ufo';
import type { RouteRecord } from 'vue-router';

export function createInternalLinkResolver(absolute?: boolean) {
    return (path: string) => {
        const fixedSlash = withTrailingSlash(path);
        if (absolute) return withBase(fixedSlash, '/');

        return fixedSlash;
    };
}

function getBreadcrumbs(input: string) {
    const startNode = parseURL(input);
    const appendsTrailingSlash = hasTrailingSlash(startNode.pathname);

    const stepNode = (node: ParsedURL, nodes: string[] = []) => {
        const fullPath = stringifyParsedURL(node);
        // the pathname will always be without the trailing slash
        const currentPathName = node.pathname;
        // when we hit the root the path will be an empty string; we swap it out for a slash
        nodes.push(fullPath || '/');
        // strip the last path segment (/my/cool/path -> /my/cool)
        node.pathname = currentPathName.substring(0, currentPathName.lastIndexOf('/'));
        // if the input was provided with a trailing slash we need to honour that
        if (appendsTrailingSlash) {
            node.pathname = withTrailingSlash(node.pathname.substring(0, node.pathname.lastIndexOf('/')));
        }

        // if we still have a pathname, and it's different, traverse
        if (node.pathname !== currentPathName) stepNode(node, nodes);
        return nodes;
    };
    return stepNode(startNode);
}

type Opts = {
    hideIndex: boolean;
};

const defaultOpts: Opts = {
    hideIndex: true,
};

export function useBreadcrumbs(opts?: Partial<Opts>) {
    if (!opts) {
        opts = defaultOpts;
    }
    const router = useRouter();
    const resolveUrl = createInternalLinkResolver();

    return computed(() => {
        const routes = router.getRoutes();
        const route = router.currentRoute.value;
        const bs = getBreadcrumbs(route.path);
        return bs
            .reverse()
            .map((path, idx, a) => {
                if (a.length - 1 === idx) {
                    return {
                        path,
                        meta: route?.meta,
                    };
                } else {
                    return {
                        path,
                        meta: routes.find(
                            (route: RouteRecord) => withoutTrailingSlash(route.path) === withoutTrailingSlash(path),
                        )?.meta,
                    };
                }
            })
            .filter(({ meta }) => meta !== undefined)
            .filter(({ path }) => path !== '/' && !!opts?.hideIndex)
            .map(({ path, meta }) => {
                // title case string regex
                let title = meta?.title;
                if (!title) {
                    if (path === '/') {
                        title = 'Home';
                        // pop last url segment and title case it
                    } else {
                        title = toTitleCase(withoutTrailingSlash(path).split('/').pop() || '');
                    }
                }

                return {
                    to: resolveUrl(path),
                    title: title,
                    icon: meta?.icon as string | undefined,
                };
            });
    });
}
