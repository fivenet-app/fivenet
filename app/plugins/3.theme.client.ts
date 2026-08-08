import { setTabletColors } from '~/composables/nui';
import { useSettingsStore } from '~/stores/settings';

const paletteShades = [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950] as const;

function setPaletteVars(root: HTMLElement, uiColor: 'primary' | 'neutral', colorName: string): void {
    for (const shade of paletteShades) {
        root.style.setProperty(`--ui-color-${uiColor}-${shade}`, `var(--color-${colorName}-${shade})`);
        root.style.setProperty(`--ui-color-${uiColor}-${shade}-rgb`, `var(--color-${colorName}-${shade}-rgb)`);
    }

    root.style.setProperty(`--ui-color-${uiColor}`, `var(--color-${colorName}-500)`);
    root.style.setProperty(`--ui-color-${uiColor}-rgb`, `var(--color-${colorName}-500-rgb)`);
}

export default defineNuxtPlugin({
    name: 'theme-sync',
    dependsOn: ['config', 'nui'],
    parallel: true,

    setup() {
        const appConfig = useAppConfig();
        const colorMode = useColorMode();
        const settingsStore = useSettingsStore();
        const { design } = storeToRefs(settingsStore);
        const themeColor = computed(() => (colorMode.value === 'dark' ? '#111827' : '#fff'));

        useHead({
            meta: [{ key: 'theme-color', name: 'theme-color', content: themeColor }],
        });

        const applyThemeColors = (): void => {
            const primary = design.value.ui.primary;
            const gray = design.value.ui.gray;

            appConfig.ui.colors.primary = primary;
            appConfig.ui.colors.neutral = gray;
            void setTabletColors(appConfig.ui.colors.primary, appConfig.ui.colors.neutral);

            // Only set CSS variables for Chrome 103+ due to lack of support in earlier versions.
            const root = document.documentElement;
            if (!root.classList.contains('polyfills')) return;

            setPaletteVars(root, 'primary', primary);
            setPaletteVars(root, 'neutral', gray);
        };

        watch(design, applyThemeColors, { deep: true, immediate: true });
    },
});
