import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { ObservableV2 } from 'lib0/observable';
import { applyAwarenessUpdate, Awareness, encodeAwarenessUpdate, removeAwarenessStates } from 'y-protocols/awareness';
import * as Y from 'yjs';
import { ClientPacket, type ServerPacket } from '~~/gen/ts/resources/collab/collab';

export type StreamConnectFn = (options?: RpcOptions) => DuplexStreamingCall<ClientPacket, ServerPacket>;

const logger = useLogger('📞 yjs:grpc');

interface GrpcProviderOpts {
    /** Target id (e.g., document id, page id) */
    targetId: number;
}

type Events = {
    sync(synced: boolean, doc: Y.Doc): void;
    promoted(): void;
    saved(): void;
    loading(loading: boolean): void;
};

export default class GrpcProvider extends ObservableV2<Events> {
    public readonly ydoc: Y.Doc;
    public readonly awareness: Awareness;

    private readonly opts: GrpcProviderOpts;
    private clientId: number | undefined;
    private streamConnect: StreamConnectFn;
    private stream: DuplexStreamingCall<ClientPacket, ServerPacket> | undefined;
    private connected = false;
    private reconnectAttempt = 1;
    private authoritative = false;
    private synced = false;
    private destroyed = false;
    private connecting = false;
    private reconnectTimer: ReturnType<typeof setTimeout> | undefined;

    constructor(doc: Y.Doc, streamProvider: StreamConnectFn, opts: GrpcProviderOpts) {
        super();

        this.opts = opts;
        this.ydoc = doc;
        this.awareness = new Awareness(doc);

        this.streamConnect = streamProvider;

        // Setup local listeners
        this.ydoc.on('update', this.handleDocUpdate);
        this.awareness.on('update', this.handleAwarenessUpdate);
    }

    public get isAuthoritative() {
        return this.authoritative;
    }

    public get isSynced() {
        return this.synced;
    }

    // Public helpers
    override destroy() {
        if (this.destroyed) return;

        this.destroyed = true;
        // Clear local user state
        this.clientId && removeAwarenessStates(this.awareness, [this.clientId], 'app closed');

        setTimeout(() => {
            if (this.reconnectTimer !== undefined) {
                clearTimeout(this.reconnectTimer);
                this.reconnectTimer = undefined;
            }
            this.stream?.requests.complete();
            this.ydoc.off('update', this.handleDocUpdate);
            this.awareness.off('update', this.handleAwarenessUpdate);
        }, 0);

        super.destroy();
        logger.debug('Destroyed grpc provider');
    }

