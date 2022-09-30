/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
  "/": {
    get: {
      responses: {
        /** OK */
        200: unknown;
      };
    };
  };
  "/stream_history_timestamp": {
    get: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history_timestamp.id"];
          time?: parameters["rowFilter.stream_history_timestamp.time"];
          /** Filtering Columns */
          select?: parameters["select"];
          /** Ordering */
          order?: parameters["order"];
          /** Limiting and Pagination */
          offset?: parameters["offset"];
          /** Limiting and Pagination */
          limit?: parameters["limit"];
        };
        header: {
          /** Limiting and Pagination */
          Range?: parameters["range"];
          /** Limiting and Pagination */
          "Range-Unit"?: parameters["rangeUnit"];
          /** Preference */
          Prefer?: parameters["preferCount"];
        };
      };
      responses: {
        /** OK */
        200: {
          schema: definitions["stream_history_timestamp"][];
        };
        /** Partial Content */
        206: unknown;
      };
    };
    post: {
      parameters: {
        body: {
          /** stream_history_timestamp */
          stream_history_timestamp?: definitions["stream_history_timestamp"];
        };
        query: {
          /** Filtering Columns */
          select?: parameters["select"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** Created */
        201: unknown;
      };
    };
    delete: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history_timestamp.id"];
          time?: parameters["rowFilter.stream_history_timestamp.time"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** No Content */
        204: never;
      };
    };
    patch: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history_timestamp.id"];
          time?: parameters["rowFilter.stream_history_timestamp.time"];
        };
        body: {
          /** stream_history_timestamp */
          stream_history_timestamp?: definitions["stream_history_timestamp"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** No Content */
        204: never;
      };
    };
  };
  "/stream_history": {
    get: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history.id"];
          stream_id?: parameters["rowFilter.stream_history.stream_id"];
          user_id?: parameters["rowFilter.stream_history.user_id"];
          user_login?: parameters["rowFilter.stream_history.user_login"];
          user_name?: parameters["rowFilter.stream_history.user_name"];
          game_id?: parameters["rowFilter.stream_history.game_id"];
          game_name?: parameters["rowFilter.stream_history.game_name"];
          stream_title?: parameters["rowFilter.stream_history.stream_title"];
          viewer_count?: parameters["rowFilter.stream_history.viewer_count"];
          started_at?: parameters["rowFilter.stream_history.started_at"];
          language?: parameters["rowFilter.stream_history.language"];
          is_mature?: parameters["rowFilter.stream_history.is_mature"];
          stream_history_timestamp_id?: parameters["rowFilter.stream_history.stream_history_timestamp_id"];
          /** Filtering Columns */
          select?: parameters["select"];
          /** Ordering */
          order?: parameters["order"];
          /** Limiting and Pagination */
          offset?: parameters["offset"];
          /** Limiting and Pagination */
          limit?: parameters["limit"];
        };
        header: {
          /** Limiting and Pagination */
          Range?: parameters["range"];
          /** Limiting and Pagination */
          "Range-Unit"?: parameters["rangeUnit"];
          /** Preference */
          Prefer?: parameters["preferCount"];
        };
      };
      responses: {
        /** OK */
        200: {
          schema: definitions["stream_history"][];
        };
        /** Partial Content */
        206: unknown;
      };
    };
    post: {
      parameters: {
        body: {
          /** stream_history */
          stream_history?: definitions["stream_history"];
        };
        query: {
          /** Filtering Columns */
          select?: parameters["select"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** Created */
        201: unknown;
      };
    };
    delete: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history.id"];
          stream_id?: parameters["rowFilter.stream_history.stream_id"];
          user_id?: parameters["rowFilter.stream_history.user_id"];
          user_login?: parameters["rowFilter.stream_history.user_login"];
          user_name?: parameters["rowFilter.stream_history.user_name"];
          game_id?: parameters["rowFilter.stream_history.game_id"];
          game_name?: parameters["rowFilter.stream_history.game_name"];
          stream_title?: parameters["rowFilter.stream_history.stream_title"];
          viewer_count?: parameters["rowFilter.stream_history.viewer_count"];
          started_at?: parameters["rowFilter.stream_history.started_at"];
          language?: parameters["rowFilter.stream_history.language"];
          is_mature?: parameters["rowFilter.stream_history.is_mature"];
          stream_history_timestamp_id?: parameters["rowFilter.stream_history.stream_history_timestamp_id"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** No Content */
        204: never;
      };
    };
    patch: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history.id"];
          stream_id?: parameters["rowFilter.stream_history.stream_id"];
          user_id?: parameters["rowFilter.stream_history.user_id"];
          user_login?: parameters["rowFilter.stream_history.user_login"];
          user_name?: parameters["rowFilter.stream_history.user_name"];
          game_id?: parameters["rowFilter.stream_history.game_id"];
          game_name?: parameters["rowFilter.stream_history.game_name"];
          stream_title?: parameters["rowFilter.stream_history.stream_title"];
          viewer_count?: parameters["rowFilter.stream_history.viewer_count"];
          started_at?: parameters["rowFilter.stream_history.started_at"];
          language?: parameters["rowFilter.stream_history.language"];
          is_mature?: parameters["rowFilter.stream_history.is_mature"];
          stream_history_timestamp_id?: parameters["rowFilter.stream_history.stream_history_timestamp_id"];
        };
        body: {
          /** stream_history */
          stream_history?: definitions["stream_history"];
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferReturn"];
        };
      };
      responses: {
        /** No Content */
        204: never;
      };
    };
  };
  "/stream_history_timestamp_with_total_view_count": {
    get: {
      parameters: {
        query: {
          id?: parameters["rowFilter.stream_history_timestamp_with_total_view_count.id"];
          time?: parameters["rowFilter.stream_history_timestamp_with_total_view_count.time"];
          total_viewer_count?: parameters["rowFilter.stream_history_timestamp_with_total_view_count.total_viewer_count"];
          total_streams?: parameters["rowFilter.stream_history_timestamp_with_total_view_count.total_streams"];
          /** Filtering Columns */
          select?: parameters["select"];
          /** Ordering */
          order?: parameters["order"];
          /** Limiting and Pagination */
          offset?: parameters["offset"];
          /** Limiting and Pagination */
          limit?: parameters["limit"];
        };
        header: {
          /** Limiting and Pagination */
          Range?: parameters["range"];
          /** Limiting and Pagination */
          "Range-Unit"?: parameters["rangeUnit"];
          /** Preference */
          Prefer?: parameters["preferCount"];
        };
      };
      responses: {
        /** OK */
        200: {
          schema: definitions["stream_history_timestamp_with_total_view_count"][];
        };
        /** Partial Content */
        206: unknown;
      };
    };
  };
  "/rpc/insert_streams_with_timestamp": {
    post: {
      parameters: {
        body: {
          args: {
            /** Format: jsonb */
            streams: unknown;
          };
        };
        header: {
          /** Preference */
          Prefer?: parameters["preferParams"];
        };
      };
      responses: {
        /** OK */
        200: unknown;
      };
    };
  };
}

