// API Base URL
const API_BASE = "/api";

// State
let currentTab = "overview";

// Initialize app
document.addEventListener("DOMContentLoaded", () => {
  loadStats();
  loadContainers();
  loadNotifications();
  loadEvents();

  // Refresh data every 30 seconds
  setInterval(() => {
    loadStats();
    loadContainers();
    loadEvents();
  }, 30000);
});

// Tab Management
function showTab(tabName) {
  // Hide all tabs
  document.querySelectorAll(".tab-content").forEach((tab) => {
    tab.classList.add("hidden");
  });

  // Remove active state from all buttons
  document.querySelectorAll(".tab-button").forEach((button) => {
    button.classList.remove("border-blue-500", "text-blue-500");
    button.classList.add("border-transparent", "text-gray-400");
  });

  // Show selected tab
  document.getElementById(`${tabName}-tab`).classList.remove("hidden");

  // Set active button
  event.target.classList.add("border-blue-500", "text-blue-500");
  event.target.classList.remove("border-transparent", "text-gray-400");

  currentTab = tabName;

  // Refresh data for the tab
  if (tabName === "containers") loadContainers();
  if (tabName === "notifications") loadNotifications();
  if (tabName === "events") loadEvents();
}

// Load Statistics
async function loadStats() {
  try {
    const response = await fetch(`${API_BASE}/stats`);
    const data = await response.json();

    document.getElementById("containers-count").textContent = data.containers_count || 0;
    document.getElementById("notifications-count").textContent = data.notifications_count || 0;
    document.getElementById("events-count").textContent = data.events_count || 0;
  } catch (error) {
    console.error("Error loading stats:", error);
  }
}

// Load Containers
async function loadContainers() {
  try {
    const response = await fetch(`${API_BASE}/containers`);
    const containers = await response.json();

    const containersList = document.getElementById("containers-list");

    if (containers.length === 0) {
      containersList.innerHTML = '<p class="text-gray-400 text-center py-8">No containers found</p>';
      return;
    }

    containersList.innerHTML = containers
      .map(
        (container) => `
            <div class="border border-dark-border rounded-lg p-4 hover:bg-dark-hover transition">
                <div class="flex items-center justify-between">
                    <div class="flex-1">
                        <h3 class="font-semibold text-white">${container.name}</h3>
                        <p class="text-sm text-gray-400 mt-1">${container.image}</p>
                        <div class="flex items-center mt-2 space-x-2">
                            <span class="px-2 py-1 rounded text-xs ${getStatusClass(container.state)}">${
          container.state
        }</span>
                            <span class="text-xs text-gray-500">${container.id.substring(0, 12)}</span>
                        </div>
                    </div>
                    <div class="flex flex-col space-y-2">
                        <label class="flex items-center space-x-2 cursor-pointer">
                            <input type="checkbox" ${container.notify_on_success ? "checked" : ""} 
                                   onchange="updateContainerSettings('${
                                     container.id
                                   }', 'notify_on_success', this.checked)"
                                   class="w-4 h-4 text-blue-600 bg-dark-bg border-dark-border rounded focus:ring-blue-500">
                            <span class="text-sm text-gray-300">✅ Success</span>
                        </label>
                        <label class="flex items-center space-x-2 cursor-pointer">
                            <input type="checkbox" ${container.notify_on_failure ? "checked" : ""} 
                                   onchange="updateContainerSettings('${
                                     container.id
                                   }', 'notify_on_failure', this.checked)"
                                   class="w-4 h-4 text-blue-600 bg-dark-bg border-dark-border rounded focus:ring-blue-500">
                            <span class="text-sm text-gray-300">❌ Failure</span>
                        </label>
                    </div>
                </div>
            </div>
        `
      )
      .join("");
  } catch (error) {
    console.error("Error loading containers:", error);
    document.getElementById("containers-list").innerHTML =
      '<p class="text-red-400 text-center py-8">Error loading containers</p>';
  }
}

