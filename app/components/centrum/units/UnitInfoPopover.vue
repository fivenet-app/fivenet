<script lang="ts" setup>
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { DispatchAssignment } from '~~/gen/ts/resources/centrum/dispatches';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const centrumStore = useCentrumStore();
const { timeCorrection } = storeToRefs(centrumStore);

withDefaults(
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
</script>

<template>
    <template v-if="!unit">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>{{ $t('common.na') }}</span>
        </span>
    </template>
    <UPopover v-else>
        <UButton variant="outline" :padded="false" size="xs" class="inline-flex items-center gap-1 p-0.5">
            <slot name="before" />

            <span>
                <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                <template v-else>
                    {{ unit.initials }}
                </template>
            </span>

            <template v-if="assignment?.expiresAt">
                <UIcon name="i-mdi-timer" class="size-4 text-amber-600" />
            </template>
        </UButton>

        <template #panel>
            <div class="p-4">
                <p class="text-base font-semibold leading-none">{{ unit.name }} ({{ unit.initials }})</p>
                <p v-if="assignment?.expiresAt" class="inline-flex items-center justify-center text-sm font-normal">
                    {{
                        useLocaleTimeAgo(toDate(assignment.expiresAt, timeCorrection), {
                            showSecond: true,
                            updateInterval: 1_000,
                        }).value
                    }}
                </p>
                <div class="text-gray-900 dark:text-white">
                    <p class="text-sm font-medium leading-none">
                        {{ $t('common.members') }}
                    </p>
                    <template v-if="unit.users.length === 0">
                        <p class="text-xs font-normal">
                            {{ $t('common.units', 0) }}
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
