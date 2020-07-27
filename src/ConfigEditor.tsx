import React, { PureComponent } from "react";

import { DataSourcePluginOptionsEditorProps } from "@grafana/data";
import { TelegrafDataSourceOptions } from "./types";
import { LegacyForms } from "@grafana/ui";

interface Props
  extends DataSourcePluginOptionsEditorProps<TelegrafDataSourceOptions> {}

export class ConfigEditor extends PureComponent<Props> {
  onComponentDidMount() {}

  onFormChange = (optionValue: string, optionKey: string) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        [optionKey]: optionValue
      }
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData } = options;

    return (
      <div>
        <p>
          This is a work-in-progress while grafana core improves the core
          streaming support. In the current form, the telegraph plugin must
          connect to a running grafana instance over websockets and push each
          message
        </p>

        <div className="gf-form">
          <LegacyForms.FormField
            label="Channel"
            labelWidth={13}
            inputWidth={24}
            tooltip={"this should match the channel configured in telegraf"}
            onChange={e => this.onFormChange(e.target.value, "channel")}
            value={jsonData.channel}
            placeholder="The channel name"
          />
        </div>
      </div>
    );
  }
}
