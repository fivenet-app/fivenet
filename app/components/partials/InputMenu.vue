<script
    lang="ts"
    setup
    generic="
        // Generic
        A extends ArrayOrNested<InputMenuItem>,
        VK extends GetItemKeys<A> | undefined = undefined,
        M extends boolean = false,
        Mod extends Omit<ModelModifiers, 'lazy'> = Omit<ModelModifiers, 'lazy'>,
        C extends boolean | object = boolean | object,
        T extends NestedItem<A> = NestedItem<A>
    "
>
import type { InputMenuEmits, InputMenuItem, InputMenuProps, InputMenuSlots } from '@nuxt/ui';
import type { ModelModifiers } from '@nuxt/ui/runtime/types/input.js';
import type { ArrayOrNested, GetItemKeys, NestedItem } from '@nuxt/ui/runtime/types/utils.js';

interface Props<
    T extends ArrayOrNested<InputMenuItem>,
    VK extends GetItemKeys<T> | undefined = undefined,
    M extends boolean = false,
    Mod extends Omit<ModelModifiers, 'lazy'> = Omit<ModelModifiers, 'lazy'>,
    C extends boolean | object = boolean | object,
> extends /* @vue-ignore */ InputMenuProps<T, VK, M, Mod, C> {
    searchableKey?: string;
    searchable?: (q: string) => Promise<T>;
}

const props = defineProps<Props<A, VK, M, Mod, C>>();
defineSlots</* @vue-ignore */ InputMenuSlots<A, VK, M, Mod, C, T>>();
defineEmits</* @vue-ignore */ InputMenuEmits<A, VK, M, Mod, C>>();

const loading = ref(false);

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

const { data } = await useAsyncData(
    () => `${props.searchableKey}-${searchTermDebounced.value}`,
    async () => {
        if (props.searchable === undefined) return [];

        loading.value = true;
        const items = await props.searchable(searchTermDebounced.value);
        loading.value = false;

        return items;
    },
    {
        immediate: !!props.searchable && props.items === undefined,
    },
);
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model:search-term="searchTerm"
            :loading="loading"
            :items="props.items ?? data"
            ignore-filter
            v-bind="$attrs"
        >
            <template v-for="(_, name) in $slots" #[name]="slotProps">
                <!-- @vue-expect-error to type or not to type, the `name` attribute is correct but not correct -->
                <slot :name="name" v-bind="slotProps" />
            </template>
        </UInputMenu>
    </ClientOnly>
</template>
