import { Component } from "vue";
import { PatchConfigPatch } from "@vervstack/matreshka";

export interface EnvVar {
  getOriginalName(): string;

  isChanged(): boolean;

  getComponent(): Component;

  getChanges(): PatchConfigPatch[]
}

export enum DataType {
  STRING = "string",
  INT = "int",
  FLOAT = "float",
  BOOLEAN = "bool",
}
