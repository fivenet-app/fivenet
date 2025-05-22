<script lang="ts" setup>
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    unit: Unit;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const completorStore = useCompletorStore();

const usersLoading = ref(false);

const schema = z.object({
    users: z.custom<UserShort>().array().max(10),
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

        const call = $grpc.centrum.centrum.assignUnit({
            unitId: unitId,
            toAdd: toAdd,
            toRemove: toRemove,
        });
        await call;

        isOpen.value = false;
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
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        padding: 'px-1 py-2 sm:p-2',
                    },
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.centrum.assign_unit.title') }}: {{ unit.name }} ({{ unit.initials }})
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div class="flex flex-1 flex-col justify-between gap-2">
                        <div class="divide-y divide-gray-100 px-2 sm:px-6 dark:divide-gray-800">
                            <UFormGroup class="flex-1" name="users" :label="$t('common.colleague', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.users"
                                        multiple
                                        :searchable="
                                            async (query: string) => {
                                                usersLoading = true;
                                                const colleagues = await completorStore.completeCitizens({
                                                    search: query,
                                                });
                                                usersLoading = false;
                                                return colleagues;
                                            }
                                        "
                                        searchable-lazy
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['firstname', 'lastname']"
                                        block
                                        :placeholder="$t('common.search')"
                                        trailing
                                        by="userId"
                                        :disabled="!canSubmit"
                                    >
                                        <template #option="{ option: user }">
                                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.colleague', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>

                            <div class="mt-2 overflow-hidden rounded-md bg-neutral-100 dark:bg-base-900">
                                <ul
                                    class="divide-y divide-gray-100 text-sm font-medium text-gray-100 dark:divide-gray-800"
                                    role="list"
                                >
                                    <li
                                        v-for="user in state.users"
                                        :key="user.userId"
                                        class="inline-flex items-center px-4 py-2"
                                    >
                                        <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
