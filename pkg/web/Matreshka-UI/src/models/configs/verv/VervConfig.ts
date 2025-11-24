import { Component } from "vue";
import { Node, PatchConfigPatch } from "@vervstack/matreshka";

import VervConfigView from "@/components/config/verv/VervConfigView.vue";

import ConfigContent from "@/models/configs/ConfigContent.ts";
import AppInfoClass from "@/models/configs/verv/info/VervConfig.ts";
import DataSourceClass from "@/models/configs/verv/resources/Resource.ts";
import { mapDataSources } from "@/models/configs/verv/resources/Mapping.ts";
import { mapServer } from "@/models/configs/verv/servers/Mapping.ts";
import mapEnvVar from "@/models/configs/verv/env_vars/Mapping.ts";
import ServerClass from "@/models/configs/verv/servers/Servers.ts";
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";

export default class VervConfig implements ConfigContent {
  appInfo: AppInfoClass;
  dataSources: DataSourceClass[];
  servers: ServerClass[];
  envVars: EnvVar[] = [];

  constructor(root: Node) {
    let appInfo: AppInfoClass | undefined;
    let dataSources: DataSourceClass[] = [];
    let servers: ServerClass[] = [];

    root.innerNodes?.map((node: Node) => {
      switch (node.name) {
        case "APP-INFO":
          appInfo = new AppInfoClass(node);
          break;
        case "DATA-SOURCES":
          dataSources = mapDataSources(node);
          break;
        case "SERVERS":
          servers = mapServer(node);
          break
        case "ENVIRONMENT":
          this.envVars = mapEnvVar(node);
          break;
      }
    });

    if (!appInfo) {
      throw { message: "No app info found in env" };
    }

    this.appInfo = appInfo;
    this.dataSources = dataSources;
    this.servers = servers;
  }

  public isChanged(): boolean {
    return this.getChanges().length != 0;
  }

  public getChanges(): PatchConfigPatch[] {
    const changes: PatchConfigPatch[] = [];
    changes.push(...this.appInfo.getChanges());

    this.dataSources.map((ds) => changes.push(...ds.getChanges()));

    this.servers.map((s) => changes.push(...s.getChanges()));

    this.envVars.map((e) => {
      changes.push(...e.getChanges());
    });

    return changes;
  }

  public getChangedDataSourcesNames(): string[] {
    const changedDataSourceNames: string[] = [];
    this.dataSources.map((ds) => {
      if (ds.isChanged()) {
        changedDataSourceNames.push(ds.resourceName);
      }
    });

    return changedDataSourceNames;
  }

  public getChangedServersNames(): string[] {
    const changedServerNames: string[] = [];
    this.servers.map((serv) => {
      if (serv.isChanged()) {
        changedServerNames.push(serv.name);
      }
    });
    return changedServerNames;
  }

  public getChangedEnvVarNames(): string[] {
    const changedEnvVarNames: string[] = [];
    this.envVars.map((e: EnvVar) => {
      if (e.isChanged()) {
        changedEnvVarNames.push(e.getOriginalName());
      }
    });
    return changedEnvVarNames;
  }

  public rollback() {
    this.appInfo.rollback();
    this.dataSources.map((ds) => ds.rollback());
    this.servers.map((s) => s.rollback());
  }

  getComponent(): Component {
    return VervConfigView;
  }
}
