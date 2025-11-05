# NotifyPipe 🔔

**Free & Open Source Docker Notification Manager**

NotifyPipe is a self-hosted application that listens to Docker events and sends push notifications (via Slack, Telegram, Discord, etc.) whenever containers are deployed or fail. Stop paying for what you can build yourself!

## ✨ Features

- 🐳 **Docker Event Monitoring** - Real-time container deployment tracking
- 📬 **Multi-Channel Notifications** - Slack, Telegram, Discord, Email support
- 🎨 **Clean Dashboard** - Beautiful, dark-themed UI inspired by Beszel
- 🔧 **Easy Setup** - First-run wizard to get started in minutes
- 💾 **Lightweight** - PocketBase-powered local storage
- 🚀 **Self-Hosted** - Your data, your infrastructure
- 🆓 **100% Free** - No subscriptions, no vendor lock-in

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose installed
- Access to `/var/run/docker.sock`

### Installation

1. Clone the repository:

```bash
git clone https://github.com/fatlirmorina/notifypipe.git
cd notifypipe
```

2. Start the application:

```bash
docker-compose up -d
```

3. Open your browser and navigate to:

```
http://localhost:8080
```

4. Follow the setup wizard to:
   - Configure Docker connection
   - Add notification channels
   - Select containers to monitor

## 📋 Notification Channels

NotifyPipe supports multiple notification providers via [Shoutrrr](https://containrrr.dev/shoutrrr/):

- **Slack** - `slack://token@channel`
- **Telegram** - `telegram://token@telegram?chats=@channel-1`
- **Discord** - `discord://token@id`
- **Email** - `smtp://username:password@host:port/?from=sender@example.com&to=recipient@example.com`

## 🏗️ Architecture

- **Backend**: Golang with Fiber framework
- **Database**: PocketBase (embedded SQLite)
- **Notifications**: Shoutrrr
- **Docker API**: Official Docker SDK for Go
- **Frontend**: HTML/CSS/JavaScript with TailwindCSS

## 🛠️ Development

### Build from source

```bash
# Install dependencies
go mod download

# Build the binary
go build -o notifypipe ./cmd/notifypipe

# Run locally
./notifypipe
```

### Environment Variables

| Variable        | Description                   | Default                 |
| --------------- | ----------------------------- | ----------------------- |
| `PORT`          | HTTP server port              | `8080`                  |
| `DOCKER_SOCKET` | Docker socket path            | `/var/run/docker.sock`  |
| `DATA_DIR`      | Data directory for PocketBase | `./data`                |
| `BASE_URL`      | Base URL for the app          | `http://localhost:8080` |

## 📚 API Endpoints

### Containers

- `GET /api/containers` - List all containers
- `PUT /api/containers/:id` - Update container notification settings

### Notifications

- `GET /api/notifications` - List notification channels
- `POST /api/notifications` - Add notification channel
- `DELETE /api/notifications/:id` - Remove notification channel
- `POST /api/notifications/test` - Test notification

### Events

- `GET /api/events` - Get recent events
- `GET /api/events/:containerId` - Get events for specific container

### System

- `GET /api/health` - Health check
- `GET /api/setup` - Check if setup is complete

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- Inspired by [Beszel](https://github.com/henrygd/beszel) for UI/UX design
- Powered by [Shoutrrr](https://containrrr.dev/shoutrrr/) for notifications
- Built with [PocketBase](https://pocketbase.io/) for data persistence

## 💡 Why NotifyPipe?

Don't pay for simple notification services when you can self-host them for free! NotifyPipe gives you complete control over your Docker monitoring and notifications.

---

**Made with ❤️ by the open-source community**