    // Internal
    public connect() {
        if (this.destroyed || this.connected || this.connecting) return;

        this.connecting = true;
        logger.info('Connecting to collab gRPC stream');

        let stream: DuplexStreamingCall<ClientPacket, ServerPacket>;
        try {
            stream = this.streamConnect({});
            this.stream = stream;
        } catch (err) {
            this.connecting = false;
            logger.error('Failed to connect to collab gRPC stream', err);
            this.scheduleReconnect();
            return;
        }

        stream.responses.onError((error) => {
            logger.warn('Collab gRPC stream ended', {
                error,
                connected: this.connected,
                synced: this.synced,
                authoritative: this.authoritative,
                clientId: this.clientId,
            });

            // A delayed callback from an older stream must not tear down a newer one.
            if (this.stream !== stream) return;

            this.connecting = false;
            this.connected = false;
            this.clientId = undefined;
            this.synced = false;
            this.authoritative = false;

            this.stream = undefined;
            this.scheduleReconnect();
        });

        stream.responses.onMessage((msg: ServerPacket) => {
            if (msg.msg.oneofKind === 'handshake' && !this.clientId) {
                logger.info('Received handshake message from server', msg.msg.handshake);
                this.clientId = msg.msg.handshake.clientId;
                this.connected = true;
                this.connecting = false;

                const sv = Y.encodeStateVector(this.ydoc);
                this.send(
                    ClientPacket.create({
                        msg: {
                            oneofKind: 'syncStep',
                            syncStep: {
                                step: 1,
                                data: sv,
                            },
                        },
                    }),
                );
                return;
            }

            if (!this.clientId) {
                logger.warn('Received message before clientId was set', msg);
                return;
            }

            // Ignore our own echoes (server broadcasts them to everyone incl. sender)
            if (msg.senderId === this.clientId) return;
            logger.debug('Received message from', msg.senderId, 'oneofKind:', msg.msg.oneofKind);

            switch (msg.msg.oneofKind) {
                case 'syncStep': {
                    if (msg.msg.syncStep.data.length === 0 || msg.msg.syncStep.step < 1 || msg.msg.syncStep.step > 2) {
                        logger.warn('Received invalid sync step', msg.msg.syncStep);
                        break;
                    }

                    logger.debug(
                        'Received sync step',
                        msg.msg.syncStep.step,
                        'from',
                        msg.senderId,
                        'length',
                        msg.msg.syncStep.data.length,
                    );
                    if (msg.msg.syncStep.step === 1) {
                        const diff = Y.encodeStateAsUpdate(this.ydoc, msg.msg.syncStep.data);

                        this.send({
                            msg: {
                                oneofKind: 'syncStep',
                                syncStep: {
                                    step: 2,
                                    data: diff,
                                    receiverId: msg.senderId,
                                },
                            },
                        });
                    } else if (msg.msg.syncStep.step === 2) {
                        Y.applyUpdate(this.ydoc, msg.msg.syncStep.data);

                        this.triggerSync();
                    }

                    break;
                }

                case 'awareness': {
                    if (msg.msg.awareness.data.length > 0) {
                        applyAwarenessUpdate(this.awareness, msg.msg.awareness.data, 'remote');
                    }
                    break;
                }
                case 'yjsUpdate': {
                    if (msg.msg.yjsUpdate.data.length > 0) {
                        logger.debug('Received Yjs update', msg.msg.yjsUpdate.data.length, 'from', msg.senderId);
                        Y.applyUpdate(this.ydoc, msg.msg.yjsUpdate.data);
                    }
                    break;
                }

                case 'targetSaved': {
                    logger.info('Received target saved message', msg.msg.targetSaved);

                    this.emit('saved', []);
                    break;
                }

                case 'promote': {
                    // Only act if we were *not* authoritative so far.
                    if (!this.authoritative) {
                        this.authoritative = true;
                        this.emit('promoted', []);
                    }

                    logger.info('Received promote: we are the new first client. Authoritative:', this.authoritative);

                    // Trigger step-0 again to seed the room and force sync
                    this.triggerSync();
                    break;
                }
            }
        });

        this.sendHello();

        logger.debug('Connect call completed, waiting for handshake');
    }

    private scheduleReconnect() {
        if (this.destroyed) return;
        if (this.reconnectTimer !== undefined) return;

        const delay = Math.min(this.reconnectAttempt * 750, 10_000);

        logger.info('Scheduling collab reconnect', {
            reconnectAttempt: this.reconnectAttempt,
            delay,
            destroyed: this.destroyed,
            connected: this.connected,
            synced: this.synced,
            authoritative: this.authoritative,
            clientId: this.clientId,
        });

        this.emit('sync', [false, this.ydoc]);
        this.emit('loading', [true]);

        if (delay >= 10_000) {
            logger.info('Max reconnect delay reached, resetting attempt counter');
            this.reconnectAttempt = 1;
        }

        this.reconnectAttempt++;
        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = undefined;
            this.connect();
        }, delay);
    }

    private triggerSync() {
        // If we were still waiting for sync, flip the flag and emit events
        if (!this.synced) {
            this.synced = true;
            this.reconnectAttempt = 1;

            logger.info('Provider sync emit');
            this.ydoc.emit('sync', [true, this.ydoc]);
            this.emit('sync', [true, this.ydoc]);
            logger.info('Post sync emit');
        }

        this.emit('loading', [false]);
    }

    // Yjs to Server
    private handleDocUpdate = (update: Uint8Array) => {
        if (!this.connected) return;

        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'yjsUpdate',
                yjsUpdate: {
                    data: update,
                },
            },
        });

        logger.debug('Send yjs update', update.length);
        this.send(msg);
    };

    private handleAwarenessUpdate = (
        { added, updated, removed }: { added: number[]; updated: number[]; removed: number[] },
        _origin: unknown,
    ) => {
        const changed = added.concat(updated, removed);
        if (changed.length === 0 || !this.connected) return;

        const update = encodeAwarenessUpdate(this.awareness, changed);
        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'awareness',
                awareness: {
                    data: update,
                },
            },
        });

        logger.debug('Send awareness update', update.length);
        this.send(msg);
    };

    private sendHello() {
        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'hello',
                hello: {
                    targetId: this.opts.targetId,
                },
            },
        });

        logger.debug('Send hello message', this.opts.targetId);
        this.send(msg);
    }

    private async send(msg: ClientPacket) {
        try {
            await this.stream?.requests.send(msg);
        } catch (_) {
            // swallow if stream closed mid-send
        }
    }
}
