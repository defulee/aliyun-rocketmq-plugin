import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, InlineFieldRow, InlineField, Select } from '@grafana/ui';

import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { RocketMqQuery, SlsDataSourceOptions, Action } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, RocketMqQuery, SlsDataSourceOptions>;

const actionOptions = [
  { label: 'ConsumerAccumulate', value: Action.ConsumerAccumulate },
  { label: 'TrendTopicInputTps', value: Action.TrendTopicInputTps },
  { label: 'TrendGroupOutputTps', value: Action.TrendGroupOutputTps },
];

export class QueryEditor extends PureComponent<Props> {
  onActionChange = (event: SelectableValue<Action>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, action: event.value });
    onRunQuery();
  };

  onGroupIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, groupId: event.target.value });
  };

  onTopicChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, topic: event.target.value });
  };

  render() {
    const query = defaults(this.props.query);

    return (
      <>
        <div className="gf-form">
          <InlineFieldRow>
            <InlineField label="Action" grow>
              <Select options={actionOptions} onChange={this.onActionChange} value={query.action} />
            </InlineField>
          </InlineFieldRow>

          {(query.action === Action.ConsumerAccumulate || query.action === Action.TrendGroupOutputTps) && (
            <FormField
              labelWidth={8}
              value={query.groupId}
              onChange={this.onGroupIdChange}
              label="GroupId"
              tooltip="rocket mq consumer group id"
            />
          )}

          {(query.action === Action.TrendTopicInputTps || query.action === Action.TrendGroupOutputTps) && (
            <FormField
              labelWidth={8}
              value={query.topic}
              onChange={this.onTopicChange}
              label="Topic"
              tooltip="rocket mq topic."
            />
          )}
        </div>
      </>
    );
  }
}
