export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();
    $update.on('update', (version) => {
        // TODO add translations
        if (confirm(`New FiveNet version ${version} available. Reload?`)) {
            location.reload();
        }
    });
});
