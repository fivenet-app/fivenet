<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { ColleagueProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Labels } from '~~/gen/ts/resources/jobs/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse, SetColleaguePropsResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    modelValue?: Labels;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', labels: Labels | undefined): void;
    (e: 'refresh'): void;
}>();

const labels = useVModel(props, 'modelValue', emit);

const notifications = useNotificationsStore();

const jobsJobsClient = await getJobsJobsClient();

async function getColleagueLabels(search?: string): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await jobsJobsClient.getColleagueLabels({
            search: search,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const changed = ref(false);

const schema = z.object({
    reason: z.string().min(3).max(255),
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.string().min(1),
            color: z.string().length(7),
            order: z.coerce.number().nonnegative().default(0),
        })
        .array()
        .max(10)
        .default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    labels: labels.value?.list.map((l) => ({ ...l, selected: true })) ?? [],
});

async function setUserJobProp(userId: number, values: Schema): Promise<SetColleaguePropsResponse> {
    const jobsUserProps: ColleagueProps = {
        userId: userId,
        job: '',
        labels: {
            list: values.labels,
        },
    };

    try {
        const call = jobsJobsClient.setColleagueProps({
            props: jobsUserProps,
            reason: values.reason,
        });
        const { response } = await call;

        changed.value = false;
        editing.value = false;
        state.reason = '';
        emit('refresh');

        state.labels = labels.value?.list.map((l) => ({ ...l, selected: true })) ?? [];

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setUserJobProp(props.userId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(props, () => (state.labels = labels.value?.list !== undefined ? labels.value?.list.slice() : []));

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

const editing = ref(false);
</script>

<template>
    <UForm class="flex flex-1 flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <div>
            <UTooltip v-if="!editing" :text="$t('common.edit')">
                <UButton icon="i-mdi-pencil" @click="editing = true" />
            </UTooltip>
            <UTooltip v-else :text="$t('common.cancel')">
                <UButton
                    icon="i-mdi-cancel"
                    color="error"
                    @click="
                        state.labels = labels?.list.map((l) => ({ ...l, selected: true })) ?? [];
                        editing = false;
                    "
                />
            </UTooltip>
        </div>

        <div class="flex max-w-72 flex-row flex-wrap gap-1">
            <p v-if="!state.labels.length" class="text-sm leading-6">
                {{ $t('common.none', [$t('common.label', 2)]) }}
            </p>
            <template v-else>
                <UBadge
                    v-for="(attribute, idx) in state.labels"
                    :key="attribute.name"
                    class="justify-between gap-2"
                    :class="isColorBright(hexToRgb(attribute.color, RGBBlack)!) ? 'text-black!' : 'text-white!'"
                    :style="{ backgroundColor: attribute.color }"
                    size="lg"
                >
                    <span class="truncate">
                        {{ attribute.name }}
                    </span>

                    <UTooltip v-if="editing" :text="$t('common.remove')">
                        <UButton
                            :class="
                                isColorBright(hexToRgb(attribute.color, RGBBlack)!)
                                    ? 'bg-white/20! text-black!'
                                    : 'bg-black/20! text-white!'
                            "
                            variant="link"
                            icon="i-mdi-close"
                            @click="
                                changed = true;
                                state.labels.splice(idx, 1);
                            "
                        />
                    </UTooltip>
                </UBadge>
            </template>
        </div>

        <UFormField v-if="editing" name="labels">
            <ClientOnly>
                <USelectMenu
                    v-model="state.labels"
                    multiple
                    :searchable="async (q: string) => (await getColleagueLabels(q))?.labels ?? []"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :filter-fields="['name']"
                    clear-search-on-close
                >
                    <template #item-label="{ item }">
                        <span v-if="item.length" class="inline-flex flex-wrap gap-1 truncate">
                            <UBadge
                                v-for="label in item"
                                :key="label.id"
                                class="truncate"
                                :class="isColorBright(label.color) ? 'text-black!' : 'text-white!'"
                                :style="{ backgroundColor: label.color }"
                                :label="label.name"
                            />
                        </span>
                        <span v-else>&nbsp;</span>
                    </template>

                    <template #item="{ item }">
                        <UBadge
                            class="truncate"
                            :class="isColorBright(item.color) ? 'text-black!' : 'text-white!'"
                            :style="{ backgroundColor: item.color }"
                        >
                            {{ item.name }}
                        </UBadge>
                    </template>

                    <template #empty>
                        {{ $t('common.not_found', [$t('common.label', 2)]) }}
                    </template>
                </USelectMenu>
            </ClientOnly>
        </UFormField>

        <template v-if="editing">
            <UFormField name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" :disabled="!changed" />
            </UFormField>

            <UButton type="submit" icon="i-mdi-content-save" :disabled="!changed || !canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
