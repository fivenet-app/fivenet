import { useAuthStore } from '~/stores/auth';
import slug from '~/utils/slugify';
import type { Perms } from '~~/gen/ts/perms';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';

export type canMode = 'oneof' | 'all';

// Wrapper around auth store to make live easier
const _useAuth = () => {
    const authStore = useAuthStore();
    const { activeChar, attributes, isSuperuser, jobProps, username, permissions } = storeToRefs(authStore);

    function checkPerm(permissions: Permission[], perm: string | string[], mode: canMode = 'oneof'): boolean {
        if (mode === undefined) {
            mode = 'oneof';
        }

        if (permissions.find((p) => p.guardName === 'superuser')) {
            return true;
        }

        const input: string[] = [];
        if (typeof perm === 'string') {
            input.push(perm.replaceAll('/', '.'));
        } else {
            const vals = perm as string[];
            input.push(...vals.map((v) => v.replaceAll('/', '.')));
        }

        let ok = false;
        // Iterate over permissions and check in "OR" condition manner
        for (let idx = 0; idx < input.length; idx++) {
            const val = slug(input[idx] as string);
            if (permissions.find((p) => p.guardName === val) || val === '') {
                // Permission found
                if (mode === 'oneof') {
                    return true;
                }

                ok = true;
            } else if (mode === 'all') {
                // Permission not found and mode requires all to be found
                return false;
            }
        }

        return ok;
    }

    /**
     * @param perm one or more perms to check
     * @param mode default 'oneof'
     * @returns boolean
     */
    const can = (perm: Perms | Perms[], mode: canMode = 'oneof') => {
        return computed(() => {
            if (isSuperuser.value === true) return true;

            return checkPerm(permissions.value, perm, mode);
        });
    };

    const getAttr = (perm: Perms, key: string) =>
        computed(() => {
            const split = perm.split('/');
            return attributes.value.find((a) => a.category === split[0]! && a.name === split[1]! && a.key === key);
        });

    const attr = (perm: Perms, key: string, val: string) =>
        computed(() => {
            if (isSuperuser.value === true) return true;

            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'stringList') {
                return a.value.validValues.stringList.strings.includes(val);
            } else if (a?.value?.validValues.oneofKind === 'jobList') {
                return a.value.validValues.jobList.strings.includes(val);
            }

            return false;
        });

    const attrStringList = (perm: Perms, key: string) =>
        computed(() => {
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'stringList') {
                return a.value.validValues.stringList.strings;
            }
            return [];
        });

    const attrJobList = (perm: Perms, key: string) =>
        computed(() => {
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'jobList') {
                return a.value.validValues.jobList.strings;
            }
            return [];
        });

    const attrJobGradeList = (perm: Perms, key: string) =>
        computed(() => {
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'jobGradeList') {
                return a.value.validValues.jobGradeList;
            }

            return {
                fineGrained: false,
                jobs: {},
                grades: {},
            };
        });

    return {
        // Getters
        activeChar,
        isSuperuser,
        jobProps,
        username,

        // Funcs
        can,
        attr,

        attrStringList,
        attrJobList,
        attrJobGradeList,
    };
};

export const useAuth = createSharedComposable(_useAuth);