export interface definitions {
  /** @description The time streams are fetched from the Twitch API */
  stream_history_timestamp: {
    /**
     * Format: bigint
     * @description Note:
     * This is a Primary Key.<pk/>
     */
    id: number;
    /**
     * Format: timestamp with time zone
     * @default now()
     */
    time: string;
  };
  /** @description Streams fetched from the Twitch API */
  stream_history: {
    /**
     * Format: bigint
     * @description Note:
     * This is a Primary Key.<pk/>
     */
    id: number;
    /** Format: text */
    stream_id: string;
    /** Format: text */
    user_id: string;
    /** Format: text */
    user_login: string;
    /** Format: text */
    user_name: string;
    /** Format: text */
    game_id: string;
    /** Format: text */
    game_name: string;
    /** Format: text */
    stream_title: string;
    /** Format: integer */
    viewer_count: number;
    /** Format: timestamp with time zone */
    started_at: string;
    /** Format: text */
    language: string;
    /** Format: boolean */
    is_mature: boolean;
    /**
     * Format: integer
     * @description Note:
     * This is a Foreign Key to `stream_history_timestamp.id`.<fk table='stream_history_timestamp' column='id'/>
     */
    stream_history_timestamp_id: number;
  };
  stream_history_timestamp_with_total_view_count: {
    /**
     * Format: bigint
     * @description Note:
     * This is a Primary Key.<pk/>
     */
    id?: number;
    /** Format: timestamp with time zone */
    time?: string;
    /** Format: bigint */
    total_viewer_count?: number;
    /** Format: bigint */
    total_streams?: number;
  };
}

export interface parameters {
  /**
   * @description Preference
   * @enum {string}
   */
  preferParams: "params=single-object";
  /**
   * @description Preference
   * @enum {string}
   */
  preferReturn: "return=representation" | "return=minimal" | "return=none";
  /**
   * @description Preference
   * @enum {string}
   */
  preferCount: "count=none";
  /** @description Filtering Columns */
  select: string;
  /** @description On Conflict */
  on_conflict: string;
  /** @description Ordering */
  order: string;
  /** @description Limiting and Pagination */
  range: string;
  /**
   * @description Limiting and Pagination
   * @default items
   */
  rangeUnit: string;
  /** @description Limiting and Pagination */
  offset: string;
  /** @description Limiting and Pagination */
  limit: string;
  /** @description stream_history_timestamp */
  "body.stream_history_timestamp": definitions["stream_history_timestamp"];
  /** Format: bigint */
  "rowFilter.stream_history_timestamp.id": string;
  /** Format: timestamp with time zone */
  "rowFilter.stream_history_timestamp.time": string;
  /** @description stream_history */
  "body.stream_history": definitions["stream_history"];
  /** Format: bigint */
  "rowFilter.stream_history.id": string;
  /** Format: text */
  "rowFilter.stream_history.stream_id": string;
  /** Format: text */
  "rowFilter.stream_history.user_id": string;
  /** Format: text */
  "rowFilter.stream_history.user_login": string;
  /** Format: text */
  "rowFilter.stream_history.user_name": string;
  /** Format: text */
  "rowFilter.stream_history.game_id": string;
  /** Format: text */
  "rowFilter.stream_history.game_name": string;
  /** Format: text */
  "rowFilter.stream_history.stream_title": string;
  /** Format: integer */
  "rowFilter.stream_history.viewer_count": string;
  /** Format: timestamp with time zone */
  "rowFilter.stream_history.started_at": string;
  /** Format: text */
  "rowFilter.stream_history.language": string;
  /** Format: boolean */
  "rowFilter.stream_history.is_mature": string;
  /** Format: integer */
  "rowFilter.stream_history.stream_history_timestamp_id": string;
  /** @description stream_history_timestamp_with_total_view_count */
  "body.stream_history_timestamp_with_total_view_count": definitions["stream_history_timestamp_with_total_view_count"];
  /** Format: bigint */
  "rowFilter.stream_history_timestamp_with_total_view_count.id": string;
  /** Format: timestamp with time zone */
  "rowFilter.stream_history_timestamp_with_total_view_count.time": string;
  /** Format: bigint */
  "rowFilter.stream_history_timestamp_with_total_view_count.total_viewer_count": string;
  /** Format: bigint */
  "rowFilter.stream_history_timestamp_with_total_view_count.total_streams": string;
}

export interface operations {}

export interface external {}