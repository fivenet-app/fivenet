<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
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

const { $grpc } = useNuxtApp();

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
    if (values.search === '') {
        return {
            results: [],
        };
    }

    try {
        const call = $grpc.internet.internet.search({
            search: values.search,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    const response = await searchInternet(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    searchResults.value = response;
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
                wrapper: 'py-6 sm:py-6 md:py-6',
            }"
        >
            <template #title>
                <img src="/images/components/internet/homepage/logo.png" :alt="$t('common.logo')" class="mx-auto h-48" />
            </template>
        </ULandingHero>

        <ULandingSection :ui="{ wrapper: 'pt-0 sm:pt-0 pb-6 sm:pb-6' }">
            <div class="flex w-full flex-1 flex-col items-center gap-1">
                <UForm :schema="schema" :state="state" class="inline-flex gap-1" @submit="onSubmitThrottle">
                    <UFormGroup name="search">
                        <UInput
                            v-model="state.search"
                            type="text"
                            class="w-full"
                            size="xl"
                            :placeholder="$t('common.search')"
                            :disabled="!canSubmit"
                            autocomplete="off"
                            :ui="{ icon: { trailing: { pointer: '' } } }"
                        >
                            <template #trailing>
                                <UButton
                                    v-show="state.search !== ''"
                                    color="gray"
                                    variant="link"
                                    icon="i-heroicons-x-mark-20-solid"
                                    :padded="false"
                                    @click="state.search = ''"
                                />
                            </template>
                        </UInput>
                    </UFormGroup>

                    <UButton type="submit" icon="i-mdi-search" size="xl" :disabled="!canSubmit" />
                </UForm>

                <div>
                    <p v-if="searchResults" class="text-sm">
                        {{ $t('common.found', [$t('common.search_results', searchResults.results.length)]) }}
                    </p>
                </div>

                <div
                    v-if="searchResults"
                    class="grid w-full max-w-lg grid-cols-1 divide-y divide-gray-200 dark:divide-gray-800"
                >
                    <UCard
                        v-for="result in searchResults.results"
                        :key="result.id"
                        :ui="{
                            divide: '',
                            ring: '',
                            shadow: '',
                            header: { padding: 'px-4 py-1 sm:p-1 sm:px-4' },
                            body: { padding: 'px-4 py-1 sm:p-1 sm:px-4' },
                            footer: { padding: 'px-4 py-1 sm:p-1 sm:px-4' },
                        }"
                    >
                        <template #header>
                            <ULink
                                class="inline-flex items-center gap-1 hover:underline"
                                @click="internetStore.goTo(`${result.domain?.name}.${result.domain?.tld?.name}`, result.path)"
                            >
                                <UIcon name="i-mdi-web" class="size-6" />
                                <h3 class="text-xl font-semibold">{{ result.title }}</h3>
                            </ULink>
                        </template>

                        <p class="text-sm">{{ result.description }}</p>
                    </UCard>
                </div>
            </div>
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
