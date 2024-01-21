<script lang="ts" setup>
import { ArchiveIcon } from 'mdi-vue3';
import { computedAsync } from '@vueuse/core';
import { useCentrumStore } from '~/store/centrum';
import DispatchListEntry from '~/components/centrum/dispatches/DispatchListEntry.vue';
import { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import GenericTime from '~/components/partials/elements/GenericTime.vue';

const props = withDefaults(
    defineProps<{
        dispatches?: Dispatch[];
        showButton?: boolean;
        hideActions?: boolean;
        alwaysShowDay?: boolean;
    }>(),
    {
        dispatches: undefined,
        showButton: false,
        hideActions: false,
        alwaysShowDay: false,
    },
);

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { getSortedDispatches } = storeToRefs(centrumStore);

type GroupedDispatches = { date: Date; key: string; dispatches: Dispatch[] }[];

const grouped = computedAsync(async () => {
    const groups: GroupedDispatches = [];
    (props.dispatches ?? getSortedDispatches.value).forEach((e) => {
        const date = toDate(e.createdAt);
        const idx = groups.findIndex((g) => g.key === dateToDateString(date));
        if (idx === -1) {
            groups.push({
                date,
                dispatches: [e],
                key: dateToDateString(date),
            });
        } else {
            groups[idx].dispatches.push(e);
        }
    });

    return groups;
});
</script>

<template>
    <div class="h-full overflow-y-auto px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
            <div class="inline-flex items-center sm:flex-auto">
                <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6 text-gray-100">
                    {{ $t('common.dispatches') }}

                    <NuxtLink
                        v-if="showButton"
                        :to="{ name: 'centrum-dispatches' }"
                        :title="$t('common.dispatches')"
                        class="ml-2"
                    >
                        <ArchiveIcon class="h-5 w-5" />
                    </NuxtLink>
                </h2>
                <h2 v-if="dispatches === undefined" class="text-base font-semibold text-gray-100">
                    {{ $t('components.centrum.livemap.total_dispatches') }}:
                    {{ getSortedDispatches.length }}
                </h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <table class="min-w-full divide-y divide-base-600">
                        <thead>
                            <tr>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.action', 2) }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.id') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.created_at') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.status') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.postal') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.units') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.citizen') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.attributes', 2) }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.message') }}
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-base-800">
                            <template v-for="(group, idx) in grouped" :key="group.key">
                                <tr v-if="alwaysShowDay || idx !== 0">
                                    <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300" colspan="5">
                                        <GenericTime :value="group.date" type="date" />
                                    </td>
                                </tr>
                                <DispatchListEntry
                                    v-for="dispatch in group.dispatches"
                                    :key="dispatch.id"
                                    :dispatch="dispatch"
                                    :hide-actions="hideActions"
                                    @goto="$emit('goto', $event)"
                                />
                            </template>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
