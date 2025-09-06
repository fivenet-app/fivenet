<script lang="ts" setup>
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    unit: Unit;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const completorStore = useCompletorStore();

const centrumCentrumClient = await getCentrumCentrumClient();

const schema = z.object({
    users: z.custom<UserShort>().array().max(10).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    users: props.unit.users.filter((u) => u !== undefined).map((u) => u.user!),
});

async function assignUnit(unitId: number): Promise<void> {
    try {
        const toAdd: number[] = [];
        const toRemove: number[] = [];
        state.users?.forEach((u) => {
            toAdd.push(u.userId);
        });
        props.unit.users?.forEach((u) => {
            const idx = state.users.findIndex((su) => su.userId === u.userId);
            if (idx === -1) {
                toRemove.push(u.userId);
            }
        });

        const call = centrumCentrumClient.assignUnit({
            unitId: unitId,
            toAdd: toAdd,
            toRemove: toRemove,
        });
        await call;

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(props, () => (state.users = props.unit.users.filter((u) => u !== undefined).map((u) => u.user!)));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignUnit(props.unit.id).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="`${$t('components.centrum.assign_unit.title')}: ${unit.name} (${unit.initials})`">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div class="flex flex-1 flex-col justify-between gap-2">
                    <div class="divide-y divide-gray-100 px-2 sm:px-6 dark:divide-gray-800">
                        <UFormField class="flex-1" name="users" :label="$t('common.colleague', 2)">
                            <SelectMenu
                                v-model="state.users"
                                multiple
                                :searchable="
                                    async (q: string) =>
                                        await completorStore.completeCitizens({
                                            search: q,
                                            userIds: state.users.map((u) => u.userId),
                                        })
                                "
                                searchable-key="completor-citizens"
                                :search-input="{ placeholder: $t('common.search_field') }"
                                :filter-fields="['firstname', 'lastname']"
                                block
                                :placeholder="$t('common.search')"
                                trailing
                                :disabled="!canSubmit"
                                class="w-full"
                            >
                                <template #item="{ item }">
                                    {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.colleague', 2)]) }} </template>
                            </SelectMenu>
                        </UFormField>

                        <div class="mt-2 overflow-hidden rounded-md bg-neutral-100 dark:bg-neutral-900">
                            <ul
                                class="divide-y divide-gray-100 text-sm font-medium text-gray-100 dark:divide-gray-800"
                                role="list"
                            >
                                <li v-for="user in state.users" :key="user.userId" class="inline-flex items-center px-4 py-2">
                                    <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.update')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
