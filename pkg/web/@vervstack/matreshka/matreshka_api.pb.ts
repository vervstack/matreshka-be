/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum ConfigType {
  plain = "plain",
  verv = "verv",
  minio = "minio",
  pg = "pg",
  nginx = "nginx",
  kv = "kv",
}

export enum Format {
  yaml = "yaml",
  env = "env",
}

export enum SortType {
  default = "default",
  by_name = "by_name",
  by_updated_at = "by_updated_at",
}

export type Config = {
  id?: number;
  name?: string;
  createdAtUtcTimestamp?: string;
  updatedAtUtcTimestamp?: string;
  versions?: string[];
};

export type Paging = {
  limit?: string;
  offset?: string;
};

export type ApiVersionRequest = Record<string, never>;

export type ApiVersionResponse = {
  version?: string;
};

export type ApiVersion = Record<string, never>;

export type GetConfigRequest = {
  configName?: string;
  version?: string;
  format?: Format;
};

export type GetConfigResponse = {
  config?: Uint8Array;
};

export type GetConfig = Record<string, never>;

export type PatchConfigRequest = {
  configName?: string;
  version?: string;
  patches?: PatchConfigPatch[];
};

export type PatchConfigResponse = Record<string, never>;

type BasePatchConfigPatch = {
  fieldName?: string;
};

export type PatchConfigPatch = BasePatchConfigPatch &
  OneOf<{
    rename: string;
    updateValue: string;
    delete: boolean;
  }>;

export type PatchConfig = Record<string, never>;

export type StoreConfigRequest = {
  format?: Format;
  configName?: string;
  version?: string;
  config?: Uint8Array;
};

export type StoreConfigResponse = Record<string, never>;

export type StoreConfig = Record<string, never>;

export type ListConfigsRequest = {
  paging?: Paging;
  searchPattern?: string;
  sort?: Sort;
};

export type ListConfigsResponse = {
  configs?: Config[];
  totalRecords?: string;
};

export type ListConfigs = Record<string, never>;

export type Node = {
  name?: string;
  value?: string;
  innerNodes?: Node[];
};

export type GetConfigNodeRequest = {
  configName?: string;
  version?: string;
};

export type GetConfigNodeResponse = {
  root?: Node;
  versions?: string[];
};

export type GetConfigNode = Record<string, never>;

export type CreateConfigRequest = {
  configName?: string;
  configType?: ConfigType;
};

export type CreateConfigResponse = Record<string, never>;

export type CreateConfig = Record<string, never>;

export type RenameConfigRequest = {
  configName?: string;
  newName?: string;
};

export type RenameConfigResponse = {
  newName?: string;
};

export type RenameConfig = Record<string, never>;

export type Sort = {
  type?: SortType;
  desc?: boolean;
};

export type SubscribeOnChangesRequest = {
  subscribeConfigNames?: string[];
  unsubscribeConfigNames?: string[];
};

export type SubscribeOnChangesResponse = {
  configName?: string;
  timestamp?: number;
  patches?: PatchConfigPatch[];
};

export type SubscribeOnChanges = Record<string, never>;

export type DeleteConfigRequest = {
  configName?: string;
  configVersion?: string;
};

export type DeleteConfigResponse = Record<string, never>;

export type DeleteConfig = Record<string, never>;

export class MatreshkaApi {
  static ApiVersion(this:void, req: ApiVersionRequest, initReq?: fm.InitReq): Promise<ApiVersionResponse> {
    return fm.fetchRequest<ApiVersionResponse>(`/api/version?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
  static ListConfigs(this:void, req: ListConfigsRequest, initReq?: fm.InitReq): Promise<ListConfigsResponse> {
    return fm.fetchRequest<ListConfigsResponse>(`/api/config/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateConfig(this:void, req: CreateConfigRequest, initReq?: fm.InitReq): Promise<CreateConfigResponse> {
    return fm.fetchRequest<CreateConfigResponse>(`/api/config/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetConfig(this:void, req: GetConfigRequest, initReq?: fm.InitReq): Promise<GetConfigResponse> {
    return fm.fetchRequest<GetConfigResponse>(`/api/config/${req.configName}?${fm.renderURLSearchParams(req, ["configName"])}`, {...initReq, method: "GET"});
  }
  static GetConfigNodes(this:void, req: GetConfigNodeRequest, initReq?: fm.InitReq): Promise<GetConfigNodeResponse> {
    return fm.fetchRequest<GetConfigNodeResponse>(`/api/config/nodes`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static PatchConfig(this:void, req: PatchConfigRequest, initReq?: fm.InitReq): Promise<PatchConfigResponse> {
    return fm.fetchRequest<PatchConfigResponse>(`/api/config/${req.configName}/patch`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static StoreConfig(this:void, req: StoreConfigRequest, initReq?: fm.InitReq): Promise<StoreConfigResponse> {
    return fm.fetchRequest<StoreConfigResponse>(`/api/config/${req.configName}/store`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RenameConfig(this:void, req: RenameConfigRequest, initReq?: fm.InitReq): Promise<RenameConfigResponse> {
    return fm.fetchRequest<RenameConfigResponse>(`/api/config/${req.configName}/rename/${req.newName}`, {...initReq, method: "POST"});
  }
  static DeleteConfig(this:void, req: DeleteConfigRequest, initReq?: fm.InitReq): Promise<DeleteConfigResponse> {
    return fm.fetchRequest<DeleteConfigResponse>(`/api/config/${req.configName}/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}