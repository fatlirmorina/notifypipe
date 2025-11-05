# NotifyPipe Documentation

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [API Reference](#api-reference)
4. [Notification Channels](#notification-channels)
5. [Docker Events](#docker-events)
6. [Troubleshooting](#troubleshooting)

## Installation

### Using Docker Compose (Recommended)

```bash
git clone https://github.com/fatlirmorina/notifypipe.git
cd notifypipe
docker-compose up -d
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/fatlirmorina/notifypipe.git
cd notifypipe

# Install dependencies
go mod download

# Build
make build

# Run
./bin/notifypipe
```

## Configuration

### Environment Variables

| Variable        | Description                    | Default                 |
| --------------- | ------------------------------ | ----------------------- |
| `PORT`          | HTTP server port               | `8080`                  |
| `DOCKER_SOCKET` | Docker socket path             | `/var/run/docker.sock`  |
| `DATA_DIR`      | Data directory                 | `./data`                |
| `BASE_URL`      | Base URL                       | `http://localhost:8080` |
| `LOG_LEVEL`     | Log level (info, debug, error) | `info`                  |

### Configuration File

You can also use a `.env` file:

```bash
cp .env.example .env
# Edit .env with your settings
```

## API Reference

### Health Check

```http
GET /api/health
```

Response:

```json
{
  "status": "ok",
  "service": "NotifyPipe",
  "version": "1.0.0"
}
```

### Containers

#### List Containers

```http
GET /api/containers
```

#### Get Container

```http
GET /api/containers/:id
```

#### Update Container Settings

```http
PUT /api/containers/:id
Content-Type: application/json

{
  "notify_on_success": true,
  "notify_on_failure": true
}
```

### Notifications

#### List Notification Channels

```http
GET /api/notifications
```

#### Create Notification Channel

```http
POST /api/notifications
Content-Type: application/json

{
  "name": "My Slack Channel",
  "type": "slack",
  "url": "slack://token@channel"
}
```

#### Update Notification Channel

```http
PUT /api/notifications/:id
Content-Type: application/json

{
  "name": "Updated Name",
  "enabled": true
}
```

#### Delete Notification Channel

```http
DELETE /api/notifications/:id
```

#### Test Notification

```http
POST /api/notifications/test
Content-Type: application/json

{
  "url": "slack://token@channel"
}
```

### Events

#### List Events

```http
GET /api/events
```

#### Get Container Events

```http
GET /api/events/:containerId
```

### Statistics

```http
GET /api/stats
```

## Notification Channels

NotifyPipe uses [Shoutrrr](https://containrrr.dev/shoutrrr/) for sending notifications. Here are examples for different services:

### Slack

```
slack://token@channel
```

Example:

```
slack://xoxb-your-token@general
```

### Telegram

```
telegram://token@telegram?chats=@channel
```

Example:

```
telegram://123456789:ABCdefGHIjklMNOpqrsTUVwxyz@telegram?chats=@mychannel
```

### Discord

```
discord://token@id
```

Example:

```
discord://webhook-token@webhook-id
```

### Email (SMTP)

```
smtp://username:password@host:port/?from=sender@example.com&to=recipient@example.com
```

Example:

```
smtp://user:pass@smtp.gmail.com:587/?from=notify@example.com&to=admin@example.com
```

## Docker Events

NotifyPipe monitors the following Docker events:

### Container Start (Success)

Triggered when a container starts successfully and is running.

**Notification**: "✅ Container 'name' deployed successfully"

### Container Die (Failure)

Triggered when a container exits with a non-zero exit code.

**Notification**: "❌ Container 'name' failed to deploy. Exit code: X"

### Container Create

Automatically tracks new containers in the system.

## Troubleshooting

### NotifyPipe can't connect to Docker

**Issue**: `Error connecting to Docker daemon`

**Solution**:

1. Ensure Docker is running
2. Check socket permissions: `ls -l /var/run/docker.sock`
3. Add user to docker group: `sudo usermod -aG docker $USER`

### Notifications not sending

**Issue**: Notifications are not being received

**Solution**:

1. Test the notification URL using the "Test" button in the UI
2. Check the notification URL format (see Shoutrrr docs)
3. Verify the notification channel is enabled
4. Check container notification settings

### Database errors

**Issue**: PocketBase errors or database locked

**Solution**:

1. Stop NotifyPipe
2. Remove `data/pb_data` directory
3. Restart NotifyPipe (database will be recreated)

### Port already in use

**Issue**: `bind: address already in use`

**Solution**:

1. Change the PORT environment variable
2. Or stop the service using port 8080: `lsof -ti:8080 | xargs kill`

## Getting Help

- Open an issue on [GitHub](https://github.com/fatlirmorina/notifypipe/issues)
- Check existing issues for solutions
- Read the [Contributing Guide](CONTRIBUTING.md)

---

**Need more help?** Feel free to ask in the issues section!
