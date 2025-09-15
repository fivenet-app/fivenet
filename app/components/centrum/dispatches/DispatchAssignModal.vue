<script lang="ts" setup>
import { z } from 'zod';
import { type GroupedUnits, statusOrder, unitStatusToBGColor } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    dispatchId: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches, getSortedUnits } = storeToRefs(centrumStore);

const dispatch = computed(() => dispatches.value.get(props.dispatchId));

const centrumCentrumClient = await getCentrumCentrumClient();

const schema = z.object({
    units: z.custom<number>().array().max(10).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    units: [],
});

async function assignDispatch(): Promise<void> {
    if (dispatch.value === undefined) return;

    try {
        const toAdd: number[] = [];
        const toRemove: number[] = [];
        state.units.forEach((u) => {
            toAdd.push(u);
        });
        dispatch.value?.units?.forEach((u) => {
            const idx = state.units.findIndex((su) => su === u.unitId);
            if (idx === -1) {
                toRemove.push(u.unitId);
            }
        });

        const call = centrumCentrumClient.assignDispatch({
            dispatchId: props.dispatchId,
            toAdd: toAdd,
            toRemove: toRemove,
        });
        await call;

        state.units.length = 0;

        emit('close', false);
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
        state.units.splice(idx, 1);
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

function updateDispatchUnits(): void {
    state.units = [...(dispatch.value?.units?.map((du) => du.unitId) ?? [])];
}

watch(dispatch, () => updateDispatchUnits);
updateDispatchUnits();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignDispatch().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.centrum.assign_dispatch.title')">
        <template #actions>
            <IDCopyBadge :id="dispatch?.id ?? dispatchId" class="ml-2" prefix="DSP" />
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div class="flex flex-1 flex-col justify-between gap-1 px-2">
                    <template v-for="group in grouped" :key="group.key">
                        <h3 class="text-sm">
                            {{ $t(`enums.centrum.StatusUnit.${StatusUnit[group.status]}`) }}
                        </h3>

                        <div class="grid grid-cols-2 gap-2 lg:grid-cols-3">
                            <UButton
                                v-for="unit in group.units"
                                :key="unit.name"
                                class="inline-flex flex-row items-center gap-x-1 rounded-md p-1.5 text-sm font-medium hover:transition-all"
                                :class="[
                                    unitStatusToBGColor(unit.status?.status),
                                    unit.users.length === 0 ? '!bg-error-600' : '',
                                ]"
                                :disabled="unit.users.length === 0"
                                @click="selectUnit(unit)"
                            >
                                <UIcon v-if="state.units.includes(unit.id)" class="size-5" name="i-mdi-check" />
                                <UIcon v-else-if="unit.users.length > 0" class="size-5" name="i-mdi-checkbox-blank-outline" />
                                <UIcon v-else class="size-5" name="i-mdi-cancel" />

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
