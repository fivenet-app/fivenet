import * as Y from 'yjs';
import { collabDrivers, type CollabCategory } from './yjs/collab-drivers';
import GrpcProvider from './yjs/yjs';

export async function useCollabDoc(category: CollabCategory, targetId: number) {
    const ydoc = new Y.Doc();

    const driver = collabDrivers[category];
    const driverFn = await driver();
    const provider = new GrpcProvider(ydoc, driverFn, {
        targetId: targetId,
    });

    onBeforeUnmount(() => provider.destroy());

    return { ydoc, provider };
}
