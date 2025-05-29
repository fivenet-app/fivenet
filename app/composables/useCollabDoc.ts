import * as Y from 'yjs';
import { collabDrivers, type CollabCategory } from './yjs/collab-drivers';
import GrpcProvider from './yjs/yjs';

export function useCollabDoc(category: CollabCategory, targetId: number) {
    const ydoc = new Y.Doc();

    const driver = collabDrivers[category];
    const provider = new GrpcProvider(ydoc, driver, {
        targetId: targetId,
    });

    return { ydoc, provider };
}
