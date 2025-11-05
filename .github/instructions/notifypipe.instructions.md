We are building a self-hosted app that listens to Docker events and sends push notifications (via Slack, Telegram, Discord, etc.) whenever containers are successfully deployed or fail to deploy.

Think of it as a “notifications manager for Docker” — visually inspired by Beszel, technically powered by:
• Golang backend
• PocketBase (for lightweight local data & API)
• Shoutrrr for sending notifications to external channels

The goal: Let users self-host the app easily (via Docker Compose) and manage which containers trigger notifications.

⸻

🧩 Core Components

1. Backend (Golang)
   • Responsibilities:
   • Listen to Docker events using the Docker API socket (/var/run/docker.sock).
   • Detect:
   • container start → mark as deployed successfully
   • container die → detect failed deployment
   • Send notifications via Shoutrrr when relevant events occur.
   • Provide REST API endpoints for:
   • Listing containers and their notification settings
   • Managing connected notification channels
   • Managing global app configuration (webhook URLs, tokens, etc.)
   • Health/status endpoints for UI
   • Dependencies:
   • Docker SDK for Go
   • PocketBase SDK
   • Shoutrrr (https://containrrr.dev/shoutrrr/)
   • Gorilla Mux / Fiber (for API routing)
2. Database Layer (PocketBase)
   • Purpose: lightweight persistence for configuration
   • Collections:
   • notifications → { type, url, enabled_events }
   • containers → { name, container_id, notify_on_success, notify_on_failure }
   • events_log → { container_id, status, timestamp }
   • PocketBase runs embedded or as a separate lightweight container.

⸻

3. Notifications (Shoutrrr)
   • Use Shoutrrr to send messages to multiple destinations:
   • Slack
   • Telegram
   • Discord
   • Email (optional)
   • Each notification target has its URL saved in PocketBase.
   • Allow users to test notification via API endpoint /api/notifications/test.

Example Shoutrrr usage:

```gorouter, _ := shoutrrr.CreateSender("slack://TOKEN@channel")
router.Send("✅ Container 'web' deployed successfully!")
```

4. UI / UX (Inspired by Beszel)
   • Clean, dark-neutral dashboard.
   • TailwindCSS for styling (reuse Beszel-like layout).
   • Built-in onboarding setup flow on first run:
   1. Configure Docker socket access.
   2. Add at least one notification channel (Slack, Telegram, etc.)
   3. Choose which containers to monitor.
   4. Save and start listening for events.
      • Dashboard sections:
      • Overview → recent events
      • Containers → toggle success/failure notifications
      • Notifications → list + edit channels
      • Settings → environment setup (Docker, app info)
      • Logs → last 50 events

⸻

5. First-Run Setup (Initialization Flow)

When the app is started for the first time: 1. Detect if configuration exists (pocketbase.db missing). 2. Serve /setup wizard (in UI):
• Step 1: Connect Docker socket
• Step 2: Add notification channel(s)
• Step 3: Select which containers to monitor
• Step 4: Test and Finish 3. Once setup completes → redirect to dashboard.

The setup can also be re-run anytime from the settings panel.

6. Deployment (Docker Compose)

A minimal docker-compose.yml example:

````version: '3.9'
services:
  notify-manager:
    image: ghcr.io/yourorg/notify-manager:latest
    container_name: notify-manager
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/app/data
    ports:
      - "8080:8080"
    environment:
      - POCKETBASE_URL=http://localhost:8090
      - BASE_URL=http://localhost:8080
      ```
User simply runs:
```bash
docker-compose up -d
````

Then visits http://localhost:8080 to complete setup.

7. Notifications Logic
   Event Type
   Trigger Condition
   Example Message
   Deployment Success
   Container starts without error (status=running)
   ✅ Container 'web' deployed successfully.
   Deployment Failure
   Container dies immediately or exits with code ≠ 0
   ❌ Container 'web' failed to deploy. Exit code: 1

Each container has user-configurable toggles:
{
"container": "web",
"notify_on_success": true,
"notify_on_failure": false
}

🚀 Future Enhancements
• Support for resource monitoring alerts (CPU, memory, disk)
• Integration with webhooks (custom HTTP POST on event)
• Support for multi-server setups
• Role-based access or API keys
• Option to export logs to Prometheus/Grafana

⸻

🧭 Summary for the Agent

You are to: 1. Build the Golang backend that listens to Docker events and sends notifications via Shoutrrr. 2. Integrate PocketBase as embedded DB and API layer. 3. Develop a Beszel-inspired UI (dark dashboard with smooth UX). 4. Implement first-run setup wizard. 5. Expose API endpoints for managing containers, channels, and events. 6. Provide a Docker Compose setup for easy self-hosted installation.

⸻
