<script lang="ts">
import { NuxtImg, UCheckbox, UIcon } from '#components';
import { defineComponent, getCurrentInstance, h, Text, type Component, type VNode } from 'vue';
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';

export default defineComponent({
    name: 'HTMLContentRenderer',

    props: {
        value: {
            type: Object as () => JSONNode,
            required: true,
        },
    },

    setup(props) {
        // Self-reference for recursion
        const self = getCurrentInstance()?.type as Component;

        // Optional tag remapping
        const tagRemapping: Record<string, Component> = {
            img: NuxtImg,
        };

        return (): VNode => {
            const value = props.value;

            // 1. Text Node
            if (value.text) {
                return h(Text, null, value.text);
            }

            // 2. Checkbox input
            if (value.tag === 'input' && value.attrs?.type === 'checkbox') {
                return h(UCheckbox, {
                    disabled: true,
                    modelValue: !!value.attrs.checked,
                    ui: {
                        wrapper: '',
                        container: '',
                        base: 'h-5 w-5',
                    },
                });
            }

            // 3. <br> tag
            if (value.tag === 'br') {
                return h('br', value.attrs);
            }

            // 4. Resolve tag/component
            const tag = value.tag === 'body' ? 'div' : (tagRemapping[value.tag] ?? value.tag);

            // 5. Recursively render children
            const children = (value.content || []).map((child: JSONNode, idx: number) =>
                h(self, {
                    key: idx,
                    value: child,
                }),
            );

            // 6. Append external link icon
            if (value.tag === 'a' && !value.attrs?.href?.startsWith('/')) {
                children.push(
                    h(UIcon, {
                        class: 'ml-0.5 size-4',
                        name: 'i-mdi-open-in-new',
                    }),
                );
            }

            // 7. Return final tag/component with attributes
            return h(
                tag,
                {
                    id: value.id || undefined,
                    disabled: value.tag === 'input' ? true : undefined,
                    ...value.attrs,
                },
                children,
            );
        };
    },
});
</script>
