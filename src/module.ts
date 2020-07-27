import { DataSourcePlugin } from "@grafana/data";
import { DataSource } from "./DataSource";
import { ConfigEditor } from "./ConfigEditor";
import { QueryEditor } from "./QueryEditor";
import { TelegrafQuery, TelegrafDataSourceOptions } from "./types";

export const plugin = new DataSourcePlugin<
  DataSource,
  TelegrafQuery,
  TelegrafDataSourceOptions
>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
