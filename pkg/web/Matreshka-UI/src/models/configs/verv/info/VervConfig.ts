import { Node, PatchConfigPatch } from "@vervstack/matreshka";

import { ConfigValue, extractStringValue } from "@/models/shared/Values.ts";

export default class AppInfoClass {
  name: ConfigValue<string>;
  serviceVersion: ConfigValue<string>;
  public updated_at?: Date;

  versions: string[] = [];

  constructor(root: Node) {
    let appName: ConfigValue<string> | undefined;
    let appVersion: ConfigValue<string> | undefined;

    root.innerNodes?.map((n) => {
      if (!n.name || !root.name) {
        return;
      }

      const name = n.name.slice(root.name.length + 1);
      switch (name) {
        case "NAME":
          const name = extractStringValue(n);
          appName = new ConfigValue<string>(name.envName, name.value);
          break;
        case "VERSION":
          const version = extractStringValue(n);
          appVersion = new ConfigValue<string>(version.envName, version.value);
          break;
      }

      return;
    });

    if (!appName) {
      throw { message: "no app name provided" };
    }

    if (!appVersion) {
      throw { message: "no app version provided" };
    }

    this.name = appName;
    this.serviceVersion = appVersion;
  }

  getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];
    changes.push(...this.name.getChanges());
    changes.push(...this.serviceVersion.getChanges());
    return changes;
  }

  isChanged(): boolean {
    return this.getChanges().length != 0;
  }

  rollback() {
    this.name.rollback();
    this.serviceVersion.rollback();
  }
}

export enum SourceCodeSystem {
  unknown = 0,
  github = 1,
}

export function ExtractSourceCodeSystemFromServiceName(name: string): SourceCodeSystem | undefined {
  if (name.includes("github")) {
    return SourceCodeSystem.github;
  }
}

const scsToPiIcon = new Map<SourceCodeSystem, string>();
scsToPiIcon.set(SourceCodeSystem.github, "pi-github");

export function PiIconFromSourceCodeSystem(
  sourceCodeSystem: SourceCodeSystem | undefined
): string | undefined {
  if (!sourceCodeSystem) {
    return;
  }

  return scsToPiIcon.get(sourceCodeSystem);
}
