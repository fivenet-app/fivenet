<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useCompletorStore } from '~/stores/completor';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Labels } from '~~/gen/ts/resources/users/labels';
import type { UserProps } from '~~/gen/ts/resources/users/props';

const props = defineProps<{
    modelValue?: Labels;
    userId: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', labels: Labels | undefined): void;
}>();

const { $grpc } = useNuxtApp();

const { attr, can } = useAuth();

const labels = useVModel(props, 'modelValue', emit);

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();

const canDo = computed(() => ({
    set:
        can('citizens.CitizensService.SetUserProps').value &&
        attr('citizens.CitizensService.SetUserProps', 'Fields', 'Labels').value,
}));

const labelsLoading = ref(false);

const changed = ref(false);

const schema = z.object({
    labels: z
        .object({
            id: z.number(),
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

async function setJobProp(userId: number, values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: userId,
        labels: {
            list: values.labels,
        },
    };

    try {
        const call = $grpc.citizens.citizens.setUserProps({
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
</script>

<template>
    <UForm class="flex flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <p v-if="!state.labels.length" class="text-sm leading-6">
            {{ $t('common.none', [$t('common.label', 2)]) }}
        </p>
        <template v-else>
            <div class="flex max-w-72 flex-row flex-wrap gap-1">
                <UBadge
                    v-for="(attribute, idx) in state.labels"
                    :key="attribute.name"
                    class="justify-between gap-2"
                    :class="isColorBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                    :style="{ backgroundColor: attribute.color }"
                    size="lg"
                >
                    <span class="truncate">
                        {{ attribute.name }}
                    </span>

                    <UButton
                        v-if="canDo.set"
                        :class="
                            isColorBright(hexToRgb(attribute.color, RGBBlack)!)
                                ? '!bg-white/20 !text-black'
                                : '!bg-black/20 !text-white'
                        "
                        variant="link"
                        icon="i-mdi-close"
                        :padded="false"
                        :ui="{ rounded: 'rounded-full' }"
                        @click="
                            changed = true;
                            state.labels.splice(idx, 1);
                        "
                    />
                </UBadge>
            </div>
        </template>

        <UFormGroup v-if="canDo.set && can('completor.CompletorService.CompleteCitizenLabels').value" name="labels">
            <ClientOnly>
                <USelectMenu
                    v-model="state.labels"
                    multiple
                    :searchable="
                        async (query: string) => {
                            labelsLoading = true;
                            const colleagues = await completorStore.completeCitizenLabels(query);
                            labelsLoading = false;
                            return colleagues;
                        }
                    "
                    searchable-lazy
                    :searchable-placeholder="$t('common.search_field')"
                    :search-labels="['name']"
                    option-attribute="name"
                    by="name"
                    clear-search-on-close
                >
                    <template #label>
                        {{ $t('common.selected', state.labels.length) }}
                    </template>

                    <template #option="{ option }">
                        <span
                            class="truncate"
                            :class="isColorBright(hexToRgb(option.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                            :style="{ backgroundColor: option.color }"
                            >{{ option.name }}</span
                        >
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
