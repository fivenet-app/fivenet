<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import type { DispatchAssignment } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusUnit, type Unit } from '~~/gen/ts/resources/centrum/units';
import { unitStatusToBGColor } from '../helpers';

const centrumStore = useCentrumStore();
const { timeCorrection } = storeToRefs(centrumStore);

const props = withDefaults(
    defineProps<{
        unit: Unit | undefined;
        initialsOnly?: boolean;
        assignment?: DispatchAssignment;
    }>(),
    {
        initialsOnly: false,
        assignment: undefined,
    },
);

const unitStatusColor = computed(() => unitStatusToBGColor(props.unit?.status?.status));
</script>

<template>
    <template v-if="!unit">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>{{ $t('common.na') }}</span>
        </span>
    </template>
    <UPopover v-else>
        <UButton class="inline-flex items-center gap-1 p-0.5" variant="outline" size="xs">
            <slot name="before" />

            <span>
                <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                <template v-else>
                    {{ unit.initials }}
                </template>
            </span>

            <UIcon v-if="assignment?.expiresAt" class="size-4 text-amber-600" name="i-mdi-timer" />
        </UButton>

        <template #content>
            <div class="inline-flex min-w-48 flex-col gap-1 p-4">
                <p class="text-base leading-none font-semibold">{{ unit.name }} ({{ unit.initials }})</p>

                <p v-if="unit.jobLabel" class="text-base leading-none font-semibold">({{ unit.jobLabel }})</p>

                <UBadge class="rounded-sm font-semibold" :class="unitStatusColor" size="xs">
                    {{ $t(`enums.centrum.StatusUnit.${StatusUnit[unit.status?.status ?? 0]}`) }}
                </UBadge>

                <p v-if="assignment?.expiresAt" class="inline-flex items-center gap-1 text-sm font-normal">
                    <UIcon class="size-4 text-amber-600" name="i-mdi-timer" />
                    {{
                        useLocaleTimeAgo(toDate(assignment?.expiresAt, timeCorrection), {
                            showSecond: true,
                            updateInterval: 1_000,
                        }).value
                    }}
                </p>

                <div class="text-highlighted">
                    <p class="text-sm leading-none font-medium">
                        {{ $t('common.members') }}
                    </p>
                    <template v-if="unit.users.length === 0">
                        <p class="text-xs font-normal">
                            {{ $t('common.member', 0) }}
                        </p>
                    </template>
                    <ul v-else class="inline-flex flex-col text-xs font-normal">
                        <li v-for="user in unit.users" :key="user.userId" class="inline-flex items-center gap-1">
                            <span>{{ user.user?.firstname }} {{ user.user?.lastname }}</span>
                            <PhoneNumberBlock :number="user.user?.phoneNumber" :hide-number="true" />
                        </li>
                    </ul>
                </div>
            </div>
        </template>
    </UPopover>
</template>
