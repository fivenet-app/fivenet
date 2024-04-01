<script lang="ts" setup>
import { PencilIcon, TrashCanIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { conductTypesToBGColor, conductTypesToRingColor, conductTypesToTextColor } from '~/components/jobs/conduct/helpers';

defineProps<{
    conduct: ConductEntry;
}>();

defineEmits<{
    (e: 'selected'): void;
    (e: 'delete'): void;
}>();

const openMessage = ref(false);
</script>

<template>
    <tr class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1">
            <GenericTime :value="conduct.createdAt" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.expires_at') }}</dt>
                <dd class="mt-1 truncate text-accent-200">
                    <GenericTime v-if="conduct.expiresAt" class="font-semibold" :value="conduct.expiresAt" />
                    <span v-else>
                        {{ $t('components.jobs.conduct.List.no_expiration') }}
                    </span>
                </dd>
            </dl>
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-base font-medium text-accent-200 lg:table-cell">
            <GenericTime v-if="conduct.expiresAt" class="font-semibold" :value="conduct.expiresAt" />
            <span v-else>
                {{ $t('components.jobs.conduct.List.no_expiration') }}
            </span>
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <div
                class="rounded-md px-2 py-1 text-base font-medium ring-1 ring-inset"
                :class="[
                    conductTypesToBGColor(conduct.type),
                    conductTypesToRingColor(conduct.type),
                    conductTypesToTextColor(conduct.type),
                ]"
            >
                {{ $t(`enums.jobs.ConductType.${ConductType[conduct.type ?? (0 as number)]}`) }}
            </div>
        </td>
        <td class="max-w-24 whitespace-normal p-1 text-left text-accent-200 md:max-w-full">
            <p :class="openMessage ? '' : 'line-clamp-2 hover:line-clamp-6'">
                {{ conduct.message }}
            </p>
            <UButton
                v-if="conduct.message.length > 175"
                class="flex justify-center rounded-md bg-accent-500 p-1 text-sm font-semibold text-neutral transition-colors hover:bg-accent-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent-500"
                @click="openMessage = !openMessage"
            >
                {{ openMessage ? $t('common.read_less') : $t('common.read_more') }}
            </UButton>
        </td>
        <td class="whitespace-nowrap p-1 text-left text-accent-200">
            <CitizenInfoPopover :user="conduct.targetUser" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.creator') }}</dt>
                <dd class="mt-1 truncate text-accent-200">
                    <CitizenInfoPopover :user="conduct.creator" />
                </dd>
            </dl>
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-accent-200 lg:table-cell">
            <CitizenInfoPopover :user="conduct.creator" />
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-base font-medium sm:pr-0">
            <div class="flex flex-col justify-end gap-1 md:flex-row">
                <UButton
                    v-if="can('JobsConductService.UpdateConductEntry')"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="$emit('selected')"
                >
                    <PencilIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                </UButton>

                <UButton
                    v-if="can('JobsConductService.DeleteConductEntry')"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="$emit('delete')"
                >
                    <TrashCanIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                </UButton>
            </div>
        </td>
    </tr>
</template>
