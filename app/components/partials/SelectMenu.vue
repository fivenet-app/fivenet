<script
    lang="ts"
    setup
    generic="
        T extends ArrayOrNested<SelectMenuItem>,
        VK extends GetItemKeys<T> | undefined = undefined,
        M extends boolean = false,
        Mod extends Omit<ModelModifiers, 'lazy'> = Omit<ModelModifiers, 'lazy'>,
        C extends boolean | object = boolean | object
    "
>
import type { SelectMenuEmits, SelectMenuItem, SelectMenuProps, SelectMenuSlots } from '@nuxt/ui';
import type { ModelModifiers } from '@nuxt/ui/runtime/types/input.js';
import type { ArrayOrNested, GetItemKeys, NestedItem } from '@nuxt/ui/runtime/types/utils.js';

interface Props<
    T extends ArrayOrNested<SelectMenuItem>,
    VK extends GetItemKeys<T> | undefined = undefined,
    M extends boolean = false,
    Mod extends Omit<ModelModifiers, 'lazy'> = Omit<ModelModifiers, 'lazy'>,
    C extends boolean | object = boolean | object,
> extends /* @vue-ignore */ SelectMenuProps<T, VK, M, Mod, C> {
    searchableKey?: string;
    searchable?: (q: string) => Promise<T>;
}

type SelectMenuDefaultSlotProps<
    A extends ArrayOrNested<SelectMenuItem>,
    VK extends GetItemKeys<A> | undefined,
    M extends boolean,
    Mod extends Omit<ModelModifiers, 'lazy'>,
    C extends boolean | object,
> = Parameters<NonNullable<SelectMenuSlots<A, VK, M, Mod, C>['default']>>[0];

interface Slots<
    A extends ArrayOrNested<SelectMenuItem> = ArrayOrNested<SelectMenuItem>,
    VK extends GetItemKeys<A> | undefined = undefined,
    M extends boolean = false,
    Mod extends Omit<ModelModifiers, 'lazy'> = Omit<ModelModifiers, 'lazy'>,
    C extends boolean | object = boolean | object,
    T extends NestedItem<A> = NestedItem<A>,
> extends /* @vue-ignore */ SelectMenuSlots<A, VK, M, Mod, C> {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    'default'(props: SelectMenuDefaultSlotProps<A, VK, M, Mod, C> & { items?: T[] }): any;
}

const props = defineProps<Props<T, VK, M, Mod, C>>();
defineSlots</* @vue-ignore */ Slots<T, VK, M, Mod, C>>();
defineEmits</* @vue-ignore */ SelectMenuEmits<T, VK, M, Mod, C>>();

const loading = ref(false);

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 175);

const { data: items } = useLazyAsyncData(
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
        <USelectMenu
            v-model:search-term="searchTerm"
            :loading="loading"
            :items="props.items ?? items"
            ignore-filter
            v-bind="$attrs"
        >
            <template v-for="(_, name) in $slots" #[name]="slotProps">
                <!-- @vue-expect-error to type or not to type, the `name` attribute is correct but not correct -->
                <slot :name="name" v-bind="name === 'default' ? { ...slotProps, items } : slotProps" />
            </template>
        </USelectMenu>
    </ClientOnly>
</template>
