import defaults from "lodash/defaults";

import React, { PureComponent } from "react";
import { QueryEditorProps, SelectableValue } from "@grafana/data";
import { DataSource } from "./DataSource";
import {
  TelegrafQuery,
  TelegrafDataSourceOptions,
  defaultQuery
} from "./types";
import { Select } from "@grafana/ui";

type Props = QueryEditorProps<
  DataSource,
  TelegrafQuery,
  TelegrafDataSourceOptions
>;

interface State {}

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  onMeasurementChange = (v: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, measurement: v.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { measurement } = query;

    const info = this.props.datasource.getStreamer();
    const measurments = info.getMeasurements();
    let current = measurments.find(m => m.value === measurement);
    if (!current && measurement) {
      current = {
        label: measurement,
        value: measurement,
        description: "Unknown measurment"
      };
      measurments.push(current);
    }

    return (
      <div className="gf-form">
        <Select
          value={current}
          options={measurments}
          onChange={this.onMeasurementChange}
          placeholder="Select measurment"
        />
      </div>
    );
  }
}
