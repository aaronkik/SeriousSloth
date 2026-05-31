import nr from 'newrelic';

interface TraceOptions<T> {
  /** Segment name shown in the New Relic trace waterfall. */
  name: string;
  /** Async work to time. */
  handler: () => Promise<T>;
  /** Record as its own metric. Defaults to true. */
  record?: boolean;
  /** Custom attributes attached to the enclosing transaction. */
  attributes?: Record<string, string | number | boolean>;
}

/**
 * Wraps async work in a New Relic segment, capturing its duration.
 * Defaults to recording a metric; pass `record: false` to opt out.
 */
export const trace = <T>({
  name,
  handler,
  record = true,
  attributes,
}: TraceOptions<T>): Promise<T> => {
  if (attributes) {
    nr.addCustomAttributes(attributes);
  }

  return nr.startSegment(name, record, handler);
};

/** Attach custom attributes to the current transaction. */
export const setAttributes = (
  attributes: Record<string, string | number | boolean>,
): void => {
  nr.addCustomAttributes(attributes);
};
