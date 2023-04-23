<script lang="ts" setup>
import { CardElement } from '~/utils/types';
import Cards from './Cards.vue';
import { useUserSettingsStore } from '~/store/usersettings';
import { LocaleObject } from 'vue-i18n-routing';

const { locale, locales, setLocale } = useI18n();

const router = useRouter();
const store = useUserSettingsStore();

const langs = new Array<CardElement>();
locales.value.forEach(l => {
    l = l as LocaleObject;
    langs.push({
        title: l.name!,
        description: `ISO-Code: ${l.iso}`,
    });
});

async function selected(idx: number): Promise<void> {
    const l = locales.value[idx] as LocaleObject;
    store.setLocale(l.iso!);
    setLocale(l.iso!);

    await router.push({ name: 'index' });
}
</script>

<template>
    <Cards :items="langs" :show-icon="false" @selected="selected($event)" />
</template>
