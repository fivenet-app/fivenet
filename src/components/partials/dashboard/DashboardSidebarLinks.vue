<!-- This is a modified version of Nuxt UI Pro's DashboardSidebarLinks.vue -->
<template>
    <ul
        :class="ui.wrapper"
        v-bind="{
            ...attrs,
        }"
        @touchend="fixActionRestriction"
    >
        <li v-for="(link, index) in links" :key="index" tag="li" :class="ui.container">
            <component
                :is="link.children?.length ? Disclosure : 'div'"
                v-slot="slotProps"
                :default-open="link.defaultOpen === undefined ? true : link.defaultOpen"
                as="div"
            >
                <component :is="link.children?.length ? DisclosureButton : 'div'" as="template">
                    <UTooltip
                        class="flex"
                        :popper="{ placement: 'right' }"
                        :prevent="!link.tooltip"
                        :ui="ui.tooltip"
                        v-bind="link.tooltip"
                    >
                        <ULink
                            v-slot="{ isActive }"
                            v-bind="
                                link.children?.length && link.collapsible !== false
                                    ? { disabled: link.disabled }
                                    : getULinkProps(link)
                            "
                            :class="[ui.base]"
                            :active-class="ui.active"
                            :inactive-class="
                                !link.to && link.collapsible === false && level === 0 && link.children?.length
                                    ? ui.static
                                    : ui.inactive
                            "
                            @click="link.click"
                        >
                            <slot name="icon" :link="link" :is-active="isActive">
                                <UIcon
                                    v-if="link.icon"
                                    :name="link.icon"
                                    :class="
                                        twMerge(
                                            twJoin(
                                                ui.icon.base,
                                                isActive
                                                    ? ui.icon.active
                                                    : !link.to &&
                                                        link.collapsible === false &&
                                                        level === 0 &&
                                                        link.children?.length
                                                      ? ui.static
                                                      : ui.icon.inactive,
                                            ),
                                            link.iconClass,
                                        )
                                    "
                                />
                                <UAvatar
                                    v-else-if="link.avatar"
                                    v-bind="{
                                        size: ui.avatar.size,
                                        ...link.avatar,
                                    }"
                                    :class="twMerge(twJoin(ui.avatar.base), link.avatarClass)"
                                />
                                <UChip
                                    v-else-if="link.chip"
                                    v-bind="{
                                        size: ui.chip.size,
                                        ...(typeof link.chip === 'string' ? { color: link.chip as ChipColor } : link.chip),
                                    }"
                                    :class="twMerge(twJoin(ui.chip.base), link.chipClass)"
                                />

                                <span v-else-if="level > 0" :class="[ui.dot.wrapper, index < links.length - 1 && ui.dot.after]">
                                    <span :class="[ui.dot.base, isActive ? ui.dot.active : ui.dot.inactive]" />
                                </span>
                            </slot>

                            <slot :link="link" :is-active="isActive">
                                <span v-if="link.label" :class="twMerge(ui.label, link.labelClass)">
                                    <span v-if="isActive" class="sr-only"> Current page: </span>

                                    {{ link.label }}
                                </span>
                            </slot>

                            <UIcon
                                v-if="link.children?.length && link.collapsible !== false"
                                :name="ui.trailingIcon.name"
                                :class="[
                                    ui.trailingIcon.base,
                                    slotProps?.open ? ui.trailingIcon.active : ui.trailingIcon.inactive,
                                ]"
                            />

                            <slot name="badge" :link="link" :is-active="isActive">
                                <UBadge
                                    v-if="link.badge"
                                    v-bind="{
                                        size: ui.badge.size,
                                        color: ui.badge.color,
                                        variant: ui.badge.variant,
                                        ...(typeof link.badge === 'string' || typeof link.badge === 'number'
                                            ? { label: link.badge }
                                            : link.badge),
                                    }"
                                    :class="ui.badge.base"
                                />
                            </slot>
                        </ULink>
                    </UTooltip>
                </component>

                <Transition
                    v-bind="ui.transition"
                    @enter="onEnter"
                    @after-enter="onAfterEnter"
                    @before-leave="onBeforeLeave"
                    @leave="onLeave"
                >
                    <DisclosurePanel
                        v-if="link.children?.length && (slotProps?.open || link.collapsible === false)"
                        static
                        as="template"
                    >
                        <UDashboardSidebarLinks
                            :level="level + 1"
                            :links="link.children"
                            :ui="ui"
                            @update:links="$emit('update:links', $event)"
                        >
                            <template v-for="(_, name) in $slots" #[name]="slotData: any">
                                <slot :name="name" v-bind="slotData" />
                            </template>
                        </UDashboardSidebarLinks>
                    </DisclosurePanel>
                </Transition>
            </component>
        </li>
    </ul>
</template>

