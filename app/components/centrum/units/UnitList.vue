<script lang="ts" setup>
import { statusOrder, type GroupedUnits } from '~/components/centrum/helpers';
import UnitListEntry from '~/components/centrum/units/UnitListEntry.vue';
import { useCentrumStore } from '~/store/centrum';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';

const { can } = useAuth();

const centrumStore = useCentrumStore();
const { getSortedUnits, abort, reconnecting } = storeToRefs(centrumStore);

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
</script>

<template>
    <div class="flex h-full grow flex-col overflow-y-auto px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex items-center text-base font-semibold leading-6 text-gray-100">
                {{ $t('common.unit', 2) }}

                <UTooltip v-if="can('CentrumService.CreateOrUpdateUnit').value" :text="$t('common.unit', 2)">
                    <UButton :to="{ name: 'centrum-units' }" icon="i-mdi-cog" variant="link" />
                </UTooltip>
            </h2>
        </div>
        <div class="@container flex-1">
            <div
                v-if="abort === undefined && !reconnecting"
                class="@md:grid-cols-2 @3xl:grid-cols-3 mt-3 grid grid-cols-1 gap-2"
            >
                <USkeleton v-for="idx in 8" :key="idx" class="h-9 w-full" />
            </div>

            <template v-for="group in grouped" v-else :key="group.key">
                <p class="-mb-1.5 text-sm">
                    {{ $t(`enums.centrum.StatusUnit.${StatusUnit[group.status]}`) }}
                </p>
                <ul role="list" class="@md:grid-cols-2 @3xl:grid-cols-3 mt-3 grid grid-cols-1 gap-2">
                    <UnitListEntry v-for="unit in group.units" :key="unit.id" :unit="unit" />
                </ul>
            </template>
        </div>
    </div>
</template>
