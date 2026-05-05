import type { Permission } from '~~/gen/ts/resources/permissions/permissions/permissions';

export type PermissionServiceGroup = {
    namespace: string;
    service: string;
    icon?: string;
    order: number;
};

export type PermissionNamespaceGroup = {
    namespace: string;
    icon?: string;
    order: number;
    services: PermissionServiceGroup[];
};

type Translate = (key: string) => string;
type TranslationExists = (key: string) => boolean;

export function buildPermissionGroups(permissions: Permission[]): PermissionNamespaceGroup[] {
    const namespaceMap = new Map<string, PermissionNamespaceGroup & { serviceMap: Map<string, PermissionServiceGroup> }>();

    permissions.forEach((perm) => {
        const namespaceGroup = namespaceMap.get(perm.namespace) ?? {
            namespace: perm.namespace,
            icon: perm.icon,
            order: perm.order ?? 999999,
            services: [],
            serviceMap: new Map<string, PermissionServiceGroup>(),
        };

        namespaceGroup.order = Math.min(namespaceGroup.order, perm.order ?? 999999);
        namespaceGroup.icon ??= perm.icon;

        const serviceGroup = namespaceGroup.serviceMap.get(perm.service) ?? {
            namespace: perm.namespace,
            service: perm.service,
            icon: perm.icon,
            order: perm.order ?? 999999,
        };

        serviceGroup.order = Math.min(serviceGroup.order, perm.order ?? 999999);
        serviceGroup.icon ??= perm.icon;
        namespaceGroup.serviceMap.set(perm.service, serviceGroup);
        namespaceMap.set(perm.namespace, namespaceGroup);
    });

    return [...namespaceMap.values()]
        .map((namespaceGroup) => {
            const { serviceMap, ...group } = namespaceGroup;

            return {
                ...group,
                services: [...serviceMap.values()].sort((a, b) => a.order - b.order),
            };
        })
        .sort((a, b) => a.order - b.order);
}

export function getPermissionServiceLabel(namespace: string, service: string, t: Translate, te: TranslationExists): string {
    const serviceKey = `perms.${namespace}.${service}.service`;
    if (te(serviceKey)) return t(serviceKey);

    const serviceNamespaceKey = `perms.${namespace}.namespace`;
    if (te(serviceNamespaceKey)) return t(serviceNamespaceKey);

    return service;
}

export function getPermissionNamespaceLabel(namespace: string, t: Translate, te: TranslationExists): string {
    const namespaceKey = `perms.${namespace}.namespace`;
    if (te(namespaceKey)) return t(namespaceKey);

    return namespace.charAt(0).toUpperCase() + namespace.slice(1);
}
