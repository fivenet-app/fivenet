export default defineNuxtPlugin(() => {
    const { $update } = useNuxtApp();
    $update.on('update', (version) => {
        if (confirm(`New FiveNet version ${version} available. Update?`)) {
            location.reload();
        }
    });
});
