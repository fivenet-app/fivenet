<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const emit = defineEmits<{
    (e: 'joined', unit: Unit): void;
    (e: 'left'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const centrumStore = useCentrumStore();
const { ownUnitId, getSortedUnits } = storeToRefs(centrumStore);

async function joinOrLeaveUnit(unitId?: string): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().joinUnit({
            unitId,
        });
        const { response } = await call;

        if (response.unit) {
            emit('joined', response.unit);
        } else {
            emit('left');
        }

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (unitID?: string) => {
    canSubmit.value = false;
    await joinOrLeaveUnit(unitID).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const queryUnit = ref('');

const filteredUnits = computed(() =>
    getSortedUnits.value
        .filter(
            (u) =>
                u.name.toLowerCase().includes(queryUnit.value.toLowerCase()) ||
                u.initials.toLowerCase().includes(queryUnit.value.toLowerCase()),
        )
        .sort((a, b) => a.name.localeCompare(b.name)),
);
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-col flex-1"
            :ui="{ body: { base: 'flex-1' }, ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.leave_unit') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <div class="divide-y divide-gray-200 px-2">
                    <div class="mt-1">
                        <dl>
                            <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                <dt class="text-sm font-medium leading-6">
                                    <div class="flex h-6 items-center">
                                        {{ $t('common.search') }}
                                    </div>
                                </dt>
                                <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                    <div class="relative flex items-center">
                                        <UInput
                                            v-model="queryUnit"
                                            type="text"
                                            name="search"
                                            :placeholder="$t('common.search')"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>
                                </dd>
                            </div>
                        </dl>
                        <div class="my-2 space-y-24">
                            <div class="flex-1">
                                <div class="grid grid-cols-2 gap-2">
                                    <UButton
                                        v-for="unit in filteredUnits"
                                        :key="unit.name"
                                        :disabled="!canSubmit"
                                        class="group flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:transition-all"
                                        :class="[
                                            !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : ownUnitId !== undefined && ownUnitId === unit.id
                                                  ? 'bg-warn-500 hover:bg-warn-100/10'
                                                  : 'bg-primary-500 hover:bg-primary-100/10',
                                        ]"
                                        @click="onSubmitThrottle(unit.id)"
                                    >
                                        <span class="mt-0.5 text-base">
                                            <span class="font-semibold">{{ unit.initials }}:</span>
                                            {{ unit.name }}
                                        </span>
                                        <span class="mt-1 text-xs">
                                            {{ $t('common.member', unit.users.length) }}
                                        </span>
                                        <span v-if="unit.description && unit.description.length > 0" class="text-xs">
                                            <span class="font-semibold">{{ $t('common.description') }}:</span>
                                            {{ unit.description }}
                                        </span>
                                    </UButton>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButton
                    v-if="ownUnitId !== undefined"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click="onSubmitThrottle()"
                >
                    {{ $t('common.leave') }}
                </UButton>
                <UButton @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
