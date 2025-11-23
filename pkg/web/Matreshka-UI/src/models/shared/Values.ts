import { Node, PatchConfigPatch } from "@vervstack/matreshka";
import { newDeletePatch, newRenamePatch, newUpdatePatch } from "@/models/shared/Patch.ts";

export type KeyMap = {
  [key: string]: any;
};

export class ConfigValue<T extends { toString(): string } | boolean> {
  envName: string;
  value: T;

  isMuted: boolean = false;
  isNew: boolean = false;

  enums: T[] = [];

  private readonly originalName: string;
  private readonly originalValue: T;

  constructor(envName: string, value: T) {
    this.originalName = envName;
    this.originalValue = value;

    this.envName = envName;
    this.value = value;
  }

  getOriginalValue(): T {
    return this.originalValue;
  }

  getOriginalName(): string {
    return this.originalName;
  }

  isChanged(): boolean {
    if (this.isMuted) {
      return false;
    }

    return this.value != this.originalValue ||
      this.envName != this.originalName;
  }

  isNameChanged(): boolean {
    if (this.isMuted) {
      return false;
    }

    return this.envName != this.originalName;
  }

  isValueChanged(): boolean {
    if (this.isMuted) {
      return false;
    }

    return this.value != this.originalValue;
  }

  getChanges(): PatchConfigPatch[] {
    if (this.isMuted) {
      return [];
    }

    const changes: PatchConfigPatch[] = [];

    if (this.isNameChanged()) {
      if (this.isValueChanged()) {
        // Changed both name and value - practically new variable
        changes.push(newDeletePatch(this.originalName));
        changes.push(newUpdatePatch(this.envName, this.value.toString()));
      } else {
        // Only name changed
        changes.push(newRenamePatch(this.originalName, this.envName));
      }
    } else if (this.isValueChanged() || this.isNew) {
      // Value changed or it's a new variable
      changes.push(newUpdatePatch(this.envName, this.value.toString()));
    }

    return changes;
  }

  rollback() {
    this.value = this.originalValue;
    this.envName = this.originalName;
  }
}

export function extractStringValue(n: Node): ConfigValue<string> {
  return new ConfigValue<string>(n.name || "", n.value || "");
}

export function extractNumberValue(n: Node): ConfigValue<number> {
  return new ConfigValue(n.name || "", Number(n.value) || 0);
}

export function extractResourceType(node: Node, root: Node): string | undefined {
  if (!node.name || !root.name) {
    return;
  }

  let name = node.name.slice(root.name.length + 1);

  const resourceNameEndIdx = name.indexOf("_");
  if (resourceNameEndIdx > 0) {
    name = name.slice(resourceNameEndIdx);
  }
  name = name.toLowerCase();

  let resType = name;
  const resourceTypeNameEndIdx = resType.indexOf("-");
  if (resourceTypeNameEndIdx > 0) {
    resType = name.slice(0, resourceTypeNameEndIdx);
  }

  return resType;
}
