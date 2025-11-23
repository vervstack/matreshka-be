import { ConfigValue } from "@/models/shared/Values.ts";
import { PatchConfigPatch } from "@vervstack/matreshka";


export default class ServerClass {
  port: ConfigValue<number> = new ConfigValue<number>("", 0);
  name: string;
  grpc: GrpcHandler[] = [];
  fs: FsHandler[] = [];

  constructor(name: string) {
    this.name = name;
  }

  public isChanged(): boolean {
    let grpcChanged: boolean = false;
    this.grpc.map((s) => (grpcChanged = grpcChanged || s.isChanged()));

    let fsChanged: boolean = false;
    this.fs.map((s) => (fsChanged = fsChanged || s.isChanged()));
    return this.port.isChanged() || grpcChanged || fsChanged;
  }

  rollback(): void {
    this.port.rollback();
    this.grpc.map((g) => g.rollback());
    this.fs.map((f) => f.rollback());
  }

  getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];
    changes.push(...this.port.getChanges());

    this.grpc.map((g) => g.getChanges());
    this.fs.map((f) => f.getChanges());

    return changes;
  }
}

export class GrpcHandler {
  module: ConfigValue<string> = new ConfigValue("", "");
  gateway: ConfigValue<string> = new ConfigValue("", "");

  isChanged(): boolean {
    return this.module.isChanged() || this.gateway.isChanged();
  }

  rollback(): void {
    this.module.rollback();
    this.gateway.rollback();
  }

  getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];

    changes.push(...this.module.getChanges());
    changes.push(...this.gateway.getChanges());

    return changes;
  }
}

export class FsHandler {
  dist: ConfigValue<string> = new ConfigValue<string>("", "");

  isChanged(): boolean {
    return this.dist.isChanged();
  }

  rollback(): void {
    this.dist.rollback();
  }

  getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];

    changes.push(...this.dist.getChanges());

    return changes;
  }
}
