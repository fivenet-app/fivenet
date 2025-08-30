<script
    lang="ts"
    setup
    generic="JobsT extends JobAccessEntry, UsersT extends UserAccessEntry, QualiT extends QualificationAccessEntry"
>
import { useCompletorStore } from '~/stores/completor';
import AccessEntry from './AccessEntry.vue';
import type {
    AccessEntryType,
    AccessLevelEnum,
    AccessType,
    JobAccessEntry,
    MixedAccessEntry,
    QualificationAccessEntry,
    UserAccessEntry,
} from './helpers';

const props = withDefaults(
    defineProps<{
        targetId: number;
        jobs?: JobsT[];
        users?: UsersT[];
        qualifications?: QualiT[];
        accessTypes?: AccessType[];
        accessRoles: AccessLevelEnum[];
        defaultAccess?: number;
        disabled?: boolean;
        showRequired?: boolean;
        defaultAccessType?: AccessEntryType;
        totalLimit?: number;
        hideGrade?: boolean;
        hideJobs?: string[];
    }>(),
    {
        jobs: () => [],
        users: () => [],
        qualifications: () => [],
        accessTypes: undefined,
        defaultAccess: 2, // All `AccessLevel` should have 0 = UNSPECIFIED, 1 = BLOCKED, 2 = VIEW levels
        disabled: false,
        showRequired: false,
        defaultAccessType: 'job',
        totalLimit: undefined,
        hideGrade: false,
        hideJobs: () => [],
    },
);

const { maxAccessEntries } = useAppConfig();

const maxEntries = computed(() => props.totalLimit || maxAccessEntries);

const { t } = useI18n();

const jobsAccess = defineModel<JobsT[]>('jobs', { default: () => [] });
const usersAccess = defineModel<UsersT[]>('users', { default: () => [] });
const qualificationsAccess = defineModel<QualiT[]>('qualifications', { default: () => [] });

const defaultAccessTypes = [
    { label: t('common.citizen', 2), value: 'user' },
    { label: t('common.job', 2), value: 'job' },
] as AccessType[];

const aTypes = ref<AccessType[]>([]);
if (props.accessTypes === undefined) {
    aTypes.value = defaultAccessTypes;
} else {
    aTypes.value = props.accessTypes;
}

const access = ref<MixedAccessEntry[]>([]);

function isEqualArray(a: MixedAccessEntry[], b: MixedAccessEntry[]): boolean {
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) {
        const entryA = a[i]!;

        const entryB = b.find((item) => item.id === entryA.id);
        if (!entryB) {
            return false;
        }

        if (
            entryA.id !== entryB.id ||
            entryA.type !== entryB.type ||
            entryA.access !== entryB.access ||
            entryA.required !== entryB.required ||
            entryA.job !== entryB.job ||
            entryA.minimumGrade !== entryB.minimumGrade ||
            entryA.userId !== entryB.userId ||
            entryA.qualificationId !== entryB.qualificationId ||
            entryA.qualification !== entryB.qualification
        ) {
            return false;
        }
    }
    return true;
}

function syncAccessFromProps() {
    // Merge all entries from jobs, users, and qualifications into access
    const merged: MixedAccessEntry[] = [];

    jobsAccess.value.forEach((a) => {
        if (!a.id) {
            // Assign a new ID if none is set
            a.id = lastId--;
        }
        merged.push({ ...a, type: 'job' });
    });
    usersAccess.value.forEach((a) => {
        if (!a.id) {
            // Assign a new ID if none is set
            a.id = lastId--;
        }
        merged.push({ ...a, type: 'user' });
    });
    qualificationsAccess.value.forEach((a) => {
        if (!a.id) {
            // Assign a new ID if none is set
            a.id = lastId--;
        }
        merged.push({ ...a, type: 'qualification' });
    });

    const newAccess = merged.map((entry) => ({
        ...entry,
        id: entry.id,
    }));

    if (!isEqualArray(access.value, newAccess)) {
        access.value = newAccess;
    }
}

// Helper to update a reactive array in-place (add, update, remove)
function syncArray<T extends { id: number }>(source: T[], target: T[]) {
    // Remove items not in target
    for (let i = source.length - 1; i >= 0; i--) {
        if (!target.find((t) => t.id === source[i]!.id)) {
            source.splice(i, 1);
        }
    }
    // Add or update items
    target.forEach((t) => {
        const idx = source.findIndex((s) => s.id === t.id);
        if (idx === -1) {
            source.push(t);
        } else {
            // Update properties
            for (const key in t) {
                if (Object.prototype.hasOwnProperty.call(t, key) && source[idx]![key] !== t[key]) {
                    source[idx]![key] = t[key];
                }
            }
        }
    });
}

function syncPropsFromAccess() {
    const newJobs = access.value
        .filter((e) => e.type === 'job')
        .map((e) => ({
            id: e.id,
            targetId: props.targetId,
            job: e.job,
            minimumGrade: e.minimumGrade,
            access: e.access,
            required: e.required,
        })) as JobsT[];

    const newUsers = access.value
        .filter((e) => e.type === 'user')
        .map((e) => ({
            id: e.id,
            targetId: props.targetId,
            userId: e.userId,
            access: e.access,
            required: e.required,
        })) as UsersT[];

    const newQualis = access.value
        .filter((e) => e.type === 'qualification')
        .map((e) => ({
            id: e.id,
            targetId: props.targetId,
            qualificationId: e.qualificationId,
            access: e.access,
            required: e.required,
            qualification: e.qualification,
        })) as QualiT[];

    syncArray(jobsAccess.value, newJobs);
    syncArray(usersAccess.value, newUsers);
    syncArray(qualificationsAccess.value, newQualis);
}

// Sync from props on mount and when any prop changes
onBeforeMount(() => syncAccessFromProps());
watch([jobsAccess, usersAccess, qualificationsAccess], syncAccessFromProps, { deep: true });

// Sync from access to props when access changes
watch(access, syncPropsFromAccess, { deep: true });

let lastId = 0;

function addNewEntry(): void {
    let idx = aTypes.value.findIndex((at) => at.value === props.defaultAccessType);
    if (idx === -1) idx = aTypes.value.length - 1;

    access.value.push({
        id: lastId--,
        type: aTypes.value[idx]?.value ?? 'job',
        access: props.defaultAccess,
    });
}

const completorStore = useCompletorStore();
const { data: jobsList } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <div>
        <AccessEntry
            v-for="(entry, idx) in access"
            :key="entry.id"
            v-model="access[idx]!"
            :access-types="aTypes"
            :access-roles="accessRoles"
            :disabled="disabled"
            :show-required="showRequired"
            :jobs="jobsList"
            v-bind="$attrs"
            :hide-grade="hideGrade"
            :hide-jobs="hideJobs"
            @delete="access?.splice(idx, 1)"
        />

        <UTooltip v-if="!disabled" :text="$t('components.access.add_entry')">
            <UButton :disabled="access.length >= maxEntries" icon="i-mdi-plus" @click="addNewEntry" />
        </UTooltip>
    </div>
</template>
