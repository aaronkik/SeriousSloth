import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

/** Public read API for Twitch emotes. */
@service(#{ title: "emotes-service" })
@info(#{ version: "1.0.0" })
@useAuth(ApiKeyAuth<ApiKeyLocation.header, "x-api-key">)
@extension("x-amazon-apigateway-api-key-source", "HEADER")
namespace EmotesService;

model Emote {
  id: string;
  name: string;
  format: ("static" | "animated")[];
  scale: ("1.0" | "2.0" | "3.0")[];
  theme_mode: ("light" | "dark")[];
  images: {
    url_1x: url;
    url_2x: url;
    url_4x: url;
  };
}

model ActiveEmote {
  emote: Emote,
  addedAt: utcDateTime
}

@error
model Error {
  message: string;
}

@route("/emotes")
@tag("Emotes")
interface Emotes {
  /** List active emotes for an channelId */
  @get
  @route("{channelId}")
  @extension(
    "x-amazon-apigateway-integration",
    #{
      type: "aws_proxy",
      httpMethod: "POST",
      payloadFormatVersion: "1.0",
      uri: "arn:\${AWS::Partition}:apigateway:\${AWS::Region}:lambda:path/2015-03-31/functions/arn:\${AWS::Partition}:lambda:\${AWS::Region}:\${AWS::AccountId}:function:__GET_EMOTES_FUNCTION_NAME__/invocations",
      passthroughBehavior: "when_no_match",
    }
  )
  list(@path channelId: string): ActiveEmote[] | Error;
}
