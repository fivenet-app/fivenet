<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { useCompletorStore } from '~/store/completor';
import type { ArrayElement } from '~/utils/types';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { UserShort } from '~~/gen/ts/resources/users/users';

type AccessType = { id: number; name: string };

const props = withDefaults(
    defineProps<{
        readOnly?: boolean;
        showRequired?: boolean;
        init: {
            id: string;
            type: number;
            values: {
                userId?: number;
                job?: string;
                minimumGrade?: number;
                accessRole?: AccessLevel;
            };
            required?: boolean;
        };
        accessTypes: AccessType[];
        accessRoles?: AccessLevel[];
        jobs: Job[] | undefined;
    }>(),
    {
        readOnly: false,
        showRequired: false,
        accessRoles: undefined,
    },
);

const emit = defineEmits<{
    (e: 'typeChange', payload: { id: string; type: number }): void;
    (
        e: 'nameChange',
        payload: {
            id: string;
            job: Job | undefined;
            char: UserShort | undefined;
            required?: boolean;
        },
    ): void;
    (e: 'rankChange', payload: { id: string; rank: JobGrade; required?: boolean }): void;
    (e: 'accessChange', payload: { id: string; access: AccessLevel; required?: boolean }): void;
    (e: 'deleteRequest', payload: { id: string }): void;
    (e: 'requiredChange', payload: { id: string; required?: boolean }): void;
}>();

const { t } = useI18n();

const completorStore = useCompletorStore();

const required = ref<boolean | undefined>(props.init.required);
const selectedAccessType = ref<AccessType>({
    id: -1,
    name: '',
});
const usersLoading = ref(false);
const selectedUser = ref<undefined | UserShort>(undefined);
const selectedJob = ref<undefined | Job>();
const selectedMinimumRank = ref<undefined | JobGrade>(undefined);
const selectedAccessRole = ref<ArrayElement<typeof entriesAccessRoles>>();

const entriesAccessRoles: {
    id: AccessLevel;
    label: string;
    value: string;
}[] = [];
if (props.accessRoles === undefined || props.accessRoles.length === 0) {
    entriesAccessRoles.push(
        ...listEnumValues(AccessLevel)
            .map((e, k) => {
                return {
                    id: k,
                    label: t(`enums.docstore.AccessLevel.${e.name}`),
                    value: e.name,
                };
            })
            .filter((e) => e.id !== 0),
    );
} else {
    props.accessRoles.forEach((e) => {
        entriesAccessRoles.push({
            id: e,
            label: t(`enums.docstore.AccessLevel.${AccessLevel[e]}`),
            value: AccessLevel[e],
        });
    });
}

async function findUser(userId?: number): Promise<UserShort[]> {
    if (userId === undefined) return [];

    return completorStore.completeCitizens({
        search: '',
        userId: userId,
    });
}

async function setFromProps(): Promise<void> {
    if (props.init.type === 0 && props.init.values.userId !== undefined) {
        const users = await findUser(props.init.values.userId);
        selectedUser.value = users.find((char) => char.userId === props.init.values.userId);
    } else if (props.init.type === 1 && props.init.values.job !== undefined && props.init.values.minimumGrade !== undefined) {
        selectedJob.value = props.jobs?.find((j) => j.name === props.init.values.job);
        if (selectedJob.value) {
            selectedMinimumRank.value = selectedJob.value.grades.find((rank) => rank.grade === props.init.values.minimumGrade);
        }
    }

    selectedAccessRole.value = entriesAccessRoles.find((type) => type.id === props.init.values.accessRole);

    const passedType = props.accessTypes.find((e) => e.id === props.init.type);
    if (passedType) {
        selectedAccessType.value = passedType;
    }
}

setFromProps();
watch(props, () => setFromProps());

watch(required, () => emit('requiredChange', { id: props.init.id, required: required.value }));

watch(selectedAccessType, async () => {
    emit('typeChange', {
        id: props.init.id,
        type: selectedAccessType.value.id,
    });

    selectedUser.value = undefined;
    selectedJob.value = undefined;
    selectedMinimumRank.value = undefined;
});

watch(selectedJob, () => {
    if (!selectedJob.value) {
        return;
    }

    emit('nameChange', {
        id: props.init.id,
        job: selectedJob.value,
        char: undefined,
        required: required.value,
    });
});

