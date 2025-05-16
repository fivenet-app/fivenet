<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { splitURL, urlHomePage } from '~/components/internet/helper';
import { useInternetStore, type Tab } from '~/stores/internet';
import type { GetPageResponse } from '~~/gen/ts/services/internet/internet';
import HomePage from './pages/HomePage.vue';
import NotFound from './pages/NotFound.vue';
import WebPage from './pages/WebPage.vue';
import { localPages, localPagesDomains } from './pages/helpers';

const props = defineProps<{
    modelValue: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', tab: Tab): void;
}>();

const tab = useVModel(props, 'modelValue', emit);

const { $grpc } = useNuxtApp();

const internetStore = useInternetStore();

const {
    data: page,
    status,
    refresh,
} = useLazyAsyncData(`internet-page-${tab.value?.domain}:${tab.value?.path}`, () => getPage(), {
    immediate: false,
});

async function getPage(): Promise<GetPageResponse | undefined> {
    if (localPages[tab.value.domain]) {
        return;
    }

    try {
        const call = $grpc.internet.internet.getPage({
            domain: tab.value.domain,
            path: tab.value.path,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        const err = e as RpcError;

        handleGRPCError(err);
        throw err;
    }
}

watch(
    tab,
    async () => {
        state.url = tab.value.domain + (tab.value.path && tab.value.path !== '' ? tab.value.path : '/');

        // Skip local pages
        if (localPagesDomains.includes(tab.value?.domain)) {
            return;
        }

        if (tab.value.active) {
            refresh();
        }
    },
    { deep: true },
);

function goToPage(domain: string, path?: string): void {
    state.url = domain + (path && path !== '' ? path : '/');
    internetStore.goTo(domain, path);
}

const schema = z.object({
    url: z.string().max(128),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    url: tab.value ? tab.value.domain + (tab.value.path && tab.value.path !== '' ? tab.value.path : '') : '',
});

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    const split = splitURL(event.data.url);
    if (!split) {
        return;
    }

    goToPage(split.domain.toLowerCase(), split.path?.toLowerCase());
}, 500);

const input = useTemplateRef('input');
</script>

<template>
    <UDashboardToolbar
        :ui="{ wrapper: 'min-h-[41px] bg-gray-100 dark:bg-gray-800 p-0 gap-x-0', container: 'gap-x-0 gap-y-0 mx-1' }"
    >
        <UForm class="flex flex-1 items-center gap-1" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UButtonGroup>
                <UTooltip :text="$t('common.back')">
                    <UButton
                        variant="ghost"
                        color="white"
                        icon="i-mdi-chevron-left"
                        :disabled="tab.history.length === 0 || status === 'pending'"
                        @click="internetStore.back()"
                    />
                </UTooltip>

                <UTooltip :text="$t('common.home')">
                    <UButton
                        :disabled="tab.domain === urlHomePage || status === 'pending'"
                        variant="ghost"
                        color="white"
                        icon="i-mdi-home"
                        @click="goToPage(urlHomePage)"
                    />
                </UTooltip>

                <UTooltip :text="$t('common.refresh')">
                    <UButton
                        variant="ghost"
                        color="white"
                        icon="i-mdi-refresh"
                        :disabled="status === 'pending'"
                        :loading="status === 'pending'"
                        @click="refresh"
                    />
                </UTooltip>
            </UButtonGroup>

            <UInput ref="input" v-model="state.url" class="flex-1" type="text" :ui="{ icon: { trailing: { pointer: '' } } }">
                <template #trailing>
                    <UButton
                        v-show="state.url !== ''"
                        color="gray"
                        variant="link"
                        icon="i-mdi-close"
                        :padded="false"
                        @click="
                            state.url = '';
                            input?.input.focus();
                        "
                    />
                </template>
            </UInput>
        </UForm>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <HomePage v-if="tab?.domain === urlHomePage || tab?.domain === ''" v-model="tab" />
        <template v-else-if="tab?.domain && localPages[tab.domain]">
            <component :is="localPages[tab.domain]" v-model="tab" />
        </template>
        <NotFound v-else-if="!page?.page" v-model="tab" @refresh="refresh" />
        <WebPage v-else v-model="tab" :page="page.page" />
    </UDashboardPanelContent>
</template>
