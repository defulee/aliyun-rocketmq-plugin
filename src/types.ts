import { DataQuery, DataSourceJsonData } from '@grafana/data';

export enum Action {
  ConsumerAccumulate = 'ConsumerAccumulate',
  TrendTopicInputTps = 'TrendTopicInputTps',
  TrendGroupOutputTps = 'TrendGroupOutputTps',
}

export interface RocketMqQuery extends DataQuery {
  action?: Action;
  groupId?: string;
  topic?: string;
}

/**
 * These are options configured for each DataSource instance.
 */
export interface SlsDataSourceOptions extends DataSourceJsonData {
  accessKeyId?: string;
  region?: string;
  instanceId?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface SlsSecureJsonData {
  accessKeySecret?: string;
}
