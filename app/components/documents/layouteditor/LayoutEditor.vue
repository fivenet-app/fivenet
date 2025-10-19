<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import Editor from 'svgedit/dist/editor/Editor.js';
import 'svgedit/dist/editor/svgedit.css';

const container = useTemplateRef('container');

onMounted(() => {
    const svgEditor = new Editor(container.value);
    svgEditor.setConfig({
        allowInitialUserOverride: true,
        extensions: [],
        noDefaultExtensions: false,
        userExtensions: [],
    });
    svgEditor.init();

    useTimeoutFn(() => {
        /* global svgEditor */
        svgEditor.addExtension('placeholders', function () {
            const svgCanvas = svgEditor.canvas;

            function insertPlaceholder(opts) {
                const { key, type = 'short_text', required = true, sample = `{{${key}}}` } = opts;
                const x = 50,
                    y = 50;

                // Create a group to keep things together
                const g = svgCanvas.addSVGElementFromJson({
                    element: 'g',
                    attr: { class: 'ph', 'data-ph-key': key, 'data-ph-type': type, 'data-ph-required': String(required) },
                });

                // Text token
                const t = svgCanvas.addSVGElementFromJson({
                    element: 'text',
                    attr: { x, y, 'xml:space': 'preserve' },
                });
                t.textContent = sample;
                g.append(t);

                // Underline for visual affordance
                const line = svgCanvas.addSVGElementFromJson({
                    element: 'line',
                    attr: { x1: x, y1: y + 6, x2: x + 140, y2: y + 6, 'stroke-dasharray': '4,3' },
                });
                g.append(line);

                svgCanvas.selectOnly([g]);
            }

            return {
                name: 'Text Placeholders',
                svgicons: null,
                buttons: [
                    {
                        id: 'placeholders_add',
                        type: 'mode',
                        title: 'Insert Placeholder',
                        events: {
                            click() {
                                // You can show a custom dialog; for demo, use a prompt
                                const key = window.prompt('Field key (e.g., patient_name):', 'patient_name');
                                if (key) insertPlaceholder({ key });
                            },
                        },
                    },
                ],
                // Optional: context menu entries, elementChanged hooks, etc.
            };
        });
    }, 1000);
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <!-- Top Toolbar -->
            <UDashboardNavbar title="Layout Editor">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div ref="container" style="width: 100%; height: 95vh"></div>
        </template>
    </UDashboardPanel>
</template>