// Update Container Settings
async function updateContainerSettings(containerId, field, value) {
  try {
    // Get current settings first
    const getResponse = await fetch(`${API_BASE}/containers/${containerId}`);
    const container = await getResponse.json();

    const settings = {
      notify_on_success: container.notify_on_success,
      notify_on_failure: container.notify_on_failure,
    };

    settings[field] = value;

    const response = await fetch(`${API_BASE}/containers/${containerId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(settings),
    });

    if (response.ok) {
      showToast("Container settings updated", "success");
    }
  } catch (error) {
    console.error("Error updating container:", error);
    showToast("Failed to update container settings", "error");
  }
}

// Load Notifications
async function loadNotifications() {
  try {
    const response = await fetch(`${API_BASE}/notifications`);
    const notifications = await response.json();

    const notificationsList = document.getElementById("notifications-list");

    if (notifications.length === 0) {
      notificationsList.innerHTML = '<p class="text-gray-400 text-center py-8">No notification channels configured</p>';
      return;
    }

    notificationsList.innerHTML = notifications
      .map(
        (notif) => `
            <div class="border border-dark-border rounded-lg p-4 hover:bg-dark-hover transition">
                <div class="flex items-center justify-between">
                    <div class="flex-1">
                        <h3 class="font-semibold text-white">${notif.name}</h3>
                        <p class="text-sm text-gray-400 mt-1">${notif.type}</p>
                        <p class="text-xs text-gray-500 mt-1 font-mono">${maskUrl(notif.url)}</p>
                    </div>
                    <div class="flex items-center space-x-3">
                        <label class="flex items-center space-x-2 cursor-pointer">
                            <input type="checkbox" ${notif.enabled ? "checked" : ""} 
                                   onchange="toggleNotification('${notif.id}', this.checked)"
                                   class="w-4 h-4 text-blue-600 bg-dark-bg border-dark-border rounded focus:ring-blue-500">
                            <span class="text-sm text-gray-300">Enabled</span>
                        </label>
                        <button onclick="deleteNotification('${notif.id}')" class="text-red-400 hover:text-red-300">
                            🗑️
                        </button>
                    </div>
                </div>
            </div>
        `
      )
      .join("");
  } catch (error) {
    console.error("Error loading notifications:", error);
  }
}

// Load Events
async function loadEvents() {
  try {
    const response = await fetch(`${API_BASE}/events`);
    const events = await response.json();

    const eventsList = document.getElementById("events-list");
    const recentEvents = document.getElementById("recent-events");

    if (events.length === 0) {
      const emptyMessage = '<p class="text-gray-400 text-center py-8">No events yet</p>';
      eventsList.innerHTML = emptyMessage;
      recentEvents.innerHTML = emptyMessage;
      return;
    }

    const eventHTML = events
      .map(
        (event) => `
            <div class="border border-dark-border rounded p-3 text-sm">
                <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                        <span class="text-lg">${getEventIcon(event.status)}</span>
                        <div>
                            <span class="font-medium text-white">${event.container_name}</span>
                            <span class="text-gray-400 mx-2">•</span>
                            <span class="text-gray-400">${event.message}</span>
                        </div>
                    </div>
                    <span class="text-gray-500 text-xs">${formatDate(event.timestamp)}</span>
                </div>
            </div>
        `
      )
      .join("");

    eventsList.innerHTML = eventHTML;
    recentEvents.innerHTML = events
      .slice(0, 5)
      .map(
        (event) => `
            <div class="flex items-center space-x-3 text-sm">
                <span class="text-lg">${getEventIcon(event.status)}</span>
                <span class="text-white font-medium">${event.container_name}</span>
                <span class="text-gray-400">${event.message}</span>
                <span class="text-gray-500 text-xs ml-auto">${formatDate(event.timestamp)}</span>
            </div>
        `
      )
      .join("");
  } catch (error) {
    console.error("Error loading events:", error);
  }
}

// Modal Management
function showAddNotificationModal() {
  document.getElementById("notification-modal").classList.remove("hidden");
}

function hideAddNotificationModal() {
  document.getElementById("notification-modal").classList.add("hidden");
  document.getElementById("notification-form").reset();
}

// Add Notification
document.getElementById("notification-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const name = document.getElementById("notification-name").value;
  const type = document.getElementById("notification-type").value;
  const url = document.getElementById("notification-url").value;

  try {
    const response = await fetch(`${API_BASE}/notifications`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, type, url }),
    });

    if (response.ok) {
      showToast("Notification channel added", "success");
      hideAddNotificationModal();
      loadNotifications();
      loadStats();
    }
  } catch (error) {
    console.error("Error adding notification:", error);
    showToast("Failed to add notification channel", "error");
  }
});

// Test Notification
async function testNotification() {
  const url = document.getElementById("notification-url").value;

  if (!url) {
    showToast("Please enter a URL", "error");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/notifications/test`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });

    const data = await response.json();

    if (data.success) {
      showToast("Test notification sent!", "success");
    } else {
      showToast(`Test failed: ${data.error}`, "error");
    }
  } catch (error) {
    console.error("Error testing notification:", error);
    showToast("Failed to send test notification", "error");
  }
}

// Toggle Notification
async function toggleNotification(id, enabled) {
  try {
    const response = await fetch(`${API_BASE}/notifications/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ enabled }),
    });

    if (response.ok) {
      showToast(`Notification ${enabled ? "enabled" : "disabled"}`, "success");
    }
  } catch (error) {
    console.error("Error toggling notification:", error);
    showToast("Failed to update notification", "error");
  }
}

// Delete Notification
async function deleteNotification(id) {
  if (!confirm("Are you sure you want to delete this notification channel?")) {
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/notifications/${id}`, {
      method: "DELETE",
    });

    if (response.ok) {
      showToast("Notification channel deleted", "success");
      loadNotifications();
      loadStats();
    }
  } catch (error) {
    console.error("Error deleting notification:", error);
    showToast("Failed to delete notification channel", "error");
  }
}

// Utility Functions
function getStatusClass(status) {
  const classes = {
    running: "bg-green-500/20 text-green-400",
    exited: "bg-red-500/20 text-red-400",
    created: "bg-blue-500/20 text-blue-400",
    paused: "bg-yellow-500/20 text-yellow-400",
  };
  return classes[status] || "bg-gray-500/20 text-gray-400";
}

function getEventIcon(status) {
  const icons = {
    success: "✅",
    failure: "❌",
    stopped: "⏸️",
    created: "🆕",
  };
  return icons[status] || "📋";
}

function maskUrl(url) {
  return url.replace(/:[^:@]*@/, ":***@").replace(/\/\/[^@]*@/, "//***@");
}

function formatDate(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now - date;

  if (diff < 60000) return "just now";
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
  return date.toLocaleDateString();
}

function showToast(message, type = "info") {
  const toast = document.createElement("div");
  toast.className = `fixed bottom-4 right-4 px-6 py-3 rounded-lg text-white ${
    type === "success" ? "bg-green-600" : type === "error" ? "bg-red-600" : "bg-blue-600"
  } shadow-lg z-50`;
  toast.textContent = message;

  document.body.appendChild(toast);

  setTimeout(() => {
    toast.remove();
  }, 3000);
}
