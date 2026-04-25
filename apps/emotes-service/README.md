# Emotes Service

Service responsible for ingesting Twitch emotes. Infrastructure managed with Pulumi (Go) on AWS.

## Prerequisites

- Go
- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- AWS credentials configured
- [Task](https://taskfile.dev/) for running tasks

For exact versions needed see the [mise config](../../mise.toml).

## Local Development

Ephemeral stacks use [Pulumi ESC](https://www.pulumi.com/docs/esc/).

```bash
# Create your ephemeral stack (one-time)
task create
```

```sh
# Preview infrastructure changes
task preview
```

```sh
# Deploy
task deploy
```

```sh
# Tear down and remove stack
task destroy:local
```

## Configuration

Configurations can be found from the following locations:

| Source                                    | Layer                                            |
|-------------------------------------------|--------------------------------------------------|
| Pulumi.yaml                               | Project-level                                    |
| emotes-service/local-dev                  | ESC environment Pulumi Cloud (local development) |
| Dynamically via tasks in `Taskfile.yml`   | Stack-level                                      |

Stack-level config overrides ESC values. ESC values override project-level.
