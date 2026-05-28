import { defineAsyncComponent, markRaw, type DefineComponent } from 'vue';

export type IconEntry = {
    name: string;
    label: string;
};

export type IconComponent = ReturnType<typeof defineAsyncComponent<DefineComponent>>;

type IconModule = Record<string, DefineComponent> & {
    default?: DefineComponent;
};

type IconLoader = () => Promise<IconModule>;

const modules = import.meta.glob('../../../node_modules/mdi-vue3/icons/*.js', { eager: false }) as Record<string, IconLoader>;
const iconLoadersByName = new Map<string, IconLoader>();
const iconComponentsByName = new Map<string, IconComponent>();

const toLabel = (name: string): string => name.replace(/([a-z])([A-Z])/g, '$1 $2').replace(/Icon$/, '');

export const availableIcons: IconEntry[] = Object.entries(modules).map(([path, loader]) => {
    const name = path.split('/').pop()!.replace(/\.js$/, '');
    const entry: IconEntry = {
        name: name,
        label: toLabel(name),
    };

    iconLoadersByName.set(name, loader);

    return entry;
});

export const fallbackIconName = 'MapMarkerQuestionIcon';

function createIconComponent(name: string, loader: IconLoader): IconComponent {
    return markRaw(
        defineAsyncComponent<DefineComponent>({
            loader: async () => {
                const mod = await loader();
                // Grab default if present, otherwise the named export.
                const comp = mod.default || mod[name];
                if (!comp) {
                    throw new Error(`mdi-vue3: icon "${name}" not found in module`);
                }
                return comp;
            },
        }),
    ) as IconComponent;
}

function resolveOrCreateIconComponent(name: string): IconComponent | undefined {
    const cached = iconComponentsByName.get(name);
    if (cached) {
        return cached;
    }

    const loader = iconLoadersByName.get(name);
    if (!loader) {
        return undefined;
    }

    const component = createIconComponent(name, loader);
    iconComponentsByName.set(name, component);
    return component;
}

export function resolveIconComponent(iconName?: string | null): IconComponent {
    if (iconName) {
        const component = resolveOrCreateIconComponent(iconName);
        if (component) {
            return component;
        }
    }

    const fallbackComponent = resolveOrCreateIconComponent(fallbackIconName);
    if (!fallbackComponent) {
        throw new Error(`mdi-vue3: fallback icon "${fallbackIconName}" not found`);
    }

    return fallbackComponent;
}
export const fallbackIconComponent = resolveIconComponent(fallbackIconName);
