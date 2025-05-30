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
    }>(),
    {
        jobs: () => [],
        users: () => [],
        qualifications: () => [],
        accessTypes: undefined,
        defaultAccess: 2, // All `AccessLevel` should have 0 = UNSPECIFIED, 1 = BLOCKED, 2 = VIEW
        disabled: false,
        showRequired: false,
        defaultAccessType: 'job',
        totalLimit: undefined,
    },
);

const { maxAccessEntries } = useAppConfig();

const maxEntries = computed(() => props.totalLimit || maxAccessEntries);

const emit = defineEmits<{
    (e: 'update:jobs', jobs: JobsT[]): void;
    (e: 'update:users', users: UsersT[]): void;
    (e: 'update:qualifications', qualifications: QualiT[]): void;
}>();

const { t } = useI18n();

const jobsAccess = useVModel(props, 'jobs', emit, {
    deep: true,
});
const usersAccess = useVModel(props, 'users', emit, {
    deep: true,
});
const qualificationsAccess = useVModel(props, 'qualifications', emit, {
    deep: true,
});

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

watch(
    access,
    () => {
        usersAccess.value.length = 0;
        jobsAccess.value.length = 0;
        qualificationsAccess.value.length = 0;

        access.value.forEach((e) => {
            if (e.type === 'user') {
                if (e.userId === undefined) {
                    return;
                }

                const idx = usersAccess.value.findIndex((a) => a.id === e.id);
                if (idx === -1) {
                    usersAccess.value.push({
                        id: e.id,
                        targetId: props.targetId,
                        userId: e.userId,
                        access: e.access,
                        required: e.required,
                    } as UsersT);
                }
            } else if (e.type === 'job') {
                if (e.job === undefined || e.minimumGrade === undefined) {
                    return;
                }

                const idx = jobsAccess.value.findIndex((a) => a.id === e.id);
                if (idx === -1) {
                    jobsAccess.value.push({
                        id: e.id,
                        targetId: props.targetId,
                        job: e.job,
                        minimumGrade: e.minimumGrade,
                        access: e.access,
                        required: e.required,
                    } as JobsT);
                }
            } else if (e.type === 'qualification') {
                if (e.qualificationId === undefined) {
                    return;
                }

                const idx = qualificationsAccess.value.findIndex((a) => a.id === e.id);
                if (idx === -1) {
                    qualificationsAccess.value.push({
                        id: e.id,
                        targetId: props.targetId,
                        qualificationId: e.qualificationId,
                        access: e.access,
                        required: e.required,
                        qualification: e.qualification,
                    } as QualiT);
                }
            }
        });
    },
    {
        deep: true,
    },
);

const lastId = ref(0);

function setFromPropsJobs(): void {
    access.value?.push(
        ...jobsAccess.value
            .filter((a) => !access.value.find((ac) => ac.id === a.id))
            .map((a) => {
                if (a.id === 0) {
                    a.id = lastId.value;
                    lastId.value++;
                }
                return a;
            })
            .map((a) => ({ ...a, type: 'job' }) as MixedAccessEntry),
    );
}

function setFromPropsUsers(): void {
    access.value?.push(
        ...usersAccess.value
            .filter((a) => !access.value.find((ac) => ac.id === a.id))
            .map((a) => {
                if (a.id === 0) {
                    a.id = lastId.value;
                    lastId.value++;
                }
                return a;
            })
            .map((a) => ({ ...a, type: 'user' }) as MixedAccessEntry),
    );
}

function setFromPropsQualifications(): void {
    access.value?.push(
        ...usersAccess.value
            .filter((a) => !access.value.find((ac) => ac.id === a.id))
            .map((a) => {
                if (a.id === 0) {
                    a.id = lastId.value;
                    lastId.value++;
                }
                return a;
            })
            .map((a) => ({ ...a, type: 'user' }) as MixedAccessEntry),
    );
}

onBeforeMount(() => {
    setFromPropsJobs();
    setFromPropsUsers();
    setFromPropsQualifications();
});

watch(jobsAccess, setFromPropsJobs);
watch(usersAccess, setFromPropsUsers);
watch(qualificationsAccess, setFromPropsQualifications);

function addNewEntry(): void {
    let idx = aTypes.value.findIndex((at) => at.type === props.defaultAccessType);
    if (idx === -1) {
        idx = aTypes.value.length - 1;
    }

    access.value.push({
        id: lastId.value,
        type: aTypes.value[idx]?.type ?? 'job',
        access: props.defaultAccess,
    });
    lastId.value++;
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
            @delete="access?.splice(idx, 1)"
        />

        <UTooltip v-if="!disabled" :text="$t('components.access.add_entry')">
            <UButton
                :disabled="access.length >= maxEntries"
                :ui="{ rounded: 'rounded-full' }"
                icon="i-mdi-plus"
                @click="addNewEntry"
            />
        </UTooltip>
    </div>
</template>
