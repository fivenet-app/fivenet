<script lang="ts" setup>
import { z } from 'zod';
import { statusOrder, unitStatusToBGColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/store/centrum';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';
import type { GroupedUnits } from '~/components/centrum/helpers';

const centrumStore = useCentrumStore();
const { dispatches, getSortedUnits } = storeToRefs(centrumStore);

const props = defineProps<{
    dispatchId: string;
}>();

const dispatch = computed(() => dispatches.value.get(props.dispatchId)!);

const { isOpen } = useModal();

const schema = z.object({
    units: z.custom<string>().array().max(10),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    units: [...dispatch.value.units.map((du) => du.unitId)],
});

async function assignDispatch(): Promise<void> {
    try {
        const toAdd: string[] = [];
        const toRemove: string[] = [];
        state.units.forEach((u) => {
            toAdd.push(u);
        });
        dispatch.value.units?.forEach((u) => {
            const idx = state.units.findIndex((su) => su === u.unitId);
            if (idx === -1) {
                toRemove.push(u.unitId);
            }
        });

        const call = getGRPCCentrumClient().assignDispatch({
            dispatchId: props.dispatchId,
            toAdd,
            toRemove,
        });
        await call;

        state.units.length = 0;

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function selectUnit(item: Unit): void {
    const idx = state.units.findIndex((u) => u === item.id);
    if (idx === -1) {
        state.units.push(item.id);
    } else {
        delete state.units[idx];
    }
}

const grouped = computedAsync(async () => {
    const groups: GroupedUnits = [];
    getSortedUnits.value.forEach((e) => {
        const idx = groups.findIndex((g) => g.key === e.status?.status.toString());
        if (idx === -1) {
            groups.push({
                status: e.status?.status ?? 0,
                units: [e],
                key: e.status?.status.toString() ?? '',
            });
        } else {
            groups[idx]!.units.push(e);
        }
    });

    groups
        .sort((a, b) => statusOrder.indexOf(a.status) - statusOrder.indexOf(b.status))
        .forEach((group) =>
            group.units.sort((a, b) => {
                if (a.users.length === b.users.length) {
                    return 0;
                } else if (a.users.length === 0) {
                    return 1;
                } else if (b.users.length === 0) {
                    return -1;
                } else {
                    return a.name.localeCompare(b.name);
                }
            }),
        );

    return groups;
});

watch(dispatch.value, () => (state.units = [...dispatch.value.units.map((du) => du.unitId)]));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignDispatch().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl', margin: 'sm:my-2' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 h-full max-h-[calc(100vh-(3*var(--header-height)))] overflow-y-auto',
                        padding: 'px-1 py-2 sm:p-2',
                    },
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex items-center text-2xl font-semibold leading-6">
                            {{ $t('components.centrum.assign_dispatch.title') }}:
                            <IDCopyBadge :id="dispatch.id" class="ml-2" prefix="DSP" />
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div class="flex flex-1 flex-col justify-between gap-1 px-2">
                    <template v-for="group in grouped" :key="group.key">
                        <h3 class="text-sm">
                            {{ $t(`enums.centrum.StatusUnit.${StatusUnit[group.status]}`) }}
                        </h3>

                        <div class="grid grid-cols-2 gap-2 lg:grid-cols-3">
                            <UButton
                                v-for="unit in group.units"
                                :key="unit.name"
                                :disabled="unit.users.length === 0"
                                class="hover:bg-primary-100/10 inline-flex flex-row items-center gap-x-1 rounded-md p-1.5 text-sm font-medium hover:transition-all"
                                :class="[
                                    unitStatusToBGColor(unit.status?.status),
                                    unit.users.length === 0 ? '!bg-error-600' : '',
                                ]"
                                @click="selectUnit(unit)"
                            >
                                <UIcon name="i-mdi-check" v-if="state.units.includes(unit.id)" class="size-5" />
                                <UIcon name="i-mdi-checkbox-blank-outline" v-else-if="unit.users.length > 0" class="size-5" />
                                <UIcon name="i-mdi-cancel" v-else class="size-5" />

                                <div class="ml-0.5 flex w-full flex-col place-items-start">
                                    <span class="font-bold">
                                        {{ unit.initials }}
                                    </span>
                                    <span class="text-xs">
                                        {{ unit.name }}
                                    </span>
                                    <span class="mt-1 text-xs">
                                        <span class="block">
                                            {{ $t('common.member', unit.users.length) }}
                                        </span>
                                    </span>
                                </div>
                            </UButton>
                        </div>
                    </template>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
