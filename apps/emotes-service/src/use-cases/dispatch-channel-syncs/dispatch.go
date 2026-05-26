package dispatchchannelsyncs

import (
	"context"
	"emotes-service/src/adapters/secondary/channels_store"
	"emotes-service/src/environment"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const sqsBatchSize = 10

type SyncMessage struct {
	TwitchId string `json:"twitchId"`
}

var sqsClient *sqs.Client

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	sqsClient = sqs.NewFromConfig(cfg)
}

// Execute scans the channels registry and enqueues one SQS message per
// channel. The downstream channel-sync worker consumes the queue with
// reserved concurrency, capping in-flight Twitch calls.
func Execute(ctx context.Context) error {
	txn := newrelic.FromContext(ctx)
	queueUrl := environment.GetOrFatal("CHANNEL_SYNC_QUEUE_URL")

	seg := txn.StartSegment("channels_store.QueryAll")
	channels, err := channels_store.QueryAll(ctx)
	seg.End()
	if err != nil {
		return err
	}

	if len(channels) == 0 {
		slog.InfoContext(ctx, "no channels registered, nothing to dispatch")
		return nil
	}

	slog.InfoContext(ctx, "dispatching channel syncs", "count", len(channels), "channels", channels)

	for start := 0; start < len(channels); start += sqsBatchSize {
		end := min(start+sqsBatchSize, len(channels))
		chunk := channels[start:end]

		entries := make([]types.SendMessageBatchRequestEntry, 0, len(chunk))
		for i, c := range chunk {
			body, err := json.Marshal(SyncMessage{TwitchId: c.TwitchId})
			if err != nil {
				slog.ErrorContext(ctx, "marshal sync message failed", "twitchId", c.TwitchId, "error", err)
				return err
			}
			entries = append(entries, types.SendMessageBatchRequestEntry{
				Id:          aws.String("ch-" + strconv.Itoa(start+i)),
				MessageBody: aws.String(string(body)),
			})
		}

		sqsSeg := txn.StartSegment("sqs.SendMessageBatch")
		out, err := sqsClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: aws.String(queueUrl),
			Entries:  entries,
		})
		sqsSeg.End()
		if err != nil {
			slog.ErrorContext(ctx, "SendMessageBatch failed", "chunkStart", start, "error", err)
			return err
		}
		if len(out.Failed) > 0 {
			slog.ErrorContext(ctx, "partial batch failures", "failed", out.Failed)
			return fmt.Errorf("sqs SendMessageBatch had %d failed entries", len(out.Failed))
		}
	}

	return nil
}