watch(selectedUser, () => {
    if (!selectedUser.value) {
        return;
    }

    emit('nameChange', {
        id: props.init.id,
        job: undefined,
        char: selectedUser.value,
        required: required.value,
    });
});

watch(selectedMinimumRank, () => {
    if (!selectedMinimumRank.value) {
        return;
    }

    emit('rankChange', { id: props.init.id, rank: selectedMinimumRank.value });
});

watch(selectedAccessRole, () => {
    if (!selectedAccessRole.value) {
        return;
    }

    emit('accessChange', {
        id: props.init.id,
        access: selectedAccessRole.value.id,
    });
});
</script>

<template>
    <div class="my-2 flex flex-row items-center gap-1">
        <UCheckbox
            v-if="showRequired"
            v-model="required"
            class="flex-initial"
            :disabled="readOnly"
            :title="$t('common.require')"
            name="required"
        />

        <UFormGroup class="w-60 flex-initial">
            <UInput v-if="accessTypes.length === 1" type="text" disabled :model-value="accessTypes[0]?.name" />
            <ClientOnly v-else>
                <USelectMenu
                    v-model="selectedAccessType"
                    :disabled="readOnly"
                    :options="accessTypes"
                    :placeholder="$t('common.type')"
                    :searchable-placeholder="$t('common.search_field')"
                >
                    <template #label>
                        <span v-if="selectedAccessType" class="truncate">{{ selectedAccessType.name }}</span>
                    </template>
                    <template #option="{ option }">
                        <span class="truncate">{{ option.name }}</span>
                    </template>
                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>
                    <template #empty>
                        {{ $t('common.not_found', [$t('common.type')]) }}
                    </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <template v-if="selectedAccessType?.id === 1">
            <UFormGroup name="selectedJob" class="flex-1">
                <ClientOnly>
                    <USelectMenu
                        v-model="selectedJob"
                        :disabled="readOnly"
                        class="flex-1"
                        option-attribute="label"
                        searchable
                        :search-attributes="['name', 'label']"
                        :options="jobs ?? []"
                        :placeholder="$t('common.job')"
                        :searchable-placeholder="$t('common.search_field')"
                    >
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>
                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>

            <UFormGroup name="selectedMinimumRank" class="flex-1">
                <ClientOnly>
                    <USelectMenu
                        v-model="selectedMinimumRank"
                        :disabled="readOnly || !selectedJob"
                        class="flex-1"
                        option-attribute="label"
                        searchable
                        :search-attributes="['name', 'label']"
                        :options="selectedJob?.grades"
                        :placeholder="$t('common.rank')"
                        :searchable-placeholder="$t('common.search_field')"
                    >
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>
                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>
        </template>

        <template v-else>
            <UFormGroup name="selectedUser" class="flex-1">
                <ClientOnly>
                    <USelectMenu
                        v-model="selectedUser"
                        :searchable="
                            async (query: string) => {
                                usersLoading = true;
                                const users = await completorStore.completeCitizens({
                                    search: query,
                                });
                                usersLoading = false;
                                return users;
                            }
                        "
                        searchable-lazy
                        :searchable-placeholder="$t('common.search_field')"
                        :search-attributes="['firstname', 'lastname']"
                        class="flex-1"
                        :placeholder="$t('common.citizen', 1)"
                        trailing
                        by="userId"
                    >
                        <template #label>
                            <template v-if="selectedUser">
                                {{ usersToLabel([selectedUser]) }}
                            </template>
                        </template>
                        <template #option="{ option: user }">
                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                        </template>
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>
                        <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>
        </template>

        <UFormGroup class="w-60 flex-initial">
            <ClientOnly>
                <USelectMenu
                    v-model="selectedAccessRole"
                    :disabled="readOnly"
                    class="flex-1"
                    option-attribute="label"
                    searchable
                    :search-attributes="['label']"
                    :options="entriesAccessRoles"
                    :placeholder="$t('common.na')"
                    :searchable-placeholder="$t('common.search_field')"
                >
                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>
                    <template #empty> {{ $t('common.not_found', [$t('common.access', 2)]) }} </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <UButton
            v-if="!readOnly"
            :ui="{ rounded: 'rounded-full' }"
            class="flex-initial"
            icon="i-mdi-close"
            @click="$emit('deleteRequest', { id: props.init.id })"
        />
    </div>
</template>
