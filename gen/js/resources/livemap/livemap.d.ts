import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';

export interface IMarker {
    getX(): number;

    getY(): number;

    getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;

    getId(): number;

    getName(): string;

    getIcon(): string;

    getIconColor(): string;

    toObject(includeInstance?: boolean): IMarker.AsObject;
    static toObject(includeInstance: boolean, msg: DispatchMarker): IMarker.AsObject;
}

export namespace IMarker {
  export type AsObject = {
    x: number,
    y: number,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    id: number,
    name: string,
    icon: string,
    iconColor: string,
  }
}
