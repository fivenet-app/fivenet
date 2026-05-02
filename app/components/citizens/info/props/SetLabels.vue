<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Labels } from '~~/gen/ts/resources/citizens/labels/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { UserProps } from '~~/gen/ts/resources/users/props/props';

const props = defineProps<{
    userId: number;
}>();

const { attr, can } = useAuth();

const labels = defineModel<Labels | undefined>();

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const canDo = computed(() => ({
    set:
        can('citizens.CitizensService/SetUserProps').value &&
        attr('citizens.CitizensService/SetUserProps', 'Fields', 'Labels').value,
}));

const changed = ref(false);

const schema = z.object({
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.coerce.string().min(1),
            color: z.coerce.string().length(7),
            icon: z.coerce.string().max(255).optional(),
            expiresAt: z.custom<Timestamp>().optional(),
            // TODO expiration settings needed, but only when enabled for the label
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
        state.reason = '';
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobProp(props.userId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    changed.value = false;
}, 1000);

watch(labels, () => (state.labels = labels.value?.list !== undefined ? labels.value?.list.slice() : []));

watch(state, () => {
    if (
        state.labels.length === labels.value?.list.length &&
        state.labels.every((el, idx) => el.name === labels.value?.list[idx]?.name)
    ) {
        changed.value = false;
    } else {
        changed.value = true;
    }
});

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UForm ref="formRef" class="flex flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UFormField v-if="canDo.set && can('completor.CompletorService/CompleteCitizenLabels').value" name="labels">
            <SelectMenu
                v-model="state.labels"
                class="w-full"
                multiple
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
            >
                <template #default>
                    {{ $t('common.selected', state.labels.length) }}
                </template>

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
            <UFieldGroup v-for="(label, idx) in state.labels" :key="label.name">
                <UBadge
                    class="w-full"
                    :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                    :style="{ backgroundColor: label.color }"
                    size="md"
                    :icon="label.icon && label.icon !== '' ? convertComponentIconNameToDynamic(label.icon) : undefined"
                    :label="label.name"
                />

                <UTooltip v-if="canDo.set" :text="$t('common.remove')">
                    <UButton
                        variant="outline"
                        color="neutral"
                        size="sm"
                        icon="i-mdi-remove"
                        @click="
                            changed = true;
                            state.labels.splice(idx, 1);
                        "
                    />
                </UTooltip>
            </UFieldGroup>
        </div>

        <template v-if="changed">
            <UFormField name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" class="w-full" type="text" />
            </UFormField>

            <UButton
                block
                icon="i-mdi-content-save"
                :disabled="!canSubmit"
                :loading="!canSubmit"
                :label="$t('common.save')"
                @click="formRef?.submit()"
            />
        </template>
    </UForm>
</template>
