<script lang="ts" setup>
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import type { User } from '~~/gen/ts/resources/users/user';
import LabelBadge from '../labels/LabelBadge.vue';

const props = defineProps<{ user: User }>();
const { attr } = useAuth();
const numberFormatter = useDisplayNumberFormat();
const user = computed(() => props.user);
</script>

<template>
    <div class="h-full min-h-0 overflow-y-auto">
        <UContainer class="w-full py-4 sm:py-6">
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-4">
                <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
                    <UCard :class="user.props?.wanted ? 'border-error-500/50' : ''" :ui="{ body: 'p-3 sm:p-4' }">
                        <div class="flex items-center gap-3">
                            <UIcon
                                class="size-5 shrink-0"
                                :class="user.props?.wanted ? 'text-error-500' : 'text-success-500'"
                                :name="user.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-check'"
                            />
                            <div class="min-w-0">
                                <p class="text-xs font-medium text-muted">{{ $t('common.wanted') }}</p>
                                <p class="truncate font-semibold text-highlighted">
                                    {{ user.props?.wanted ? $t('common.yes') : $t('common.no') }}
                                </p>
                            </div>
                        </div>
                    </UCard>
                    <UCard
                        v-if="
                            attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value
                        "
                        :ui="{ body: 'p-3 sm:p-4' }"
                    >
                        <div class="flex items-center gap-3">
                            <UIcon class="size-5 shrink-0 text-muted" name="i-mdi-counter" />
                            <div class="min-w-0">
                                <p class="text-xs font-medium text-muted">{{ $t('common.traffic_infraction_points', 2) }}</p>
                                <p
                                    class="truncate font-semibold text-highlighted"
                                    :class="(user.props?.trafficInfractionPoints ?? 0) >= 10 ? 'text-error-500' : ''"
                                >
                                    {{ $t('common.point', user.props?.trafficInfractionPoints ?? 0) }}
                                </p>
                            </div>
                        </div>
                    </UCard>
                    <UCard
                        v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                        :ui="{ body: 'p-3 sm:p-4' }"
                    >
                        <div class="flex items-center gap-3">
                            <UIcon class="size-5 shrink-0 text-muted" name="i-mdi-cash-remove" />
                            <div class="min-w-0">
                                <p class="text-xs font-medium text-muted">{{ $t('common.fine', 2) }}</p>
                                <p
                                    class="truncate font-semibold text-highlighted"
                                    :class="(user.props?.openFines ?? 0) > 0 ? 'text-error-500' : ''"
                                >
                                    {{ numberFormatter.format(user.props?.openFines ?? 0) }}
                                </p>
                            </div>
                        </div>
                    </UCard>
                    <UCard :ui="{ body: 'p-3 sm:p-4' }">
                        <div class="flex items-center gap-3">
                            <UIcon class="size-5 shrink-0 text-muted" name="i-mdi-briefcase" />
                            <div class="min-w-0">
                                <p class="text-xs font-medium text-muted">{{ $t('common.job') }}</p>
                                <p class="truncate font-semibold text-highlighted">{{ user.jobLabel }}</p>
                            </div>
                        </div>
                    </UCard>
                </div>

                <div class="grid grid-cols-1 gap-4 xl:grid-cols-2">
                    <UCard :ui="{ body: 'p-0 sm:p-0' }">
                        <template #header
                            ><h2 class="font-semibold text-highlighted">{{ $t('common.personal') }}</h2></template
                        >
                        <dl class="divide-y divide-default">
                            <div class="flex items-start gap-4 px-4 py-3 sm:px-5">
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.date_of_birth') }}</dt>
                                <dd class="min-w-0 text-sm text-toned">{{ user.dateofbirth || $t('common.unknown') }}</dd>
                            </div>
                            <div class="flex items-center gap-4 px-4 py-3 sm:px-5">
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.sex') }}</dt>
                                <dd class="inline-flex min-w-0 items-center gap-2 text-sm text-toned">
                                    <span>{{ $t(`common.sex_mapping.${user.sex?.toLowerCase() ?? 'n'}`) }}</span
                                    ><CharSexBadge :sex="user.sex ?? ''" />
                                </dd>
                            </div>
                            <div class="flex items-center gap-4 px-4 py-3 sm:px-5">
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.height') }}</dt>
                                <dd class="min-w-0 text-sm text-toned">
                                    {{ user.height ? `${user.height}cm` : $t('common.unknown') }}
                                </dd>
                            </div>
                            <div
                                v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.BloodType').value"
                                class="flex items-center gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.blood_type') }}</dt>
                                <dd class="min-w-0 text-sm text-toned">{{ user.props?.bloodType ?? $t('common.na') }}</dd>
                            </div>
                        </dl>
                    </UCard>

                    <UCard :ui="{ body: 'p-0 sm:p-0' }">
                        <template #header
                            ><h2 class="font-semibold text-highlighted">
                                {{ $t('common.phone') }} &amp; {{ $t('common.mail') }}
                            </h2></template
                        >
                        <dl class="divide-y divide-default">
                            <div
                                v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'PhoneNumber').value"
                                class="flex items-start gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.phone_number') }}</dt>
                                <dd class="min-w-0 text-sm text-toned"><PhoneNumberBlock :number="user.phoneNumber" /></dd>
                            </div>
                            <div
                                v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Email').value"
                                class="flex items-start gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.mail', 1) }}</dt>
                                <dd class="min-w-0 text-sm text-toned"><EmailInfoPopover :email="user.props?.email" /></dd>
                            </div>
                            <div v-if="user.visum" class="flex items-center gap-4 px-4 py-3 sm:px-5">
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.visum') }}</dt>
                                <dd class="min-w-0 text-sm text-toned">{{ user.visum }}</dd>
                            </div>
                        </dl>
                    </UCard>

                    <UCard :ui="{ body: 'p-0 sm:p-0' }">
                        <template #header
                            ><h2 class="font-semibold text-highlighted">{{ $t('common.status') }}</h2></template
                        >
                        <dl class="divide-y divide-default">
                            <div
                                v-if="
                                    attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints')
                                        .value
                                "
                                class="flex items-center gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">
                                    {{ $t('common.traffic_infraction_points', 2) }}
                                </dt>
                                <dd
                                    class="min-w-0 text-sm text-toned"
                                    :class="
                                        (user.props?.trafficInfractionPoints ?? 0) >= 10 ? 'font-semibold text-error-500' : ''
                                    "
                                >
                                    {{ $t('common.point', user.props?.trafficInfractionPoints ?? 0) }}
                                </dd>
                            </div>
                            <div
                                v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                                class="flex items-center gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.fine', 2) }}</dt>
                                <dd class="min-w-0 text-sm text-toned">
                                    <span v-if="(user.props?.openFines ?? 0) <= 0">{{ $t('common.no_open_fine') }}</span
                                    ><span v-else class="font-semibold text-error-500">{{
                                        numberFormatter.format(user.props?.openFines ?? 0)
                                    }}</span>
                                </dd>
                            </div>
                            <div
                                v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Labels').value"
                                class="flex items-start gap-4 px-4 py-3 sm:px-5"
                            >
                                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ $t('common.label', 2) }}</dt>
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
                        v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'Licenses').value"
                        :ui="{ body: 'p-0 sm:p-0' }"
                    >
                        <template #header
                            ><h2 class="font-semibold text-highlighted">{{ $t('common.license', 2) }}</h2></template
                        >
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
