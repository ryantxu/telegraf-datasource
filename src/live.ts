import { Observable, Subject, ReplaySubject } from "rxjs";

import {
  DataQueryResponse,
  KeyValue,
  FieldType,
  CircularDataFrame,
  SelectableValue
} from "@grafana/data";
import { TelegrafDataSourceOptions } from "types";
import { getGrafanaLiveSrv } from "@grafana/runtime";

interface MeasureState {
  name: string;
  subject: Subject<DataQueryResponse>;
  frame: CircularDataFrame;
}

export interface MeasurementStreamer {
  getMeasurements(): Array<SelectableValue<string>>;
  getMeasurement(name: string): Observable<DataQueryResponse>;
}

class RecentMetrics implements MeasurementStreamer {
  cache: KeyValue<MeasureState> = {};

  constructor(private bufferSize: number, channel: string) {
    const srv = getGrafanaLiveSrv();
    if (!srv) {
      console.error('Grafana live not running, enable "live" feature toggle');
    }
    srv.initChannel(channel, {
      onPublish: (msg: any) => {
        return msg; // could pre-process the message
      }
    });
    srv.getChannelStream(channel).subscribe({
      next: (v: any) => {
        console.log("GOT", v);
      }
    });
  }

  getState(name: string): MeasureState {
    let v = this.cache[name];
    if (v) {
      return v;
    }

    const df = new CircularDataFrame({
      append: "tail",
      capacity: this.bufferSize
    });
    df.name = name;
    df.addField({ name: "timestamp", type: FieldType.time }, 0);
    this.cache[name] = v = {
      name,
      frame: df,
      subject: new ReplaySubject(1)
    };
    return v;
  }

  getMeasurement(name: string): Observable<DataQueryResponse> {
    return this.getState(name).subject.asObservable();
  }

  getMeasurements(): Array<SelectableValue<string>> {
    return Object.values(this.cache).map( v => {
      return {
        label: v.name,
        value: v.name,
        description: v.frame.fields.map( f => f.name ).join(', '),
      }
    });
  }
}

let singleton: RecentMetrics;

export function getMeasurementStreamer(
  cfg: TelegrafDataSourceOptions
): MeasurementStreamer {
  if (!singleton) {
    singleton = new RecentMetrics(
      cfg.buffer || 1000,
      cfg.channel || "telegraf"
    );
  }
  return singleton;
}