<script setup lang="ts">
import { Disclosure, DisclosureButton, DisclosurePanel, provideUseId } from '@headlessui/vue';
import type { PropType } from 'vue';
// @ts-ignore
import { useId } from '#imports';
import type { DashboardSidebarLink } from '#ui-pro/types';
import type { ChipColor } from '#ui/types';
import { getULinkProps } from '#ui/utils';
import { twJoin, twMerge } from 'tailwind-merge';

const appConfig = useAppConfig();

const config = computed(() => ({
    wrapper: 'relative !min-h-[auto] !min-w-[auto]',
    container: '!overflow-visible',
    base: 'group relative flex items-center gap-1.5 px-2.5 py-1.5 w-full rounded-md font-medium text-sm focus:outline-none focus-visible:outline-none dark:focus-visible:outline-none focus-visible:before:ring-inset focus-visible:before:ring-2 focus-visible:before:ring-primary-500 dark:focus-visible:before:ring-primary-400 before:absolute before:inset-px before:rounded-md disabled:cursor-not-allowed disabled:opacity-75',
    active: 'text-gray-900 dark:text-white before:bg-gray-100 dark:before:bg-gray-800',
    inactive:
        'text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:before:bg-gray-50 dark:hover:before:bg-gray-800/50',
    static: 'text-gray-900 dark:text-white cursor-auto',
    icon: {
        base: 'flex-shrink-0 w-5 h-5 relative',
        active: 'text-gray-900 dark:text-white',
        inactive: 'text-gray-400 dark:text-gray-500 group-hover:text-gray-700 dark:group-hover:text-gray-200',
    },
    trailingIcon: {
        name: appConfig.ui.icons.chevron,
        base: 'ml-auto w-5 h-5 transform transition-transform duration-200 flex-shrink-0',
        active: '',
        inactive: '-rotate-90',
    },
    avatar: {
        base: 'flex-shrink-0',
        size: '2xs' as const,
    },
    chip: {
        base: 'flex-shrink-0 mx-2.5',
        size: 'sm' as const,
    },
    badge: {
        base: 'flex-shrink-0 ml-auto relative rounded',
        color: 'gray' as const,
        variant: 'solid' as const,
        size: 'xs' as const,
    },
    label: 'text-sm truncate relative',
    dot: {
        wrapper: 'w-px h-full mx-[9.5px] bg-gray-200 dark:bg-gray-700 relative',
        after: 'after:absolute after:z-[1] after:w-px after:h-full after:bg-gray-200 after:dark:bg-gray-700 after:transform after:translate-y-full after:inset-x-0',
        base: 'w-1 h-1 rounded-full absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2',
        active: 'bg-gray-900 dark:bg-white',
        inactive: 'bg-gray-400 dark:bg-gray-500 group-hover:bg-gray-700 dark:group-hover:bg-gray-200',
    },
    tooltip: {
        strategy: 'override',
        transition: {
            enterActiveClass: 'transition ease-out duration-200',
            enterFromClass: 'opacity-0',
            enterToClass: 'opacity-100',
            leaveActiveClass: 'transition ease-in duration-150',
            leaveFromClass: 'opacity-100',
            leaveToClass: 'opacity-0',
        },
    },
    transition: {
        enterActiveClass: 'overflow-hidden transition-[height] duration-200 ease-out',
        leaveActiveClass: 'overflow-hidden transition-[height] duration-200 ease-out',
    },
}));

defineOptions({
    inheritAttrs: false,
});

const props = defineProps({
    level: {
        type: Number,
        default: 0,
    },
    links: {
        type: Array as PropType<DashboardSidebarLink[]>,
        default: () => [],
    },
    class: {
        type: [String, Object, Array] as PropType<any>,
        default: undefined,
    },
    ui: {
        type: Object as PropType<Partial<typeof config.value>>,
        default: () => ({}),
    },
});

const { ui, attrs } = useUI('dashboard.sidebar.links', toRef(props, 'ui'), config, toRef(props, 'class'), true);

function onEnter(_el: Element, done: () => void) {
    const el = _el as HTMLElement;
    el.style.height = '0';
    el.offsetHeight; // Trigger a reflow, flushing the CSS changes
    el.style.height = el.scrollHeight + 'px';

    el.addEventListener('transitionend', done, { once: true });
}

function onBeforeLeave(_el: Element) {
    const el = _el as HTMLElement;
    el.style.height = el.scrollHeight + 'px';
    el.offsetHeight; // Trigger a reflow, flushing the CSS changes
}

function onAfterEnter(_el: Element) {
    const el = _el as HTMLElement;
    el.style.height = 'auto';
}

function onLeave(_el: Element, done: () => void) {
    const el = _el as HTMLElement;
    el.style.height = '0';

    el.addEventListener('transitionend', done, { once: true });
}

function fixActionRestriction() {
    document.body.classList.remove('smooth-dnd-no-user-select', 'smooth-dnd-disable-touch-action');
}

provideUseId(() => useId());
</script>
