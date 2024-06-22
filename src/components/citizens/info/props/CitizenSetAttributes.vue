<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { useNotificatorStore } from '~/store/notificator';
import type { CitizenAttributes, UserProps } from '~~/gen/ts/resources/users/users';
import { useCompletorStore } from '~/store/completor';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    modelValue?: CitizenAttributes;
    userId: number;
}>();

const emits = defineEmits<{
    (e: 'update:modelValue', attributes: CitizenAttributes | undefined): void;
}>();

const attributes = useVModel(props, 'modelValue', emits);

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();

const canDo = computed(() => ({
    set:
        can('CitizenStoreService.SetUserProps').value && attr('CitizenStoreService.SetUserProps', 'Fields', 'Attributes').value,
}));

const attributesLoading = ref(false);

const changed = ref(false);

const schema = z.object({
    attributes: z
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
    attributes: attributes.value?.list !== undefined ? attributes.value.list.slice() : [],
    reason: '',
});

async function setJobProp(userId: number, values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: userId,
        attributes: {
            list: values.attributes,
        },
    };

    try {
        const call = getGRPCCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        attributes.value = response.props?.attributes;
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

watch(props, () => (state.attributes = attributes.value?.list !== undefined ? attributes.value?.list.slice() : []));

watch(state, () => {
    if (
        state.attributes.length === attributes.value?.list.length &&
        state.attributes.every((el, idx) => el.name === attributes.value?.list[idx]?.name)
    ) {
        changed.value = false;
    } else {
        changed.value = true;
    }
});
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle" class="flex flex-col gap-2">
        <p v-if="!state.attributes.length" class="text-sm leading-6">
            {{ $t('common.none', [$t('common.attributes', 2)]) }}
        </p>
        <template v-else>
            <div class="flex max-w-72 flex-row flex-wrap gap-1">
                <UBadge
                    v-for="(attribute, idx) in state.attributes"
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
                            state.attributes.splice(idx, 1);
                        "
                    />
                </UBadge>
            </div>
        </template>

        <UFormGroup v-if="canDo.set && can('CompletorService.CompleteCitizenAttributes').value" name="attributes">
            <USelectMenu
                v-model="state.attributes"
                multiple
                :searchable="
                    async (query: string) => {
                        attributesLoading = true;
                        const colleagues = await completorStore.completeCitizensAttributes(query);
                        attributesLoading = false;
                        return colleagues;
                    }
                "
                searchable-lazy
                :searchable-placeholder="$t('common.search_field')"
                :search-attributes="['name']"
                option-attribute="name"
                by="name"
                clear-search-on-close
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            >
                <template #option="{ option }">
                    <span class="truncate" :style="{ backgroundColor: option.color }">{{ option.name }}</span>
                </template>

                <template #option-empty="{ query: search }">
                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                </template>
                <template #empty>
                    {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                </template>
            </USelectMenu>
        </UFormGroup>

        <template v-if="changed">
            <UFormGroup name="reason" :label="$t('common.reason')">
                <UInput
                    v-model="state.reason"
                    type="text"
                    name="reason"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </UFormGroup>

            <UButton type="submit" block icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                {{ $t('common.save') }}
            </UButton>
        </template>
    </UForm>
</template>
