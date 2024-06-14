<script lang="ts" setup>
import { z } from 'zod';
import { vMaska } from 'maska';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { attr } from '~/composables/can';
import { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import type { User } from '~~/gen/ts/resources/users/users';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { t } = useI18n();

const schema = z.object({
    name: z.string().optional(),
    phoneNumber: z.string().max(20).optional(),
    wanted: z.boolean().optional(),
    trafficInfractionPoints: z.coerce.number().nonnegative().optional(),
    openFines: z.coerce.number().nonnegative().optional(),
    dateofbirth: z.string().max(10).optional(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`citizens-${page.value}-${JSON.stringify(query.value)}`, () => listCitizens(), {
    transform: (input) => ({ ...input, users: wrapRows(input?.users, columns) }),
});

async function listCitizens(): Promise<ListCitizensResponse> {
    try {
        const req: ListCitizensRequest = {
            pagination: {
                offset: offset.value,
            },
            search: query.value.name ?? '',
        };
        if (query.value.wanted) {
            req.wanted = query.value.wanted;
        }
        if (query.value.phoneNumber) {
            req.phoneNumber = query.value.phoneNumber;
        }
        if (query.value.trafficInfractionPoints) {
            req.trafficInfractionPoints = query.value.trafficInfractionPoints ?? 0;
        }
        if (query.value.openFines) {
            req.openFines = query.value.openFines ?? 0;
        }
        if (query.value.dateofbirth) {
            req.dateofbirth = query.value.dateofbirth;
        }

        const call = getGRPCCitizenStoreClient().listCitizens(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 200, maxWait: 1250 });

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

function addToClipboard(user: User): void {
    clipboardStore.addUser({
        ...user,
        // @ts-expect-error wrapped table rows
        jobLabel: user.jobLabel.value,
        // @ts-expect-error wrapped table rows
        sex: user.sex.value,
        // @ts-expect-error wrapped table rows
        dateofbirth: user.dateofbirth.value,
        // @ts-expect-error wrapped table rows
        height: user.height.value,
    });

    notifications.add({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'jobLabel',
        label: t('common.job'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'sex',
        label: t('common.sex'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber').value
        ? { key: 'phoneNumber', label: t('common.phone_number') }
        : undefined,
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value
        ? {
              key: 'trafficInfractionPoints',
              label: t('common.traffic_infraction_points', 2),
          }
        : undefined,
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines').value
        ? {
              key: 'openFines',
              label: t('common.fine', 2),
          }
        : undefined,
    {
        key: 'height',
        label: t('common.height'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    can('CitizenStoreService.GetUser').value
        ? {
              key: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined) as { key: string; label: string; class?: string; rowClass?: string; sortable?: boolean }[];

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <UDashboardToolbar>
        <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
            <div class="flex w-full flex-row gap-2">
                <UFormGroup class="flex-1" :label="`${$t('common.citizen', 1)} ${$t('common.name')}`">
                    <UInput
                        ref="input"
                        v-model="query.name"
                        type="text"
                        name="name"
                        :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                        block
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                        @keydown.esc="$event.target.blur()"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup name="dateofbirth" :label="$t('common.date_of_birth')">
                    <UInput
                        v-model="query.dateofbirth"
                        v-maska
                        type="text"
                        name="dateofbirth"
                        data-maska="##.##.####"
                        :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                        block
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>

                <UFormGroup
                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Wanted').value"
                    name="wanted"
                    :label="$t('components.citizens.CitizensList.only_wanted')"
                >
                    <UToggle v-model="query.wanted">
                        <span class="sr-only">
                            {{ $t('components.citizens.CitizensList.only_wanted') }}
                        </span>
                    </UToggle>
                </UFormGroup>
            </div>

            <UAccordion
                class="mt-2"
                color="white"
                variant="soft"
                size="sm"
                :items="[{ label: $t('common.advanced_search'), slot: 'search' }]"
            >
                <template #search>
                    <div class="flex flex-row gap-2">
                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber').value"
                            name="phoneNumber"
                            :label="$t('common.phone_number')"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.phoneNumber"
                                type="tel"
                                name="phoneNumber"
                                :placeholder="$t('common.phone_number')"
                                block
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'TrafficInfractionPoints').value"
                            name="trafficInfractionPoints"
                            :label="$t('common.traffic_infraction_points', 2)"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.trafficInfractionPoints"
                                type="number"
                                name="trafficInfractionPoints"
                                min="0"
                                :placeholder="$t('common.traffic_infraction_points')"
                                block
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                            name="openFines"
                            :label="$t('components.citizens.CitizensList.open_fine')"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.openFines"
                                type="number"
                                name="openFines"
                                min="0"
                                :placeholder="`${$t('common.fine')}`"
                                block
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </UAccordion>
        </UForm>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.citizen', 2)])" :retry="refresh" />
    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="data?.users"
        :empty-state="{ icon: 'i-mdi-accounts', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
        class="flex-1"
    >
        <template #name-data="{ row: citizen }">
            <div class="inline-flex items-center gap-1 text-gray-900 dark:text-white">
                <ProfilePictureImg
                    :src="citizen.props?.mugShot?.url"
                    :name="`${citizen.firstname} ${citizen.lastname}`"
                    :alt="$t('common.mug_shot')"
                    :enable-popup="true"
                    size="sm"
                />

                <span>{{ citizen.firstname }} {{ citizen.lastname }}</span>
                <span class="lg:hidden"> ({{ citizen.dateofbirth.value }}) </span>

                <UBadge v-if="citizen.props?.wanted" color="red">
                    {{ $t('common.wanted').toUpperCase() }}
                </UBadge>
            </div>

            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                <dd class="mt-1 truncate">
                    {{ citizen.sex?.value.toUpperCase() ?? $t('common.na') }} -
                    {{ citizen.jobLabel.value ?? $t('common.na') }}
                </dd>
            </dl>
        </template>
        <template #jobLabel-data="{ row: citizen }">
            {{ citizen.jobLabel.value }}
        </template>
        <template #sex-data="{ row: citizen }">
            {{ citizen.sex?.value.toUpperCase() ?? $t('common.na') }}
        </template>
        <template #phoneNumber-data="{ row: citizen }">
            <PhoneNumberBlock :number="citizen.phoneNumber" />
        </template>
        <template #openFines-data="{ row: citizen }">
            <template v-if="(citizen.props?.openFines ?? 0) > 0">
                {{ $n(parseInt((citizen.props?.openFines ?? 0).toString()), 'currency') }}
            </template>
        </template>
        <template #dateofbirth-data="{ row: citizen }">
            {{ citizen.dateofbirth.value }}
        </template>
        <template #height-data="{ row: citizen }"> {{ citizen.height.value }}cm </template>
        <template v-if="can('CitizenStoreService.GetUser').value" #actions-data="{ row: citizen }">
            <div :key="citizen.userId" class="flex flex-col justify-end md:flex-row">
                <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(citizen)" />

                <UButton
                    variant="link"
                    icon="i-mdi-eye"
                    :to="{
                        name: 'citizens-id',
                        params: { id: citizen.userId ?? 0 },
                    }"
                />
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
