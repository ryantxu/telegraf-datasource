import { merge, Observable } from "rxjs";

import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings
} from "@grafana/data";

import {
  TelegrafQuery,
  TelegrafDataSourceOptions,
  TelegrafQueryType
} from "./types";
import { getMeasurementStreamer } from "live";

export class DataSource extends DataSourceApi<
  TelegrafQuery,
  TelegrafDataSourceOptions
> {
  constructor(
    private instanceSettings: DataSourceInstanceSettings<
      TelegrafDataSourceOptions
    >
  ) {
    super(instanceSettings);
  }

  getStreamer() {
    return getMeasurementStreamer(this.instanceSettings.jsonData);
  }
  
  query(
    options: DataQueryRequest<TelegrafQuery>
  ): Observable<DataQueryResponse> {
    const info = this.getStreamer();
    const res: Array<Observable<DataQueryResponse>> = [];

    // Return a constant for each query.
    for (const target of options.targets) {
      if (target.queryType === TelegrafQueryType.Stream) {
        const m = target.measurement;
        if (m) {
          res.push(info.getMeasurement(m));
        }
      } else {
        console.log("Unknown query:", target);
      }
    }

    if (res.length === 1) {
      return res[0];
    }
    return merge(...res);
  }

  async testDatasource() {
    return {
      status: "success",
      message: "Success"
    };
  }
}
