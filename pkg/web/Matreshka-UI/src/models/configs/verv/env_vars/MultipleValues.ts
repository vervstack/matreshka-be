import { Component } from "vue";

import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";

import EnvVariableView from "@/components/config/verv/env/EnvVariableView.vue";
import { ConfigValue } from "@/models/shared/Values.ts";
import { PatchConfigPatch } from "@vervstack/matreshka";

export class StringArrayEnvVarClass implements EnvVar {
  readonly rootName: string;
  readonly values: ConfigValue<string>[];

  enums: string[] = [];

  constructor(envName: string, value: string[]) {
    this.rootName = envName;
    this.values = value.map((v, idx) => new ConfigValue(envName + `_[${idx}]`, v));
  }

  isChanged(): boolean {
    return this.values.find((v) => v.isChanged()) != undefined;
  }

  getOriginalName(): string {
    return this.rootName;
  }

  getComponent(): Component {
    return EnvVariableView;
  }

  getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];
    this.values.forEach((v) => {
      if (v.isChanged()) {
        changes.push(...v.getChanges());
      }
    });
    return changes;
  }

  setEnums(_: string[]) {
  }
}

