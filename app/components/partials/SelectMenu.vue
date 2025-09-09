<script
    lang="ts"
    setup
    generic="
        T extends ArrayOrNested<SelectMenuItem>,
        VK extends GetItemKeys<T> | undefined = undefined,
        M extends boolean = false
    "
>
import type {
    ArrayOrNested,
    GetItemKeys,
    GetModelValue,
    NestedItem,
    SelectMenuEmits,
    SelectMenuItem,
    SelectMenuProps,
    SelectMenuSlots,
} from '@nuxt/ui';

interface Props<
    T extends ArrayOrNested<SelectMenuItem>,
    VK extends GetItemKeys<T> | undefined = undefined,
    M extends boolean = false,
> extends /* @vue-ignore */ SelectMenuProps<T, VK, M> {
    searchableKey?: string;
    searchable?: (q: string) => Promise<T>;
}

interface Slots<
    A extends ArrayOrNested<SelectMenuItem> = ArrayOrNested<SelectMenuItem>,
    VK extends GetItemKeys<A> | undefined = undefined,
    M extends boolean = false,
    T extends NestedItem<A> = NestedItem<A>,
> extends /* @vue-ignore */ SelectMenuSlots<A, VK, M> {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    'default'(props: { modelValue?: GetModelValue<A, VK, M>; open: boolean; items: T[] }): any;
}

const props = defineProps<Props<T, VK, M>>();
defineSlots</* @vue-ignore */ Slots<T, VK, M>>();
defineEmits</* @vue-ignore */ SelectMenuEmits<T, VK, M>>();

const loading = ref(false);

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

const { data: items } = await useAsyncData(
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
        <USelectMenu v-model:search-term="searchTerm" :loading="loading" :items="props.items ?? items" v-bind="$attrs">
            <template v-for="(_, name) in $slots" #[name]="slotProps">
                <!-- @vue-expect-error to type or not to type, the `name` attribute is correct but not correct -->
                <slot :name="name" v-bind="name === 'default' ? { ...slotProps, items } : slotProps" />
            </template>
        </USelectMenu>
    </ClientOnly>
</template>
