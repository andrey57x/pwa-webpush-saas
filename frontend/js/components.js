const Components = {
  Navbar(user) {
    if (!user) {
      return `
          <header class="container">
            <nav>
              <ul><li><strong>PushEngine SaaS</strong></li></ul>
              <ul>
                <li><a href="#" onclick="App.navigate('login')">Вход</a></li>
                <li><a href="#" onclick="App.navigate('register')">Регистрация</a></li>
              </ul>
            </nav>
          </header>
        `;
    }

    return `
        <header class="container">
          <nav>
            <ul>
              <li><strong>PushEngine SaaS</strong></li>
              <li><small style="color:var(--pico-muted-color);">${user.email}</small></li>
            </ul>
            <ul>
              <li><a href="#" onclick="App.navigate('apps')">Сайты</a></li>
              <li><a href="#" onclick="App.navigate('campaigns')">Рассылки</a></li>
              <li><a href="#" onclick="App.navigate('analytics')">Аналитика</a></li>
              <li><a href="#" onclick="App.navigate('apikeys')">API-Ключи</a></li>
              <li><a href="#" onclick="App.navigate('docs')">Инструкция</a></li>
              <li><a href="#" onclick="App.handleLogout()" style="color:var(--pico-del-color);">Выйти</a></li>
            </ul>
          </nav>
        </header>
      `;
  },

  LoginPage() {
    return `
        <main class="container">
          <article>
            <header><strong>Вход в Личный Кабинет</strong></header>
            <form onsubmit="App.handleLogin(event)">
              <label>Email
                <input type="email" id="login-email" placeholder="admin@domain.com" required>
              </label>
              <label>Пароль
                <input type="password" id="login-pass" placeholder="••••••••" required>
              </label>
              <button type="submit">Войти в панель</button>
            </form>
          </article>
        </main>
      `;
  },

  RegisterPage() {
    return `
        <main class="container">
          <article>
            <header><strong>Регистрация Компании</strong></header>
            <form onsubmit="App.handleRegister(event)">
              <label>Email аккаунта
                <input type="email" id="reg-email" placeholder="manager@pizza.com" required>
              </label>
              <label>Пароль
                <input type="password" id="reg-pass" placeholder="Минимум 8 символов" required>
              </label>
              <button type="submit">Зарегистрировать компанию</button>
            </form>
          </article>
        </main>
      `;
  },

  AppsPage() {
    return `
        <main class="container">
          <article>
            <div style="display:flex; justify-content:space-between; align-items:center;">
              <div>
                <h3 style="margin:0;">Подключенные Сайты</h3>
                <p style="margin:0; color:var(--pico-muted-color);">Каждый сайт имеет уникальный код и автосгенерированные ключи шифрования VAPID.</p>
              </div>
              <button onclick="App.openCreateAppModal()" style="width:auto;">+ Добавить сайт</button>
            </div>
  
            <figure style="margin-top:20px;">
              <table>
                <thead>
                  <tr>
                    <th>Название сайта</th>
                    <th>Код сайта</th>
                    <th>Публичный VAPID-ключ</th>
                  </tr>
                </thead>
                <tbody id="apps-list">
                  <tr><td colspan="3">Загрузка сайтов...</td></tr>
                </tbody>
              </table>
            </figure>
          </article>
  
          <dialog id="create-app-modal">
            <article>
              <header>
                <button aria-label="Close" rel="prev" onclick="App.closeCreateAppModal()"></button>
                <strong>Подключение Нового Сайта</strong>
              </header>
              <form onsubmit="App.handleCreateAppForm(event)">
                <label>Название сайта
                  <input type="text" id="app-name-input" placeholder="Luigi Pizza" required>
                </label>
                <label>Код сайта (англ. буквы и подчёркивание)
                  <input type="text" id="app-code-input" placeholder="pizza_app" required>
                </label>
                <footer>
                  <button type="button" class="secondary" onclick="App.closeCreateAppModal()">Отмена</button>
                  <button type="submit">Создать сайт</button>
                </footer>
              </form>
            </article>
          </dialog>
        </main>
      `;
  },

  CampaignsPage() {
    return `
        <main class="container">
          <article>
            <div style="display:flex; justify-content:space-between; align-items:center;">
              <div>
                <h3 style="margin:0;">Управление Рассылками</h3>
                <p style="margin:0; color:var(--pico-muted-color);">Запускайте push-уведомления по всем активным подписчикам выбранного сайта.</p>
              </div>
              <button onclick="App.openCreateCampaignModal()" style="width:auto;">+ Создать рассылку</button>
            </div>
  
            <div style="margin-top:20px;">
              <label>Фильтр по сайту:
                <select id="camp-app-select"></select>
              </label>
            </div>
  
            <figure style="margin-top:15px;">
              <table>
                <thead>
                  <tr>
                    <th>Заголовок</th>
                    <th>Текст</th>
                    <th>Статус</th>
                    <th>Действия</th>
                  </tr>
                </thead>
                <tbody id="campaigns-list">
                  <tr><td colspan="4">Загрузка рассылок...</td></tr>
                </tbody>
              </table>
            </figure>
          </article>
  
          <dialog id="create-campaign-modal">
            <article>
              <header>
                <button aria-label="Close" rel="prev" onclick="App.closeCreateCampaignModal()"></button>
                <strong>Новая Push-Рассылка</strong>
              </header>
              <form onsubmit="App.handleCreateCampaign(event)">
                <label>Выберите сайт
                  <select id="modal-camp-app-select" required></select>
                </label>
                <label>Заголовок уведомления
                  <input type="text" id="camp-title" placeholder="Скидка 20% на пиццу!" required>
                </label>
                <label>Текст сообщения
                  <textarea id="camp-body" rows="2" placeholder="Введите текст рассылки..." required></textarea>
                </label>
                <label>Ссылка перехода при клике
                  <input type="url" id="camp-url" placeholder="https://pizza.com/promo">
                </label>
                <footer>
                  <button type="button" class="secondary" onclick="App.closeCreateCampaignModal()">Отмена</button>
                  <button type="submit">Сохранить рассылку</button>
                </footer>
              </form>
            </article>
          </dialog>
        </main>
      `;
  },

  AnalyticsPage() {
    return `
        <main class="container">
          <article>
            <header><strong>Аналитика Доставки и CTR (%)</strong></header>
            <div class="grid">
              <label>Выберите сайт
                <select id="analytics-app-select" onchange="App.handleAnalyticsAppChange()"></select>
              </label>
              <label>Выберите рассылку
                <select id="analytics-camp-select" onchange="App.handleLoadAnalytics()"></select>
              </label>
            </div>
          </article>
    
          <article id="analytics-details" style="display:none;">
            <header><strong id="analytics-title">Метрики рассылки</strong></header>
            <table>
              <thead>
                <tr>
                  <th>Целевых подписчиков</th>
                  <th>Отправлено</th>
                  <th>Доставлено</th>
                  <th>Клики</th>
                  <th>CTR (%)</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td id="m-targeted">0</td>
                  <td id="m-sent">0</td>
                  <td id="m-delivered">0</td>
                  <td id="m-clicked">0</td>
                  <td id="m-ctr"><strong>0%</strong></td>
                </tr>
              </tbody>
            </table>
            <div style="height: 260px;">
              <canvas id="analyticsChart"></canvas>
            </div>
          </article>
        </main>
      `;
  },

  APIKeysPage() {
    return `
        <main class="container">
          <article>
            <header><strong>API-Ключи Интеграции</strong></header>
            <p>Используются для отправки рассылок с вашего внешнего бэкенда через заголовок <code>X-API-Key</code>.</p>
            <div class="grid">
              <input type="text" id="key-name-input" placeholder="Название (например: Backend Pizza)">
              <button onclick="App.handleCreateAPIKey()">Выпустить API-Ключ</button>
            </div>
            <div id="raw-key-display"></div>
    
            <figure style="margin-top:20px;">
              <table>
                <thead>
                  <tr>
                    <th>Название</th>
                    <th>Хэш ключа</th>
                    <th>Создан</th>
                  </tr>
                </thead>
                <tbody id="apikeys-list">
                  <tr><td colspan="3">Загрузка ключей...</td></tr>
                </tbody>
              </table>
            </figure>
          </article>
        </main>
      `;
  },

  DocsPage() {
    return `
        <main class="container">
          <article>
            <header><strong>Инструкция по Интеграции</strong></header>
            
            <h3>1. Установка Service Worker на ваш сайт</h3>
            <p>Создайте файл <code>sw.js</code> в корневой директории вашего сайта (например, <code>https://your-site.com/sw.js</code>):</p>
            <pre><code>// sw.js (Скопируйте 1 в 1)
  self.addEventListener("install", function (event) {
    self.skipWaiting();
  });
  
  self.addEventListener("activate", function (event) {
    event.waitUntil(self.clients.claim());
  });
  
  self.addEventListener("push", function (event) {
    let payload = { title: "Уведомление", body: "" };
    if (event.data) {
      try { payload = event.data.json(); } catch (e) { payload = { title: "Уведомление", body: event.data.text() }; }
    }
  
    const saasUrl = "http://localhost:8080/api/v1"; // Замените на URL вашей SaaS системы
  
    event.waitUntil(
      self.registration.showNotification(payload.title || "Уведомление", {
        body: payload.body || "",
        data: payload
      }).then(function () {
        if (payload.campaign_id && payload.subscription_id) {
          fetch(saasUrl + "/feedback/ping", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ campaign_id: payload.campaign_id, subscription_id: payload.subscription_id, status: "DELIVERED" })
          });
        }
      })
    );
  });
  
  self.addEventListener("notificationclick", function (event) {
    event.notification.close();
    const data = event.notification.data || {};
    const saasUrl = "http://localhost:8080/api/v1";
  
    event.waitUntil(
      clients.matchAll({ type: "window" }).then(function (clientList) {
        if (data.campaign_id && data.subscription_id) {
          fetch(saasUrl + "/feedback/ping", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ campaign_id: data.campaign_id, subscription_id: data.subscription_id, status: "CLICKED" })
          });
        }
        if (clients.openWindow) return clients.openWindow(data.url || "/");
      })
    );
  });</code></pre>
  
            <h3>2. Получение VAPID ключа вашего сайта</h3>
            <p>Выполните GET-запрос к API платформы по коду вашего сайта:</p>
            <pre><code>GET /api/v1/apps/{app_code}</code></pre>
  
            <h3>3. Отправка рассылок с вашего бэкенда</h3>
            <p>Выполните POST-запрос с использованием выписанного <code>X-API-Key</code>:</p>
            <pre><code>POST /api/v1/campaigns/{campaign_id}/send
  Headers:
    X-API-Key: pk_live_your_secret_key_here</code></pre>
          </article>
        </main>
      `;
  },
};
