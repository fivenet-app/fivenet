declare module '#app' {
    interface PageMeta {
        title?: string;
        requiresAuth?: boolean;
        permission?: String;
        authOnlyToken?: boolean;
        showQuickButtons?: boolean;
    }
}

// It is always important to ensure you import/export something when augmenting a type
export {};

//Taken from https://github.com/Pictogrammers/vue-icon/issues/10#issuecomment-1528951622
declare let SvgIcon: import('vue').DefineComponent<{
    type: {
        type: StringConstructor;
        default: string;
    };
    path: {
        type: StringConstructor;
        default: string;
    };
    size: {
        type: NumberConstructor;
        optional: boolean;
    };
    viewbox: {
        type: StringConstructor;
        optional: boolean;
    };
    flip: {
        type: StringConstructor;
        optional: boolean;
    };
    rotate: {
        type: StringConstructor;
        optional: boolean;
    };
}>;

declare module '@jamescoyle/vue-icon' {
    export default SvgIcon;
}
