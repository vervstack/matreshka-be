import { Component } from "vue";

export interface EnvVar {
  getOriginalName(): string;

  isChanged(): boolean;

  getComponent(): Component;
}

export enum DataType {
  STRING = "string",
  INT = "int",
  FLOAT = "float",
  BOOLEAN = "bool",
}
