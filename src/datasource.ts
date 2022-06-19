import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { SlsDataSourceOptions, RocketMqQuery } from './types';

export class DataSource extends DataSourceWithBackend<RocketMqQuery, SlsDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SlsDataSourceOptions>) {
    super(instanceSettings);
  }
}
