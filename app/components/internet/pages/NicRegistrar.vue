<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type { Tab } from '~/store/internet';
import type { TLD } from '~~/gen/ts/resources/internet/domain';
import type { CheckDomainAvailabilityResponse } from '~~/gen/ts/services/internet/domain';
import DomainList from './nic/DomainList.vue';
import RegisterForm from './nic/RegisterForm.vue';

const props = defineProps<{
    modelValue?: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Tab): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const tab = useVModel(props, 'modelValue', emit);

function updateTabInfo(): void {
    if (!tab.value) {
        return;
    }

    tab.value.label = t('components.internet.pages.nic_registrar.title');
    tab.value.icon = 'i-mdi-domain';
}
updateTabInfo();

const schema = z.object({
    tldId: z.number().positive().min(1),
    search: z.string().min(3).max(40),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    tldId: 1,
    search: '',
});

const { data: tlds } = useLazyAsyncData('internet-tlds', () => listTLDs());

async function listTLDs(): Promise<TLD[]> {
    try {
        const call = $grpc.internet.domain.listTLDs({});
        const { response } = await call;

        return response.tlds;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const {
    data: domainAvailability,
    error,
    refresh,
} = useLazyAsyncData(`internet-domain-check-${state.tldId}-${state.search}`, () => checkDomainAvailability(), {
    immediate: false,
});

async function checkDomainAvailability(): Promise<CheckDomainAvailabilityResponse> {
    try {
        const call = $grpc.internet.domain.checkDomainAvailability({
            tldId: state.tldId,
            name: state.search,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await refresh().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const items = [
    {
        label: t('components.internet.pages.nic_registrar.search.button'),
        icon: 'i-mdi-search',
        slot: 'search',
    },
    {
        label: t('components.internet.pages.nic_registrar.control_panel'),
        icon: 'i-mdi-administrator',
        slot: 'admin',
    },
];
</script>

<template>
    <UPage>
        <ULandingHero
            :title="$t('components.internet.pages.nic_registrar.title')"
            :description="$t('components.internet.pages.nic_registrar.description')"
            :ui="{
                wrapper: 'py-6 sm:py-16 md:py-16',
            }"
        />

        <ULandingSection :ui="{ wrapper: 'py-6 sm:py-6' }">
            <UTabs :items="items" :unmount="true">
                <template #search>
                    <UForm
                        :schema="schema"
                        :state="state"
                        class="mb-2 flex place-content-center gap-1"
                        @submit="onSubmitThrottle"
                    >
                        <UFormGroup name="search">
                            <UInput
                                v-model="state.search"
                                type="text"
                                class="w-full"
                                size="xl"
                                :placeholder="$t('common.domain', 1)"
                            />
                        </UFormGroup>

                        <UFormGroup name="tldId">
                            <USelectMenu
                                v-model="state.tldId"
                                :disabled="!tlds || tlds?.length === 0"
                                :options="tlds"
                                value-attribute="id"
                                option-attribute="name"
                                size="xl"
                            >
                                <template #label="{ selected }">
                                    <template v-if="selected">
                                        <span class="truncate">.{{ selected.name }}</span>
                                    </template>
                                </template>

                                <template #option="{ option: tld }">
                                    <span class="truncate">.{{ tld.name }}</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.tld', 2)]) }}
                                </template>
                            </USelectMenu>
                        </UFormGroup>

                        <UFormGroup>
                            <UButton
                                type="submit"
                                :label="$t('components.internet.pages.nic_registrar.search.button')"
                                trailing-icon="i-mdi-search"
                                size="xl"
                            />
                        </UFormGroup>
                    </UForm>

                    <DataErrorBlock v-if="error" :error="error" />
                    <UContainer v-else-if="domainAvailability !== undefined" class="flex flex-col gap-2">
                        <UAlert
                            v-if="!domainAvailability.transferable && !domainAvailability.available"
                            icon="i-mdi-information-outline"
                            color="red"
                            :title="$t('components.internet.pages.nic_registrar.search.not_available.title')"
                            :description="$t('components.internet.pages.nic_registrar.search.not_available.description')"
                        />
                        <template v-else>
                            <UAlert
                                v-if="domainAvailability.transferable"
                                icon="i-mdi-information-outline"
                                :title="$t('components.internet.pages.nic_registrar.search.transferable.title')"
                                :description="$t('components.internet.pages.nic_registrar.search.transferable.description')"
                            />
                            <UAlert
                                v-else-if="domainAvailability.available"
                                icon="i-mdi-information-outline"
                                :title="$t('components.internet.pages.nic_registrar.search.available.title')"
                                :description="$t('components.internet.pages.nic_registrar.search.available.description')"
                                color="green"
                            />

                            <RegisterForm :domain="state" :status="domainAvailability" />
                        </template>
                    </UContainer>
                </template>

                <template #admin>
                    <DomainList />
                </template>
            </UTabs>
        </ULandingSection>

        <ULandingSection :ui="{ wrapper: 'py-6 sm:py-6' }">
            <UPageGrid :ui="{ wrapper: 'sm:grid-cols-2 xl:grid-cols-2' }">
                <ULandingCard
                    icon="i-mdi-build"
                    :title="$t('components.internet.pages.nic_registrar.cards.builder.title')"
                    :description="$t('components.internet.pages.nic_registrar.cards.builder.description')"
                />

                <ULandingCard
                    icon="i-mdi-user-access-control"
                    :title="$t('components.internet.pages.nic_registrar.cards.access.title')"
                    :description="$t('components.internet.pages.nic_registrar.cards.access.description')"
                />
            </UPageGrid>
        </ULandingSection>
    </UPage>
</template>
