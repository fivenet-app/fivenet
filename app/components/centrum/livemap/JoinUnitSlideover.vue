<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { UnitAccessLevel } from '~~/gen/ts/resources/centrum/access';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import { checkUnitAccess } from '../helpers';

const emit = defineEmits<{
    (e: 'joined', unit: Unit): void;
    (e: 'left'): void;
}>();

const { isOpen } = useSlideover();

const centrumStore = useCentrumStore();
const { ownUnitId, getSortedUnits } = storeToRefs(centrumStore);

async function joinOrLeaveUnit(unitId?: string): Promise<void> {
    try {
        const call = getGRPCCentrumClient().joinUnit({
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
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (unitID?: string) => {
    canSubmit.value = false;
    await joinOrLeaveUnit(unitID).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const queryUnit = ref('');

const filteredUnits = computed(() => ({
    available: getSortedUnits.value
        .filter(
            (u) =>
                (u.name.toLowerCase().includes(queryUnit.value.toLowerCase()) ||
                    u.initials.toLowerCase().includes(queryUnit.value.toLowerCase())) &&
                checkUnitAccess(u.access, UnitAccessLevel.JOIN),
        )
        .sort((a, b) => a.name.localeCompare(b.name)),
    unavailable: getSortedUnits.value
        .filter(
            (u) =>
                (u.name.toLowerCase().includes(queryUnit.value.toLowerCase()) ||
                    u.initials.toLowerCase().includes(queryUnit.value.toLowerCase())) &&
                !checkUnitAccess(u.access, UnitAccessLevel.JOIN),
        )
        .sort((a, b) => a.name.localeCompare(b.name)),
}));
</script>

<template>
    <USlideover :ui="{ width: 'w-screen max-w-xl' }">
        <UCard
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-2 text-2xl font-semibold leading-6">
                        {{ $t('common.leave_unit') }}

                        <UIcon v-if="!canSubmit" name="i-mdi-loading" class="size-6 animate-spin" />
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <div class="flex flex-col gap-1">
                    <UFormGroup name="search" :label="$t('common.search')">
                        <UInput
                            v-model="queryUnit"
                            type="text"
                            name="search"
                            :placeholder="$t('common.search')"
                            leading-icon="i-mdi-search"
                        />
                    </UFormGroup>

                    <div class="grid grid-cols-2 gap-2">
                        <UButton
                            v-for="unit in filteredUnits.available"
                            :key="unit.name"
                            :color="ownUnitId !== undefined && ownUnitId === unit.id ? 'amber' : 'primary'"
                            :disabled="!canSubmit || !checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
                            class="flex flex-col"
                            @click="onSubmitThrottle(unit.id)"
                        >
                            <span class="text-base">
                                <span class="font-semibold">{{ unit.initials }}:</span>
                                {{ unit.name }}
                            </span>
                            <span class="mt-0.5 text-xs">
                                {{ $t('common.member', unit.users.length) }}
                            </span>
                            <span v-if="unit.description && unit.description.length > 0" class="text-xs">
                                <span class="font-semibold">{{ $t('common.description') }}:</span>
                                <span class="line-clamp-1">{{ unit.description }}</span>
                            </span>
                        </UButton>
                    </div>

                    <div v-if="filteredUnits.unavailable.length > 0">
                        <h3>{{ $t('common.unavailable') }}</h3>

                        <div class="grid grid-cols-2 gap-2">
                            <UButton
                                v-for="unit in filteredUnits.unavailable"
                                :key="unit.name"
                                :color="ownUnitId !== undefined && ownUnitId === unit.id ? 'amber' : 'primary'"
                                :disabled="!canSubmit || !checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
                                class="flex flex-col"
                                @click="onSubmitThrottle(unit.id)"
                            >
                                <span class="text-base">
                                    <span class="font-semibold">{{ unit.initials }}:</span>
                                    {{ unit.name }}
                                </span>
                                <span class="mt-0.5 text-xs">
                                    {{ $t('common.member', unit.users.length) }}
                                </span>
                                <span v-if="unit.description && unit.description.length > 0" class="text-xs">
                                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                                    <span class="line-clamp-1">{{ unit.description }}</span>
                                </span>
                            </UButton>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        v-if="ownUnitId !== undefined"
                        block
                        color="red"
                        class="flex-1"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle()"
                    >
                        {{ $t('common.leave') }}
                    </UButton>
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
