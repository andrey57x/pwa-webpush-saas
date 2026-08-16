importScripts("/config.js");

self.addEventListener("install", function (event) {
  self.skipWaiting();
});

self.addEventListener("activate", function (event) {
  event.waitUntil(self.clients.claim());
});

self.addEventListener("push", function (event) {
  let payload = { title: "Уведомление", body: "" };

  if (event.data) {
    try {
      payload = event.data.json();
    } catch (e) {
      payload = { title: "Уведомление", body: event.data.text() };
    }
  }

  const title = payload.title || "Уведомление";
  const options = {
    body: payload.body || "",
    data: {
      url: payload.url || "/",
      campaign_id: payload.campaign_id,
      subscription_id: payload.subscription_id,
    },
  };

  const saasUrl = self.PIZZA_CONFIG
    ? self.PIZZA_CONFIG.saasApiUrl
    : "http://localhost:8080/api/v1";

  event.waitUntil(
    self.registration.showNotification(title, options).then(function () {
      if (payload.campaign_id && payload.subscription_id) {
        fetch(`${saasUrl}/feedback/ping`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            campaign_id: payload.campaign_id,
            subscription_id: payload.subscription_id,
            status: "DELIVERED",
          }),
        }).catch(function (err) {
          console.error("Feedback ping error:", err);
        });
      }
    }),
  );
});

self.addEventListener("notificationclick", function (event) {
  event.notification.close();
  const data = event.notification.data || {};
  const targetUrl = data.url || "/";
  const saasUrl = self.PIZZA_CONFIG
    ? self.PIZZA_CONFIG.saasApiUrl
    : "http://localhost:8080/api/v1";

  event.waitUntil(
    clients.matchAll({ type: "window" }).then(function (clientList) {
      if (data.campaign_id && data.subscription_id) {
        fetch(`${saasUrl}/feedback/ping`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            campaign_id: data.campaign_id,
            subscription_id: data.subscription_id,
            status: "CLICKED",
          }),
        }).catch(function (err) {
          console.error("Feedback click error:", err);
        });
      }

      for (let i = 0; i < clientList.length; i++) {
        let client = clientList[i];
        if (client.url === targetUrl && "focus" in client) {
          return client.focus();
        }
      }
      if (clients.openWindow) {
        return clients.openWindow(targetUrl);
      }
    }),
  );
});
