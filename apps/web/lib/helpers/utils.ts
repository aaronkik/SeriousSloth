export const capitaliseFirstLetter = (string: string) =>
  string.replace(/^\w/, (str) => str.toUpperCase());

export const getEnvOrThrow = (key: string): string => {
  const env = process.env[key];

  if (!env) {
    throw new Error(`Env variable undefined: ${key}`);
  }

  return env;
};
