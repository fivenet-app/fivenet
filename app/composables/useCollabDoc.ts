import * as Y from 'yjs';
import { collabDrivers, type CollabCategory } from './yjs/collab-drivers';
import GrpcProvider from './yjs/yjs';

// useCollabDoc provides a Y.Doc and GrpcProvider for the given category and targetId.
// Remember to destroy the provider when done.
export async function useCollabDoc(category: CollabCategory, targetId: number) {
    const ydoc = new Y.Doc();

    const driver = collabDrivers[category];
    const driverFn = await driver();
    const provider = new GrpcProvider(ydoc, driverFn, {
        targetId: targetId,
    });

    return { ydoc, provider };
}
