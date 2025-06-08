import { MapMarkerQuestionIcon } from 'mdi-vue3';
import { defineAsyncComponent, type DefineComponent } from 'vue';

export type IconEntry = {
    name: string;
    component: ReturnType<typeof defineAsyncComponent>;
};

export const fallbackIcon = MapMarkerQuestionIcon;

const modules = import.meta.glob('../../../node_modules/mdi-vue3/icons/*.js', { eager: false });

export const availableIcons: IconEntry[] = Object.entries(modules).map(([path, loader]) => {
    const name = path.split('/').pop()!.replace(/\.js$/, '');

    return {
        name,
        component: defineAsyncComponent<DefineComponent>({
            loader: async () => {
                const mod = (await loader()) as Record<string, DefineComponent>;
                // Grab default if present, otherwise the named export
                const comp = mod.default || mod[name];
                if (!comp) {
                    throw new Error(`mdi-vue3: icon "${name}" not found in module`);
                }
                return comp;
            },
        }),
    };
});
