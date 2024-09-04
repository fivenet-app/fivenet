<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import { useCompletorStore } from '~/store/completor';
import type { ManageCitizenAttributesResponse } from '~~/gen/ts/services/citizenstore/citizenstore';

const { isOpen } = useModal();

const completorStore = useCompletorStore();

const schema = z
    .object({
        id: z.string(),
        name: z.string().min(1).max(24),
        color: z.string().length(7),
    })
    .array()
    .max(15);

type Schema = z.output<typeof schema>;

const state = ref<Schema>([]);

const { data: attributes } = useLazyAsyncData('citizenstore-attributes', () => completorStore.completeCitizensAttributes(''));

async function manageCitizenAttributes(values: Schema): Promise<ManageCitizenAttributesResponse> {
    try {
        const { response } = await getGRPCCitizenStoreClient().manageCitizenAttributes({
            attributes: values,
        });

        state.value = response.attributes;

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await manageCitizenAttributes(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(attributes, () => (state.value = attributes.value ?? []));
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.rector.job_props.citizen_attributes.title') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup
                    v-if="state && can('CitizenStoreService.SetUserProps').value"
                    name="citizenAttributes.list"
                    class="grid items-center gap-2"
                    :ui="{ container: '' }"
                >
                    <div class="flex flex-col gap-1">
                        <div v-for="(_, idx) in state" :key="idx" class="flex items-center gap-1">
                            <UFormGroup :name="`${idx}.name`" class="flex-1">
                                <UInput
                                    v-model="state[idx]!.name"
                                    :name="`${idx}.name`"
                                    type="text"
                                    class="w-full flex-1"
                                    :placeholder="$t('common.attributes', 1)"
                                />
                            </UFormGroup>

                            <UFormGroup :name="`${idx}.color`">
                                <ColorPicker v-model="state[idx]!.color" :name="`${idx}.color`" class="min-w-16" />
                            </UFormGroup>

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                :disabled="!canSubmit"
                                icon="i-mdi-close"
                                @click="state.splice(idx, 1)"
                            />
                        </div>
                    </div>

                    <UButton
                        :ui="{ rounded: 'rounded-full' }"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        :class="state.length ? 'mt-2' : ''"
                        @click="state.push({ id: '0', name: '', color: '#ffffff' })"
                    />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
