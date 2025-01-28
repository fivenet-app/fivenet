<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getGRPCInternetClient } from '~/composables/grpc';
import { useInternetStore, type Tab } from '~/store/internet';
import { AdType } from '~~/gen/ts/resources/internet/ads';
import type { SearchResponse } from '~~/gen/ts/services/internet/internet';
import DataPendingBlock from '../../partials/data/DataPendingBlock.vue';
import SponsoredAdCard from '../ads/SponsoredAdCard.vue';

const props = defineProps<{
    modelValue?: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Tab): void;
}>();

const { t } = useI18n();

const tab = useVModel(props, 'modelValue', emit);

const internetStore = useInternetStore();

function updateTabInfo(): void {
    if (!tab.value) {
        return;
    }

    tab.value.label = t('components.internet.pages.homepage.title');
    tab.value.icon = 'i-mdi-home';
}

updateTabInfo();
watch(tab, () => updateTabInfo());

const schema = z.object({
    search: z.string().max(40),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    search: '',
});

async function searchInternet(values: Schema): Promise<SearchResponse> {
    try {
        const call = getGRPCInternetClient().search({
            search: values.search,
        });
        const { response } = await call;

        searchResults.value = response;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await searchInternet(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const searchResults = ref<undefined | SearchResponse>(undefined);

const { data: ads, pending: loadingAds } = useLazyAsyncData(`internet-ads`, () =>
    internetStore.getAds({
        adType: AdType.SPONSORED,
        count: 3,
    }),
);
</script>

<template>
    <UPage>
        <ULandingHero
            :description="$t('components.internet.pages.homepage.description')"
            :ui="{
                wrapper: 'py-6 sm:py-16 md:py-16',
            }"
        >
            <template #title>
                <img src="/images/components/internet/homepage/logo.png" :alt="$t('common.logo')" class="mx-auto h-48" />
            </template>

            <template #links>
                <UForm :schema="schema" :state="state" class="inline-flex gap-1" @submit="onSubmitThrottle">
                    <UFormGroup name="search">
                        <UInput
                            v-model="state.search"
                            type="text"
                            class="w-full"
                            size="xl"
                            :placeholder="$t('common.search')"
                            :disabled="!canSubmit"
                        />
                    </UFormGroup>

                    <UButton type="submit" icon="i-mdi-search" size="xl" :disabled="!canSubmit" />
                </UForm>
            </template>
        </ULandingHero>

        <ULandingSection v-if="searchResults" :ui="{ wrapper: 'py-6 sm:py-6' }">
            <DataNoDataBlock v-if="searchResults.results.length === 0" :type="$t('common.result', 2)" />

            <template v-else>
                <UCard v-for="result in searchResults.results" :key="result.url">
                    <template #header>
                        <h3 class="text-xl font-semibold">{{ result.title }}</h3>
                    </template>

                    <p class="">{{ result.description }}</p>

                    <template #footer>
                        <p class="text-sm">{{ result.url }}</p>
                    </template>
                </UCard>
            </template>
        </ULandingSection>

        <ULandingSection :ui="{ wrapper: 'py-6 sm:py-6' }">
            <DataPendingBlock v-if="loadingAds" :message="$t('common.loading', [$t('common.ad', 1)])" />
            <UPageGrid v-else>
                <SponsoredAdCard v-for="(ad, idx) in ads?.ads" :key="idx" :ad="ad" />
            </UPageGrid>
        </ULandingSection>

        <ULandingSection :ui="{ wrapper: 'py-3 sm:py-3 mx-auto max-w-lg' }">
            <ULandingCard
                icon="i-mdi-web-box"
                orientation="horizontal"
                :title="$t('components.internet.pages.homepage.nic_registrar_cta.title')"
            >
                <template #description>
                    <p>{{ $t('components.internet.pages.homepage.nic_registrar_cta.description') }}</p>

                    <UButton
                        :label="$t('components.internet.pages.homepage.nic_registrar_cta.button')"
                        @click="internetStore.goTo('nic.ls')"
                    />
                </template>
            </ULandingCard>
        </ULandingSection>
    </UPage>
</template>
