<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Label, Labels } from '~~/gen/ts/resources/citizens/labels/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';
import ConfigureLabelModal from '../../labels/ConfigureLabelModal.vue';
import { AccessLevel } from '~~/gen/ts/resources/citizens/labels/access';
import type { Access } from '~~/gen/ts/resources/access/access';

const props = defineProps<{
    userId: number;
}>();

const { t } = useI18n();
const { attr, can } = useAuth();

const labels = defineModel<Labels | undefined>({ default: undefined });

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();

const overlay = useOverlay();
const removeConfirmModal = overlay.create(ConfirmModal);

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.coerce.string().min(1),
            color: z.coerce.string().length(7),
            icon: z.coerce.string().max(255).optional(),
            expiresAt: z.custom<Timestamp>().optional(),
            sortOrder: z.number().default(0),
            access: z
                .custom<Access>()
                .default({
                    jobs: [],
                    users: [],
                    qualifications: [],
                })
                .optional(),
        })
        .array()
        .max(10)
        .default([]),
    reason: z.coerce.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: labels.value?.list !== undefined ? labels.value.list.slice() : [],
    reason: '',
});

const originalLabelIds = ref<number[]>(state.labels.map((label) => label.id));

const { hasUnsavedChanges, syncSnapshot } = useSnapshotChanges(state);

function setFromProps(): void {
    state.labels = labels.value?.list !== undefined ? labels.value.list.slice() : [];
    state.reason = '';
    originalLabelIds.value = state.labels.map((label) => label.id);
    syncSnapshot();
}

async function setJobProp(userId: number, values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: userId,
        labels: {
            list: values.labels,
        },
    };

    try {
        const call = citizensCitizensClient.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        labels.value = response.props?.labels;
        setFromProps();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await setJobProp(props.userId, event.data);
});

watch(labels, () => setFromProps());

const selectedLabel = ref<Label | null>(null);

const configureLabelModal = overlay.create(ConfigureLabelModal);

function handleLabelUpdate(label: Label | null): void {
    if (!label) return;
    selectedLabel.value = null;

    configureLabelModal.open({
        label: label,
        onClose: ($event) => {
            if (!$event) return;

            const idx = state.labels.findIndex((l) => l.id === $event.id);
            if (idx == -1) {
                state.labels.unshift({
                    ...label,
                    ...$event,
                });
            } else {
                state.labels[idx] = {
                    ...label,
                    ...$event,
                };
            }
        },
    });
}

function confirmRemoveLabel(label: Label, index: number): void {
    removeConfirmModal.open({
        title: t('common.remove'),
        description: t('components.citizens.CitizenInfoProfile.remove_label_description', { name: label.name }),
        confirm: () => state.labels.splice(index, 1),
    });
}

const formRef = useTemplateRef('formRef');

const labelChanges = computed(() => {
    const original = new Set(originalLabelIds.value);
    const current = new Set(state.labels.map((label) => label.id));

    return {
        added: state.labels.filter((label) => !original.has(label.id)),
        removed: originalLabelIds.value.filter((id) => !current.has(id)),
    };
});
</script>

<template>
    <UForm ref="formRef" class="flex flex-col gap-2" :schema="schema" :state="state" @submit="submit">
        <UFormField
            v-if="
                can('citizens.CitizensService/SetUserProps').value &&
                attr('citizens.CitizensService/SetUserProps', 'Fields', 'Labels').value
            "
            name="labels"
        >
            <SelectMenu
                v-model="selectedLabel"
                class="w-full"
                :searchable="
                    async (q: string) =>
                        (await completorStore.completeCitizenLabels(q)).filter(
                            (l) => !state.labels.some((sl) => sl.id === l.id),
                        )
                "
                searchable-key="completor-citizens-labels"
                :search-input="{ placeholder: $t('common.search_field') }"
                :search-labels="['name']"
                :ui="{ itemLeadingIcon: 'hidden' }"
                @update:model-value="($event) => handleLabelUpdate($event)"
            >
                <template #item-label="{ item }">
                    <UBadge
                        class="truncate"
                        :class="isColorBright(hexToRgb(item.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                        :style="{ backgroundColor: item.color }"
                        :icon="item.icon && item.icon !== '' ? convertComponentIconNameToDynamic(item.icon) : undefined"
                        :label="item.name"
                    />
                </template>

                <template #empty>
                    {{ $t('common.not_found', [$t('common.label', 2)]) }}
                </template>
            </SelectMenu>
        </UFormField>

        <p v-if="!state.labels.length" class="text-sm leading-6">
            {{ $t('common.none', [$t('common.label', 2)]) }}
        </p>
        <div v-else class="flex flex-1 flex-col gap-1">
            <UFieldGroup v-for="label in state.labels" :key="label.id">
                <UBadge
                    class="w-full"
                    :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                    :style="{ backgroundColor: label.color }"
                    size="md"
                    :icon="label.icon && label.icon !== '' ? convertComponentIconNameToDynamic(label.icon) : undefined"
                >
                    <div class="inline-flex flex-col gap-1">
                        <span>{{ label.name }}</span>

                        <div v-if="label.expiresAt">
                            ({{ $t('common.expires_at') }} {{ $d(toDate(label.expiresAt), 'short') }})
                        </div>
                    </div>
                </UBadge>

                <UTooltip v-if="label.access?.jobs.find((ja) => ja.access >= AccessLevel.GIVE)" :text="$t('common.edit')">
                    <UButton
                        variant="outline"
                        color="neutral"
                        size="sm"
                        icon="i-mdi-pencil"
                        @click="() => handleLabelUpdate(label)"
                    />
                </UTooltip>

                <UTooltip v-if="label.access?.jobs.find((ja) => ja.access >= AccessLevel.REMOVE)" :text="$t('common.remove')">
                    <UButton
                        variant="outline"
                        color="neutral"
                        size="sm"
                        icon="i-mdi-remove"
                        @click="
                            confirmRemoveLabel(
                                label,
                                state.labels.findIndex((item) => item.id === label.id),
                            )
                        "
                    />
                </UTooltip>
            </UFieldGroup>
        </div>

        <template v-if="hasUnsavedChanges">
            <UAlert
                color="warning"
                icon="i-mdi-pencil-outline"
                :description="
                    $t('components.citizens.CitizenInfoProfile.label_changes_pending', {
                        added: labelChanges.added.length,
                        removed: labelChanges.removed.length,
                    })
                "
            />

            <UFormField name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" class="w-full" type="text" />
            </UFormField>

            <UButton
                block
                icon="i-mdi-content-save"
                :disabled="!canSubmit"
                :loading="isSubmitting"
                :label="$t('common.save')"
                @click="formRef?.submit()"
            />
        </template>
    </UForm>
</template>
