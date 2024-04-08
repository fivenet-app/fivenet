<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCompletorStore } from '~/store/completor';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    unit: Unit;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const usersLoading = ref(false);

const completorStore = useCompletorStore();

const selectedUsers = ref<UserShort[]>(props.unit.users.filter((u) => u !== undefined).map((u) => u.user!));

async function assignUnit(): Promise<void> {
    try {
        const toAdd: number[] = [];
        const toRemove: number[] = [];
        selectedUsers.value?.forEach((u) => {
            toAdd.push(u.userId);
        });
        props.unit.users?.forEach((u) => {
            const idx = selectedUsers.value.findIndex((su) => su.userId === u.userId);
            if (idx === -1) {
                toRemove.push(u.userId);
            }
        });

        const call = $grpc.getCentrumClient().assignUnit({
            unitId: props.unit.id,
            toAdd,
            toRemove,
        });
        await call;

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignUnit().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover :ui="{ width: 'w-screen max-w-xl' }">
        <UForm :schema="undefined" :state="{}" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 min-h-[calc(100vh-(2*var(--header-height)))] max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
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

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div class="flex flex-1 flex-col justify-between gap-2">
                        <div class="divide-y divide-gray-200 px-2 sm:px-6">
                            <UFormGroup name="selectedUsers" :label="$t('common.colleague', 2)" class="flex-1">
                                <USelectMenu
                                    v-model="selectedUsers"
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
                                    :search-attributes="['firstname', 'lastname']"
                                    block
                                    :placeholder="selectedUsers ? charsGetDisplayValue(selectedUsers) : $t('common.owner')"
                                    trailing
                                    by="userId"
                                >
                                    <template #option="{ option: user }">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </template>
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                                </USelectMenu>
                            </UFormGroup>

                            <div class="overflow-hidden rounded-md bg-base-800">
                                <ul role="list" class="divide-y divide-gray-200 text-sm font-medium text-gray-100">
                                    <li
                                        v-for="user in selectedUsers"
                                        :key="user.userId"
                                        class="inline-flex items-center px-6 py-4"
                                    >
                                        <CitizenInfoPopover :user="user" />
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </USlideover>
</template>
