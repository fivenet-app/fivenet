<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await searchInternet(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

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

        <ULandingSection v-if="searchResults && searchResults.results.length > 0" :ui="{ wrapper: 'py-6 sm:py-6' }">
            <!-- TODO show search results -->
        </ULandingSection>

        <ULandingSection :ui="{ wrapper: 'py-6 sm:py-6' }">
            <DataPendingBlock v-if="loadingAds" :message="$t('common.loading', [$t('common.ad', 1)])" />
            <UPageGrid v-else>
                <SponsoredAdCard v-for="(ad, idx) in ads?.ads" :key="idx" :ad="ad" />
            </UPageGrid>
        </ULandingSection>
    </UPage>
</template>
