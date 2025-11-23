import { Node } from "@vervstack/matreshka";

import { mapGrpc } from "@/models/configs/verv/resources/Grpc.ts";
import { newPostgres } from "@/models/configs/verv/resources/Postgres.ts";
import { mapRedis } from "@/models/configs/verv/resources/Redis.ts";
import DataSource from "@/models/configs/verv/resources/Resource.ts";
import { ResourceType } from "@/models/configs/verv/resources/ResourceTypes.ts";
import { mapSqlite } from "@/models/configs/verv/resources/Sqlite.ts";
import { mapTelegram } from "@/models/configs/verv/resources/Telegram.ts";
import { extractResourceType } from "@/models/shared/Values.ts";

const resourceMapping = new Map<string, (node: Node) => DataSource>();

resourceMapping.set(ResourceType.Postgres, newPostgres);
resourceMapping.set(ResourceType.Sqlite, mapSqlite);
resourceMapping.set(ResourceType.Redis, mapRedis);
resourceMapping.set(ResourceType.Telegram, mapTelegram);
resourceMapping.set(ResourceType.Grpc, mapGrpc);

export function mapDataSources(root: Node): DataSource[] {
  const dataSources: DataSource[] = [];
  root.innerNodes?.map((n) => {
    const resType = extractResourceType(n, root);

    if (!resType) {
      return;
    }

    const mapper = resourceMapping.get(resType);

    if (mapper) {
      dataSources.push(mapper(n));
    }
  });

  return dataSources;
}
