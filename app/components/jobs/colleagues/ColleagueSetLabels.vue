<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import type { JobsUserProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Labels } from '~~/gen/ts/resources/jobs/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    modelValue?: Labels;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', labels: Labels | undefined): void;
}>();

const { attr, can } = useAuth();

const labels = useVModel(props, 'modelValue', emit);

const notifications = useNotificatorStore();

const canDo = computed(() => ({
    set: can('JobsService.SetJobsUserProps').value && attr('JobsService.SetJobsUserProps', 'Types', 'Labels').value,
}));

async function getColleagueLabels(): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await getGRPCJobsClient().getColleagueLabels({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const changed = ref(false);

const schema = z.object({
    labels: z
        .object({
            id: z.string(),
            name: z.string().min(1),
            color: z.string().length(7),
        })
        .array()
        .max(10),
    reason: z.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: labels.value?.list !== undefined ? labels.value.list.slice() : [],
    reason: '',
});

async function setUserJobProp(userId: number, values: Schema): Promise<void> {
    console.log('setUserJobProp labels', values.labels);

    const jobsUserProps: JobsUserProps = {
        userId: userId,
        job: '',
        labels: {
            list: values.labels,
        },
    };

    try {
        const call = getGRPCJobsClient().setJobsUserProps({
            props: jobsUserProps,
            reason: values.reason,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
    await setUserJobProp(props.userId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    changed.value = false;
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

onMounted(() => {
    state.labels = labels.value?.list.map((l) => ({ ...l, selected: true })) ?? [];
});
</script>

<template>
    <UForm :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
        <p v-if="!state.labels.length" class="text-sm leading-6">
            {{ $t('common.none', [$t('common.label', 2)]) }}
        </p>
        <template v-else>
            <div class="flex max-w-72 flex-row flex-wrap gap-1">
                <UBadge
                    v-for="(attribute, idx) in state.labels"
                    :key="attribute.name"
                    :style="{ backgroundColor: attribute.color }"
                    class="justify-between gap-2"
                    :class="isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                    size="lg"
                >
                    <span class="truncate">
                        {{ attribute.name }}
                    </span>

                    <UButton
                        v-if="canDo.set"
                        variant="link"
                        icon="i-mdi-close"
                        :padded="false"
                        :ui="{ rounded: 'rounded-full' }"
                        :class="
                            isColourBright(hexToRgb(attribute.color, RGBBlack)!)
                                ? '!bg-white/20 !text-black'
                                : '!bg-black/20 !text-white'
                        "
                        @click="
                            changed = true;
                            state.labels.splice(idx, 1);
                        "
                    />
                </UBadge>
            </div>
        </template>

        <UFormGroup name="labels">
            <ClientOnly>
                <USelectMenu
                    v-model="state.labels"
                    multiple
                    :searchable="
                        async (_: string) => {
                            return (await getColleagueLabels()).labels;
                        }
                    "
                    searchable-lazy
                    :searchable-placeholder="$t('common.search_field')"
                    :search-attributes="['name']"
                    option-attribute="name"
                    by="name"
                    clear-search-on-close
                >
                    <template #label>
                        <span v-if="state.labels.length" class="truncate">{{
                            $t('common.selected_no', [state.labels.length])
                        }}</span>
                        <span v-else>&nbsp;</span>
                    </template>

                    <template #option="{ option }">
                        <span class="truncate" :style="{ backgroundColor: option.color }">{{ option.name }}</span>
                    </template>

                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>

                    <template #empty>
                        {{ $t('common.not_found', [$t('common.label', 2)]) }}
                    </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <template v-if="changed">
            <UFormGroup name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" />
            </UFormGroup>

            <UButton type="submit" block icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
