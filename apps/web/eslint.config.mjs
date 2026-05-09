import nextCoreWebVitals from 'eslint-config-next/core-web-vitals';
import prettier from 'eslint-config-prettier';

const config = [
  ...nextCoreWebVitals,
  prettier,
  {
    rules: {
      'no-console': ['error', { allow: ['error'] }],
      'no-unused-vars': ['error'],
    },
  },
  {
    ignores: ['.next/**', 'node_modules/**', 'cypress/**'],
  },
];

export default config;
