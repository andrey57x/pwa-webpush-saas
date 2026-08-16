// Native Web Push Service Worker
self.addEventListener("push", function (event) {
  console.log("[sw.js] Входящий Push от сервера Mozilla/Google получен!");

  if (!event.data) {
    console.log("[sw.js] Пуш без payload.");
    return;
  }

  let payload;
  try {
    payload = event.data.json();
    console.log("[sw.js] Расшифрованный Payload:", payload);
  } catch (e) {
    payload = { title: "Notification", body: event.data.text() };
  }

  const title = payload.title || "Web Push Notification";
  const options = {
    body: payload.body || "",
    data: {
      url: payload.url || "/",
      campaign_id: payload.campaign_id,
      subscription_id: payload.subscription_id,
    },
  };

  if (payload.icon && payload.icon.length > 0) {
    options.icon = payload.icon;
  }

  event.waitUntil(
    self.registration.showNotification(title, options).then(() => {
      console.log("[sw.js] Системное уведомление успешно показано!");

      // Feedback Loop: Отправка статуса DELIVERED на бэкенд
      if (payload.campaign_id && payload.subscription_id) {
        fetch("/api/v1/feedback/ping", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            campaign_id: payload.campaign_id,
            subscription_id: payload.subscription_id,
            status: "DELIVERED",
          }),
        })
          .then(() => console.log("[sw.js] Feedback DELIVERED отправлен в Go!"))
          .catch((err) =>
            console.error("[sw.js] Ошибка отправки feedback:", err),
          );
      }
    }),
  );
});

self.addEventListener("notificationclick", function (event) {
  console.log("[sw.js] Пользователь КЛИКНУЛ по уведомлению!");
  event.notification.close();

  const notificationData = event.notification.data || {};
  const targetUrl = notificationData.url || "/";

  event.waitUntil(
    clients.matchAll({ type: "window" }).then(function (clientList) {
      // Feedback Loop: Отправка статуса CLICKED на бэкенд
      if (notificationData.campaign_id && notificationData.subscription_id) {
        fetch("/api/v1/feedback/ping", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            campaign_id: notificationData.campaign_id,
            subscription_id: notificationData.subscription_id,
            status: "CLICKED",
          }),
        })
          .then(() => console.log("[sw.js] Feedback CLICKED отправлен в Go!"))
          .catch((err) =>
            console.error("[sw.js] Ошибка отправки feedback click:", err),
          );
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
