<script lang="ts" setup>
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import { UnitAccessLevel } from '~~/gen/ts/resources/centrum/units_access';
import { checkUnitAccess } from '../helpers';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'joined', unit: Unit): void;
    (e: 'left'): void;
}>();

const centrumStore = useCentrumStore();
const { ownUnitId, getSortedUnits } = storeToRefs(centrumStore);

const centrumCentrumClient = await getCentrumCentrumClient();

async function joinOrLeaveUnit(unitId?: number): Promise<void> {
    try {
        const call = centrumCentrumClient.joinUnit({
            unitId: unitId,
        });
        const { response } = await call;

        if (response.unit) {
            emit('joined', response.unit);
        } else {
            emit('left');
        }

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (unitID?: number) => {
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
    <USlideover :title="$t('common.leave_unit')" :overlay="false">
        <template #actions>
            <UIcon v-if="!canSubmit" class="size-6 animate-spin" name="i-mdi-loading" />
        </template>

        <template #body>
            <div class="flex flex-col gap-1">
                <UFormField name="search" :label="$t('common.search')">
                    <UInput
                        v-model="queryUnit"
                        type="text"
                        name="search"
                        :placeholder="$t('common.search')"
                        leading-icon="i-mdi-search"
                    />
                </UFormField>

                <div class="grid grid-cols-2 gap-2">
                    <UButton
                        v-for="unit in filteredUnits.available"
                        :key="unit.name"
                        class="flex flex-col"
                        :color="ownUnitId !== undefined && ownUnitId === unit.id ? 'warning' : 'primary'"
                        :disabled="!canSubmit || !checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
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
                            class="flex flex-col"
                            :color="ownUnitId !== undefined && ownUnitId === unit.id ? 'warning' : 'primary'"
                            :disabled="!canSubmit || !checkUnitAccess(unit.access, UnitAccessLevel.JOIN)"
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
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    v-if="ownUnitId !== undefined"
                    class="flex-1"
                    block
                    color="error"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click="onSubmitThrottle()"
                >
                    {{ $t('common.leave') }}
                </UButton>
                <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </USlideover>
</template>
