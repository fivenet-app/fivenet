<script lang="ts" setup>
import { CardElement } from '~/utils/types';
import Cards from './Cards.vue';
//import { useUserSettingsStore } from '~/store/usersettings';

/*const { locales, setLocale } = useI18n();

const router = useRouter();
const store = useUserSettingsStore();*/

const langs = new Array<CardElement>();
/*locales.value.forEach(l => {
    langs.push({
        title: l.name,
        description: `ISO-Code: ${l.iso}`,
    });
});*/

async function selected(idx: number): Promise<void> {
    /*store.setLocale(locales.value[idx].iso);
    setLocale(locales.value[idx].iso);

    await router.push({ name: 'index' });*/
}
</script>

<template>
    <Cards :items="langs" :show-icon="false" @selected="selected($event)" />
</template>
