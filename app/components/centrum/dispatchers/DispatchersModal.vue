<script lang="ts" setup>
import ColleagueName from '~/components/jobs/colleagues/ColleagueName.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCentrumStore } from '~/stores/centrum';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';

const { isOpen } = useModal();

const { can } = useAuth();

const centrumStore = useCentrumStore();
const { updateDispatchers } = centrumStore;
const { dispatchers, getCurrentMode } = storeToRefs(centrumStore);
</script>

<template>
    <UModal :ui="{ width: '!max-w-2xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.dispatchers', 2) }}
                        <UBadge color="gray">
                            {{ $t('common.mode') }}: {{ $t(`enums.centrum.CentrumMode.${CentrumMode[getCurrentMode ?? 0]}`) }}
                        </UBadge>
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <DataNoDataBlock
                v-if="dispatchers && dispatchers.length === 0"
                class="mt-5"
                icon="i-mdi-monitor"
                :type="$t('common.dispatchers', 2)"
            />
            <UPageGrid v-else>
                <UPageCard
                    v-for="dispatcher in dispatchers"
                    :key="dispatcher.userId"
                    :ui="{
                        title: 'text-gray-900 dark:text-white text-base font-semibold flex items-center gap-1.5 line-clamp-2 whitespace-break-spaces',
                    }"
                >
                    <template #title>
                        <div class="flex flex-1 items-center gap-2">
                            <ProfilePictureImg
                                :src="dispatcher.avatar"
                                :name="`${dispatcher.firstname} ${dispatcher.lastname}`"
                            />
                            <ColleagueName :colleague="dispatcher" />
                        </div>
                    </template>

                    <div class="flex items-center justify-between gap-2">
                        <PhoneNumberBlock :number="dispatcher.phoneNumber" />

                        <UTooltip v-if="can('centrum.CentrumService.UpdateDispatchers').value" :text="$t('common.remove')">
                            <UButton
                                variant="ghost"
                                icon="i-mdi-trash"
                                color="error"
                                @click="updateDispatchers([dispatcher.userId])"
                            />
                        </UTooltip>
                    </div>
                </UPageCard>
            </UPageGrid>

            <template #footer>
                <UButton class="flex-1" color="black" block @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
