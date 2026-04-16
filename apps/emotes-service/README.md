# Emotes Service

Service responsible for ingesting Twitch emotes. Infrastructure managed with Pulumi (Go) on AWS.

## Prerequisites

- Go
- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- AWS credentials configured
- Bun for running scripts

For exact versions needed see the [mise config](../../mise.toml).

## Local Development

Ephemeral stacks use [Pulumi ESC](https://www.pulumi.com/docs/esc/).

```bash
# Create your ephemeral stack (one-time)
bun run local:create

# Preview infrastructure changes
bun run local:preview

# Deploy
bun run local:up

# Tear down and remove stack
bun run local:destroy
```

## Configuration

Configurations can be found from the following locations:

| Source                                    | Layer                                            |
|-------------------------------------------|--------------------------------------------------|
| Pulumi.yaml                               | Project-level                                    |
| emotes-service/local-dev                  | ESC environment Pulumi Cloud (local development) |
| Dynamically via scripts in `package.json` | Stack-level                                      |

Stack-level config overrides ESC values. ESC values override project-level.
