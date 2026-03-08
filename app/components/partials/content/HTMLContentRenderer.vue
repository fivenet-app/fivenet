<script lang="ts">
import { UCheckbox, UIcon } from '#components';
import { defineComponent, getCurrentInstance, h, Text, type Component, type VNode } from 'vue';
import type { RichTextHtmlNode } from '~~/gen/ts/resources/common/content/content';
import GenericImg from '../elements/GenericImg.vue';

export default defineComponent({
    name: 'HTMLContentRenderer',

    props: {
        value: {
            type: Object as () => RichTextHtmlNode,
            required: true,
        },
    },

    setup(props) {
        // Self-reference for recursion
        const self = getCurrentInstance()?.type as Component;

        // Optional tag remapping
        const tagRemapping: Record<string, Component> = {};

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

            // 4. img tag
            if (value.tag === 'img') {
                return h(GenericImg, {
                    ...value.attrs,
                    src: cleanupImageURL(value.attrs?.src || ''),
                });
            }

            // 4.1 Penalty calculator embed marker (template preview fallback)
            if (
                value.tag === 'div' &&
                (value.attrs?.['data-embed'] === 'penalty-calculator' ||
                    value.attrs?.['data-type'] === 'penalty-calculator' ||
                    value.attrs?.['data-type'] === 'penaltyCalculator')
            ) {
                return h('div', { class: 'my-2 rounded-md border border-dashed border-neutral-300 p-2 text-sm text-muted' }, [
                    'Penalty Calculator',
                ]);
            }

            // 5. Tag remapping
            if (tagRemapping[value.tag]) {
                return h(tagRemapping[value.tag]!, value.attrs);
            }

            const tag = value.tag === 'body' ? 'div' : value.tag;

            // 6. Recursively render children
            const children = (value.content || []).map((child: RichTextHtmlNode, idx: number) =>
                h(self, {
                    key: idx,
                    value: child,
                }),
            );

            // 7. Append external link icon
            if (value.tag === 'a' && !value.attrs?.href?.startsWith('/')) {
                children.push(
                    h(UIcon, {
                        class: 'ml-0.5 size-4',
                        name: 'i-mdi-open-in-new',
                    }),
                );
            }

            // 8. Return final tag/component with attributes
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
