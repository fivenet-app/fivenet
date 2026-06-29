import type { ComputedRef } from 'vue';
import { useAuthStore } from '~/stores/auth';
import slug from '~/utils/slugify';
import type { PermAttrKey, PermAttrKeysByType, PermAttrValue, Perms } from '~~/gen/ts/perms';
import type { JobGrades, RoleAttribute } from '~~/gen/ts/resources/permissions/attributes/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions/permissions';

export type canMode = 'oneof' | 'all';

// Wrapper around auth store to make live easier
const _useAuth = () => {
    const authStore = useAuthStore();
    const { username, accountId, activeChar, isSuperuser, jobProps, permissions, attributes, canBeSuperuser } =
        storeToRefs(authStore);

    function checkPerm(permissions: Permission[], perm: string | string[], mode: canMode = 'oneof'): boolean {
        if (mode === undefined) mode = 'oneof';

        if (permissions.find((p) => p.guardName === 'superuser')) return true;

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

    type JobGradeListValue = {
        fineGrained: boolean;
        jobs: {
            [key: string]: number;
        };
        grades: {
            [key: string]: JobGrades;
        };
    };

    function getAttr<P extends Perms, K extends PermAttrKey<P>>(perm: P, key: K): ComputedRef<RoleAttribute | undefined>;
    function getAttr(perm: Perms, key: string): ComputedRef<RoleAttribute | undefined>;
    function getAttr(perm: Perms, key: string): ComputedRef<RoleAttribute | undefined> {
        return computed(() => {
            const [serviceKey, name] = perm.split('/');
            const [namespace, service] = serviceKey?.split('.') ?? [];

            return attributes.value.find(
                (a) => a.namespace === namespace && a.service === service && a.name === name && a.key === key,
            );
        });
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

    function attr<P extends Perms, K extends PermAttrKey<P>>(perm: P, key: K, val: PermAttrValue<P, K>): ComputedRef<boolean>;
    function attr(perm: Perms, key: string, val: string): ComputedRef<boolean>;
    function attr(perm: Perms, key: string, val: unknown): ComputedRef<boolean> {
        return computed(() => {
            if (isSuperuser.value === true) return true;

            const attrValue = val as string;
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'stringList') {
                return a.value.validValues.stringList.strings.includes(attrValue);
            } else if (a?.value?.validValues.oneofKind === 'jobList') {
                return a.value.validValues.jobList.strings.includes(attrValue);
            } else if (a?.value?.validValues.oneofKind === 'jobGradeList') {
                return (
                    a.value.validValues.jobGradeList.jobs[attrValue] !== undefined ||
                    a.value.validValues.jobGradeList.grades[attrValue] !== undefined
                );
            }

            return false;
        });
    }

    function attrStringList<P extends Perms, K extends PermAttrKeysByType<P, 'stringList'>>(
        perm: P,
        key: K,
    ): ComputedRef<ReadonlyArray<PermAttrValue<P, K>>>;
    function attrStringList(perm: Perms, key: string): ComputedRef<ReadonlyArray<string>>;
    function attrStringList(perm: Perms, key: string): ComputedRef<ReadonlyArray<string>> {
        return computed(() => {
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'stringList') {
                return a.value.validValues.stringList.strings;
            }
            return [];
        });
    }

    function attrJobList<P extends Perms, K extends PermAttrKeysByType<P, 'jobList'>>(
        perm: P,
        key: K,
    ): ComputedRef<ReadonlyArray<PermAttrValue<P, K>>>;
    function attrJobList(perm: Perms, key: string): ComputedRef<ReadonlyArray<string>>;
    function attrJobList(perm: Perms, key: string): ComputedRef<ReadonlyArray<string>> {
        return computed(() => {
            const a = getAttr(perm, key).value;

            if (a?.value?.validValues.oneofKind === 'jobList') {
                return a.value.validValues.jobList.strings;
            }
            return [];
        });
    }

    function attrJobGradeList<P extends Perms, K extends PermAttrKeysByType<P, 'jobGradeList'>>(
        perm: P,
        key: K,
    ): ComputedRef<JobGradeListValue>;
    function attrJobGradeList(perm: Perms, key: string): ComputedRef<JobGradeListValue>;
    function attrJobGradeList(perm: Perms, key: string): ComputedRef<JobGradeListValue> {
        return computed(() => {
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
    }

    return {
        // Getters
        accountId,
        activeChar,
        canBeSuperuser,
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
