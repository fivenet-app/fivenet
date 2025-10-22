<script setup lang="ts">
import EditorCanvas from './EditorCanvas.vue';
import EditorSidebar from './EditorSidebar.vue';
import EditorToolbar from './EditorToolbar.vue';
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar class="print:hidden" :title="$t('pages.documents.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar class="p-1 print:hidden">
                <template #default>
                    <EditorToolbar />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <!-- Container: full-screen flex layout with top toolbar and content area -->
            <div class="flex h-full max-w-screen flex-col">
                <!-- Main content: canvas and sidebar -->
                <div class="flex flex-1 overflow-hidden">
                    <div ref="canvasContainer" class="flex-1 overflow-hidden">
                        <!-- Canvas area fills remaining space -->
                        <EditorCanvas />
                    </div>

                    <!-- Sidebar on the right with fixed width -->
                    <EditorSidebar class="w-full max-w-sm min-w-64 shrink-0 border-l border-l-default bg-default" />
                </div>
            </div>
        </template>
    </UDashboardPanel>
</template>
