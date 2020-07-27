import { DataQuery, DataSourceJsonData } from "@grafana/data";

export enum TelegrafQueryType {
  Stream = "stream"
}

export interface TelegrafQuery extends DataQuery {
  queryType: TelegrafQueryType;
  measurement?: string;
  // filter?: string;
}

export const defaultQuery: Partial<TelegrafQuery> = {
  queryType: TelegrafQueryType.Stream
};

/**
 * These are options configured for each DataSource instance
 */
export interface TelegrafDataSourceOptions extends DataSourceJsonData {
  channel?: string; // The channel everything will be sent to
  buffer?: 2000;
}
