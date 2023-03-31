declare module '#app' {
    interface PageMeta {
        requiresAuth?: boolean;
        permission?: String;
        authOnlyToken?: boolean;
    }
}
// It is always important to ensure you import/export something when augmenting a type
export {};
