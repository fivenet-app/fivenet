<script lang="ts" setup>
import PinnedPropertiesCards from '~/components/partials/properties/PinnedPropertiesCards.vue';
import PropertyGroups from '~/components/partials/properties/PropertyGroups.vue';
import type { User } from '~~/gen/ts/resources/users/user';
import LabelBadge from '../labels/LabelBadge.vue';

const props = defineProps<{
    user: User;
}>();

const user = computed(() => props.user);
const { attr } = useAuth();
const { profileProperties } = useCitizenProperties(user);
</script>

<template>
    <div data-tour="citizen-profile" class="h-full min-h-0 overflow-y-auto">
        <UContainer class="w-full py-4 sm:py-6">
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-4">
                <PinnedPropertiesCards :properties="user.properties?.properties" :user="user" />

                <div class="col-span-full grid grid-cols-1 gap-4 xl:grid-cols-2">
                    <PropertyGroups :properties="profileProperties" :groups="user.properties?.groups" />

                    <UCard
                        v-if="attr('citizens.CitizensService/ListCitizens', 'Properties', 'user.labels').value"
                        :ui="{ body: 'p-0 sm:p-0' }"
                    >
                        <template #header>
                            <h2 class="font-semibold text-highlighted">{{ $t('common.label', 2) }}</h2>
                        </template>

                        <dl class="divide-y divide-default">
                            <div class="flex items-start gap-4 px-4 py-3 sm:px-5">
                                <dd class="min-w-0 text-sm text-toned">
                                    <p v-if="!user.props?.labels?.list.length">
                                        {{ $t('common.none', [$t('common.label', 2)]) }}
                                    </p>
                                    <div v-else class="flex flex-row flex-wrap gap-1">
                                        <LabelBadge v-for="label in user.props?.labels?.list" :key="label.id" :label="label" />
                                    </div>
                                </dd>
                            </div>
                        </dl>
                    </UCard>

                    <UCard
                        v-if="attr('citizens.CitizensService/ListCitizens', 'Properties', 'user.licenses').value"
                        :ui="{ body: 'p-0 sm:p-0' }"
                    >
                        <template #header>
                            <h2 class="font-semibold text-highlighted">{{ $t('common.license', 2) }}</h2>
                        </template>

                        <div class="px-4 py-3 sm:px-5">
                            <p v-if="user.licenses.length === 0" class="text-sm text-toned">{{ $t('common.no_licenses') }}</p>
                            <ul v-else class="divide-y divide-default rounded-md border border-default" role="list">
                                <li
                                    v-for="license in user.licenses"
                                    :key="license.type"
                                    class="flex items-center gap-2 px-3 py-3 text-sm"
                                >
                                    <UIcon class="size-5 shrink-0" name="i-mdi-license" /><span
                                        class="min-w-0 truncate"
                                        :title="license.type.toUpperCase()"
                                        >{{ license.label }}</span
                                    >
                                </li>
                            </ul>
                        </div>
                    </UCard>
                </div>
            </div>
        </UContainer>
    </div>
</template>
