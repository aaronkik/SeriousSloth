'use strict';

const { loadEnvConfig } = require('@next/env');
loadEnvConfig(process.cwd());

const getDeploymentEnv = () => {
  const env = process.env.VERCEL_ENV;
  switch (env) {
    case 'production':
      return 'prod';
    case 'preview':
      return env;
    case 'development':
      return 'dev';
    default:
      return 'unknown';
  }
};

/**
 * New Relic agent configuration.
 */
exports.config = {
  app_name: [`serious-sloth-web-${getDeploymentEnv()}`],
  application_logging: {
    forwarding: {
      enabled: true,
    },
  },
  logging: {
    enabled: true,
    level: 'info',
    filepath: 'stdout',
  },
  license_key: process.env.NEW_RELIC_LICENSE_KEY,

  /**
   * When true, all request headers except for those listed in attributes.exclude
   * will be captured for all traces, unless otherwise specified in a destination's
   * attributes include/exclude lists.
   */
  allow_all_headers: true,
  attributes: {
    /**
     * Prefix of attributes to exclude from all destinations. Allows * as wildcard
     * at end.
     *
     * NOTE: If excluding headers, they must be in camelCase form to be filtered.
     *
     * @name NEW_RELIC_ATTRIBUTES_EXCLUDE
     */
    exclude: [
      'request.headers.cookie',
      'request.headers.authorization',
      'request.headers.proxyAuthorization',
      'request.headers.setCookie*',
      'request.headers.x*',
      'response.headers.cookie',
      'response.headers.authorization',
      'response.headers.proxyAuthorization',
      'response.headers.setCookie*',
      'response.headers.x*',
    ],
  },
  labels: `env:${process.env.VERCEL_ENV ?? 'unknown'};commit:${process.env.VERCEL_GIT_COMMIT_SHA ?? 'unknown'};region:${process.env.VERCEL_REGION ?? 'unknown'};`,
};
