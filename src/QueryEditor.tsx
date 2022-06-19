import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';

import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { RocketMqQuery, SlsDataSourceOptions } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, RocketMqQuery, SlsDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onGroupIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, groupId: event.target.value });
  };

  render() {
    const query = defaults(this.props.query);

    return (
      <>
        <div className="gf-form">
          <FormField
            labelWidth={8}
            value={query.groupId}
            onChange={this.onGroupIdChange}
            label="GroupId"
            tooltip="rocket mq consumer group id"
          />
        </div>
      </>
    );
  }
}
