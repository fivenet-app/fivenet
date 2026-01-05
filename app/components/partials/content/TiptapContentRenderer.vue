<script lang="ts">
import type { Extensions, JSONContent } from '@tiptap/core';
import { computed, defineComponent, h, type PropType, type VNode } from 'vue';
import { renderToVueVNode, type VueStaticRendererOptions } from '~/utils/tiptap';

export default defineComponent({
    name: 'TiptapContentRenderer',

    props: {
        value: { type: Object as PropType<JSONContent | null>, default: null },
        extensions: { type: Array as PropType<Extensions>, required: true },
        options: { type: Object as PropType<VueStaticRendererOptions>, default: undefined },
        as: { type: String, default: 'div' },
    },

    setup(props) {
        const vnode = computed(() => {
            if (!props.value) return null;

            return renderToVueVNode({
                extensions: props.extensions,
                content: props.value,
                options: props.options,
            });
        });

        /* eslint-disable @typescript-eslint/no-explicit-any */
        return (): VNode => h(props.as, null, vnode.value as any);
    },
});
</script>
