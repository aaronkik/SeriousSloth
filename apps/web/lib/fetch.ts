class HttpError<TCode extends number = number> extends Error {
  public readonly statusCode: TCode;
  public readonly url: string | undefined;
  public readonly method: string | undefined;
  public readonly data: Record<string, unknown> | undefined;

  constructor(opts: {
    url?: string;
    method?: string;
    message?: string;
    statusCode: TCode;
    cause?: unknown;
    data?: Record<string, unknown>;
  }) {
    super(opts.message ?? `HTTP Error ${opts.statusCode}`, {
      cause: opts.cause,
    });

    this.name = this.constructor.name;
    this.statusCode = opts.statusCode;
    this.url = opts.url;
    this.method = opts.method;
    this.data = opts.data;

    Error.captureStackTrace?.(this, HttpError);
  }

  public static fromRequest(
    request: Request,
    response: Response,
    parsedError: { data?: Record<string, unknown> },
  ): HttpError {
    return new HttpError({
      message: response.statusText,
      url: response.url,
      method: request.method,
      statusCode: response.status,
      data: parsedError.data,
    });
  }
}

const http = async <T>(path: string, config: RequestInit): Promise<T> => {
  const request = new Request(path, config);
  const response = await fetch(request);

  if (!response.ok) {
    const errJson = await response.json();
    throw HttpError.fromRequest(
      request,
      { ...response, statusText: errJson.message ?? response.statusText },
      errJson,
    );
  }

  return await response.json();
};

export const get = async <T>(
  path: string,
  config?: Omit<RequestInit, 'method'>,
): Promise<T> => {
  const init = { method: 'GET', ...config };
  return await http<T>(path, init);
};
