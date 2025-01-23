<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import type { Tab } from '~/store/internet';

const props = defineProps<{
    modelValue?: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Tab): void;
}>();

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
    search: z.string().max(40),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    search: '',
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    // TODO check availability
    await new Promise(() => undefined).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

// TODO
</script>

<template>
    <UPage>
        <ULandingHero
            :title="$t('components.internet.pages.nic_registrar.title')"
            :description="$t('components.internet.pages.nic_registrar.description')"
            :ui="{
                wrapper: 'py-6 sm:py-16 md:py-16',
            }"
        >
            <template #links>
                <UForm :schema="schema" :state="state" class="inline-flex gap-1" @submit="onSubmitThrottle">
                    <UFormGroup name="search">
                        <UInput v-model="state.search" type="text" class="w-full" size="xl" :placeholder="$t('common.domain')">
                            <template #trailing>
                                <span class="text-base text-gray-500 dark:text-gray-400">.ls</span>
                            </template>
                        </UInput>
                    </UFormGroup>

                    <UButton
                        type="submit"
                        :label="$t('components.internet.pages.nic_registrar.search.button')"
                        trailing-icon="i-mdi-search"
                        size="xl"
                    />
                </UForm>
            </template>
        </ULandingHero>

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
