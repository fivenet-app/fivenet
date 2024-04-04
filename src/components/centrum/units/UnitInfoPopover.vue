<script lang="ts" setup>
import { TimerIcon } from 'mdi-vue3';
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
        textClass?: unknown;
        buttonClass?: unknown;
        badge?: boolean;
        assignment?: DispatchAssignment;
    }>(),
    {
        initialsOnly: false,
        badge: false,
        assignment: undefined,
        textClass: '' as any,
        buttonClass: '' as any,
    },
);
</script>

<template>
    <template v-if="!unit">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <UPopover v-else>
        <UButton variant="link" :padded="false" class="inline-flex items-center" :class="buttonClass">
            <slot name="before" />
            <span :class="textClass">
                <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                <template v-else>
                    {{ unit.initials }}
                </template>
                <template v-if="assignment?.expiresAt">
                    <TimerIcon class="size-4 fill-warn-600" />
                </template>
            </span>
            <slot name="after" />
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
                <p class="text-sm font-medium leading-none">
                    {{ $t('common.members') }}
                </p>
                <template v-if="unit.users.length === 0">
                    <p class="text-xs font-normal">
                        {{ $t('common.unit', 0) }}
                    </p>
                </template>
                <ul v-else class="text-xs font-normal">
                    <li v-for="user in unit.users" :key="user.userId" class="inline-flex items-center gap-1">
                        <span>{{ user.user?.firstname }} {{ user.user?.lastname }}</span>
                        <PhoneNumberBlock :number="user.user?.phoneNumber" :hide-number="true" />
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
