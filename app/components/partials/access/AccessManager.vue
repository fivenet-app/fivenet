<script lang="ts" setup generic="JobsT extends JobAccessEntry, UsersT extends UserAccessEntry">
import { useCompletorStore } from '~/store/completor';
import AccessEntry from './AccessEntry.vue';
import type { AccessLevelEnum, AccessType, JobAccessEntry, MixedAccessEntry, UserAccessEntry } from './helpers';

const props = withDefaults(
    defineProps<{
        targetId: string;
        jobs?: JobsT[];
        users?: UsersT[];
        accessTypes?: AccessType[];
        accessRoles?: AccessLevelEnum[];
        defaultAccess?: number;
        disabled?: boolean;
        showRequired?: boolean;
    }>(),
    {
        jobs: () => [],
        users: () => [],
        accessTypes: undefined,
        accessRoles: () => [],
        defaultAccess: 2, // All `AccessLevel` should have 0 = UNSPECIFIED, 1 = BLOCKED, 2 = VIEW
        disabled: false,
        showRequired: false,
    },
);

const emits = defineEmits<{
    (e: 'update:jobs', jobs: JobsT): void;
    (e: 'update:users', jobs: UsersT): void;
}>();

const { t } = useI18n();

const { jobs: jobsAccess, users: usersAccess } = useVModels(props, emits);

const defaultAccessTypes = [
    { type: 'user', name: t('common.citizen', 2) },
    { type: 'job', name: t('common.job', 2) },
] as AccessType[];

const aTypes = ref<AccessType[]>([]);
if (props.accessTypes === undefined) {
    aTypes.value = defaultAccessTypes;
} else {
    aTypes.value = props.accessTypes;
}

const access = ref<MixedAccessEntry[]>([]);
watchArray(
    access,
    (newList, _, added, removed) => {
        added.forEach((e) => {
            if (e.type === 'user') {
                usersAccess.value.push({
                    id: e.id,
                    userId: e.userId,
                    access: e.access,
                    required: e.required,
                } as UsersT);
            } else if (e.type === 'job') {
                jobsAccess.value.push({
                    id: e.id,
                    job: e.job,
                    minimumGrade: e.minimumGrade,
                    access: e.access,
                    required: e.required,
                } as JobsT);
            }
        });

        removed.forEach((e) => {
            if (e.type === 'user') {
                const idx = usersAccess.value.findIndex((ua) => ua.id === e.id);
                if (idx > -1) {
                    usersAccess.value.splice(idx, 1);
                }
            } else if (e.type === 'job') {
                const idx = jobsAccess.value.findIndex((ja) => ja.id === e.id);
                if (idx > -1) {
                    jobsAccess.value.splice(idx, 1);
                }
            }
        });

        newList.forEach((e) => {
            const uaIdx = usersAccess.value.findIndex((ua) => ua.id === e.id);
            const jaIdx = jobsAccess.value.findIndex((ua) => ua.id === e.id);

            if (jaIdx > -1 && e.type === 'user') {
                jobsAccess.value.splice(jaIdx, 1);

                usersAccess.value.push({
                    id: e.id,
                    access: e.access,
                    userId: e.userId,
                    required: e.required,
                } as UsersT);
                return;
            } else if (uaIdx > -1 && e.type === 'job') {
                usersAccess.value.splice(uaIdx, 1);

                jobsAccess.value.push({
                    id: e.id,
                    access: e.access,
                    job: e.job,
                    minimumGrade: e.minimumGrade,
                    required: e.required,
                } as JobsT);
                return;
            }

            if (uaIdx > -1 && usersAccess.value[uaIdx]) {
                if (jaIdx > -1) {
                    jobsAccess.value.splice(jaIdx, 1);
                }

                usersAccess.value[uaIdx].userId = e.userId;
                usersAccess.value[uaIdx].access = e.access;
                usersAccess.value[uaIdx].required = e.required;
            } else if (jaIdx > -1 && jobsAccess.value[uaIdx]) {
                if (uaIdx > -1) {
                    usersAccess.value.splice(uaIdx, 1);
                }

                jobsAccess.value[uaIdx].job = e.job;
                jobsAccess.value[uaIdx].minimumGrade = e.minimumGrade;
                jobsAccess.value[uaIdx].userId = undefined;
                jobsAccess.value[uaIdx].access = e.access;
                jobsAccess.value[uaIdx].required = e.required;
            }
        });
    },
    {
        deep: true,
    },
);

watchOnce(jobsAccess, () => access.value?.push(...jobsAccess.value.map((a) => ({ ...a, type: 'job' }) as MixedAccessEntry)));
watchOnce(usersAccess, () => access.value?.push(...usersAccess.value.map((a) => ({ ...a, type: 'user' }) as MixedAccessEntry)));

const lastId = ref(0);
function addEntry(): void {
    access.value.push({
        id: lastId.value.toString(),
        type: aTypes.value[aTypes.value.length - 1]?.type ?? 'job',
        access: props.defaultAccess,
    });
    lastId.value--;
}

const completorStore = useCompletorStore();
const { data: jobsList } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <div>
        <AccessEntry
            v-for="(entry, idx) in access"
            :key="entry.id"
            v-bind="$attrs"
            v-model="access[idx]!"
            :access-types="aTypes"
            :access-roles="accessRoles"
            :disabled="disabled"
            :show-required="showRequired"
            :jobs="jobsList"
            @delete="access?.splice(idx, 1)"
        />

        <UButton
            :ui="{ rounded: 'rounded-full' }"
            icon="i-mdi-plus"
            :title="$t('components.documents.document_editor.add_permission')"
            @click="addEntry"
        />
    </div>
</template>
