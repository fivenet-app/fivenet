<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/stores/notificator';
import type { JobsUserProps } from '~~/gen/ts/resources/jobs/colleagues';
import type { Labels } from '~~/gen/ts/resources/jobs/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse, SetJobsUserPropsResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    modelValue?: Labels;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', labels: Labels | undefined): void;
    (e: 'refresh'): void;
}>();

const labels = useVModel(props, 'modelValue', emit);

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

async function getColleagueLabels(search?: string): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await $grpc.jobs.jobs.getColleagueLabels({
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
            id: z.number(),
            name: z.string().min(1),
            color: z.string().length(7),
            order: z.number().nonnegative().default(0),
        })
        .array()
        .max(10),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    labels: labels.value?.list.map((l) => ({ ...l, selected: true })) ?? [],
});

async function setUserJobProp(userId: number, values: Schema): Promise<SetJobsUserPropsResponse> {
    const jobsUserProps: JobsUserProps = {
        userId: userId,
        job: '',
        labels: {
            list: values.labels,
        },
    };

    try {
        const call = $grpc.jobs.jobs.setJobsUserProps({
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
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
    <UForm :schema="schema" :state="state" class="flex flex-1 flex-col gap-2" @submit="onSubmitThrottle">
        <div>
            <UButton v-if="!editing" icon="i-mdi-pencil" @click="editing = true" />
            <UButton
                v-else
                icon="i-mdi-cancel"
                color="red"
                @click="
                    state.labels = labels?.list.map((l) => ({ ...l, selected: true })) ?? [];
                    editing = false;
                "
            />
        </div>

        <div class="flex max-w-72 flex-row flex-wrap gap-1">
            <p v-if="!state.labels.length" class="text-sm leading-6">
                {{ $t('common.none', [$t('common.label', 2)]) }}
            </p>
            <template v-else>
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
                        v-if="editing"
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
            </template>
        </div>

        <UFormGroup v-if="editing" name="labels">
            <ClientOnly>
                <USelectMenu
                    v-model="state.labels"
                    multiple
                    :searchable="async (q: string) => (await getColleagueLabels(q))?.labels ?? []"
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

        <template v-if="editing">
            <UFormGroup name="reason" :label="$t('common.reason')" required>
                <UInput v-model="state.reason" type="text" :disabled="!changed" />
            </UFormGroup>

            <UButton type="submit" icon="i-mdi-content-save" :disabled="!changed || !canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
