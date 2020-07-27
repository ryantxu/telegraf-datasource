import defaults from "lodash/defaults";

import React, { PureComponent, ChangeEvent } from "react";
import { QueryEditorProps } from "@grafana/data";
import { DataSource } from "./DataSource";
import {
  TelegrafQuery,
  TelegrafDataSourceOptions,
  defaultQuery
} from "./types";
import { LegacyForms, Select } from "@grafana/ui";
import { getMeasurementStreamer } from "live";

type Props = QueryEditorProps<
  DataSource,
  TelegrafQuery,
  TelegrafDataSourceOptions
>;

interface State {}

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  onMeasurementChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, measurement: event.target.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { measurement } = query;

    const info = this.props.datasource.getStreamer();
    const measurments = info.getMeasurements();
    let current = measurments.find( m => m.value === measurement );
    if(!current) {
      current = {
        label: measurement,
        value: measurement,
        description: 'Unknown measurment',
      };
      measurments.push(current);
    }

    return (
      <div className="gf-form">
        <Select
          value={current}
          options={measurments}
          onChange={this.onMeasurementChange} />
      </div>
    );
  }
}
