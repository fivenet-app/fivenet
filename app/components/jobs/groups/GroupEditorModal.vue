<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import { type Group, GroupMembershipMode, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { groupMembershipModeItems, groupTypeItems } from './helpers';

const props = defineProps<{
    group?: Group;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'created', group: Group): void;
    (e: 'updated', group: Group): void;
}>();

const notifications = useNotificationsStore();

const jobsGroupsClient = await getJobsGroupsClient();

const schema = z.object({
    name: z.coerce.string().min(3).max(64),
    description: z.coerce.string().max(255).default(''),
    shortName: z.coerce.string().max(12).default(''),
    logoFileId: z.coerce.string().max(255).default(''),
    color: z.coerce.string().max(7).default(''),
    type: z.enum(GroupType),
    membershipMode: z.enum(GroupMembershipMode),
    sortOrder: z.coerce.number().int().min(0).default(0),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    shortName: '',
    logoFileId: '',
    color: '',
    type: GroupType.MANUAL,
    membershipMode: GroupMembershipMode.FLEXIBLE,
    sortOrder: 0,
});

const formSnapshot = computed(() => ({
    name: state.name,
    description: state.description,
    shortName: state.shortName,
    logoFileId: state.logoFileId,
    color: state.color,
    type: state.type,
    membershipMode: state.membershipMode,
    sortOrder: state.sortOrder,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot);

function setFormFromProps(): void {
    state.name = props.group?.name ?? '';
    state.description = props.group?.description ?? '';
    state.shortName = props.group?.shortName ?? '';
    state.logoFileId = props.group?.logoFileId ?? '';
    state.color = props.group?.color ?? '';
    state.type = props.group?.type ?? GroupType.MANUAL;
    state.membershipMode = props.group?.membershipMode ?? GroupMembershipMode.FLEXIBLE;
    state.sortOrder = props.group?.sortOrder ?? 0;

    syncSnapshot();
}

onBeforeMount(setFormFromProps);
watch(
    () => props.group,
    () => setFormFromProps(),
);

const formRef = useTemplateRef('formRef');

const canSubmit = ref<boolean>(true);

async function createOrUpdateGroup(values: Schema): Promise<void> {
    try {
        const payload = {
            id: props.group?.id ?? 0,
            name: values.name.trim(),
            description: values.description.trim() ? values.description.trim() : undefined,
            shortName: values.shortName.trim() ? values.shortName.trim() : undefined,
            logoFileId: values.logoFileId.trim() ? values.logoFileId.trim() : undefined,
            color: values.color.trim() ? values.color.trim() : undefined,
            type: values.type,
            membershipMode: values.membershipMode,
            sortOrder: values.sortOrder,
        };

        const call = props.group?.id
            ? jobsGroupsClient.updateGroup(payload)
            : jobsGroupsClient.createGroup({
                  ...payload,
                  jobId: 0,
                  leaderUserIds: [],
                  manualMemberUserIds: [],
                  rules: [],
              });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (props.group?.id) {
            emit('updated', response.group!);
        } else {
            emit('created', response.group!);
        }

        emit('close', false);
        syncSnapshot();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateGroup(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal :title="group?.id ? 'Update group' : 'Create group'" :close="false" :dismissible="!hasUnsavedChanges && canSubmit">
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ group?.id ? 'Update group' : 'Create group' }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div class="grid gap-4">
                    <UFormField name="name" label="Name" required>
                        <UInput v-model="state.name" class="w-full" name="name" type="text" placeholder="Name" />
                    </UFormField>

                    <UFormField name="description" label="Description">
                        <UTextarea
                            v-model="state.description"
                            class="w-full"
                            name="description"
                            :rows="3"
                            placeholder="Description"
                        />
                    </UFormField>

                    <div class="grid gap-4 sm:grid-cols-2">
                        <UFormField name="shortName" label="Short name">
                            <UInput
                                v-model="state.shortName"
                                class="w-full"
                                name="shortName"
                                type="text"
                                placeholder="Short name"
                            />
                        </UFormField>

                        <UFormField name="sortOrder" label="Sort order">
                            <UInput
                                v-model="state.sortOrder"
                                class="w-full"
                                name="sortOrder"
                                type="number"
                                min="0"
                                placeholder="0"
                            />
                        </UFormField>
                    </div>

                    <UFormField name="logoFileId" label="Logo file ID">
                        <UInput
                            v-model="state.logoFileId"
                            class="w-full"
                            name="logoFileId"
                            type="text"
                            placeholder="Logo file ID"
                        />
                    </UFormField>

                    <UFormField name="color" label="Color">
                        <ColorPicker v-model="state.color" class="w-full" block />
                    </UFormField>

                    <div class="grid gap-4 sm:grid-cols-2">
                        <UFormField name="type" label="Type" required>
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.type"
                                    class="w-full"
                                    :items="groupTypeItems"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #empty> {{ $t('common.not_found', ['types']) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField name="membershipMode" label="Membership mode" required>
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.membershipMode"
                                    class="w-full"
                                    :items="groupMembershipModeItems"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #empty> {{ $t('common.not_found', ['membership modes']) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </div>
                </div>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :loading="!canSubmit"
                    :disabled="!canSubmit"
                    :label="group?.id ? $t('common.update') : $t('common.create')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
