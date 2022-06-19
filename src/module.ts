import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { RocketMqQuery, SlsDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, RocketMqQuery, SlsDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
