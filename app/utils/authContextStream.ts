export type AuthContextStreamState = {
    abort?: AbortController;
    isStopping: () => boolean;
};

/** Restart an active stream after a committed authorization-context change. */
export async function restartForAuthContext(
    state: AuthContextStreamState,
    stopStream: () => Promise<void>,
    startStream: () => Promise<void>,
): Promise<void> {
    if (!state.abort || state.isStopping()) return;

    await stopStream();
    if (!state.isStopping()) await startStream();
}

/** Create the callback for the committed userState/transition watcher. */
export function createAuthContextRestartHandler(
    initialKey: string,
    restart: () => void,
): (state: readonly [key: string, transitioning: boolean]) => void {
    let lastRestartedKey = initialKey;

    return ([key, transitioning]) => {
        if (transitioning || key === lastRestartedKey) return;

        lastRestartedKey = key;
        restart();
    };
}
