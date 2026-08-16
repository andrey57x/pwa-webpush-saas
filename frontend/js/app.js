const App = {
  currentTab: "apps",
  currentUser: null,
  apps: [],

  init() {
    this.restoreSession();
  },

  restoreSession() {
    const token = API.getToken();
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split(".")[1]));
        this.currentUser = {
          id: payload.user_id,
          role: payload.role,
          tenant_id: payload.tenant_id,
          email: "Авторизованный пользователь",
        };

        const savedTab =
          localStorage.getItem("saas_active_tab") ||
          window.location.hash.replace("#", "") ||
          "apps";
        this.navigate(savedTab);
      } catch (e) {
        API.setToken("");
        this.navigate("login");
      }
    } else {
      this.navigate("login");
    }
  },

  navigate(tab) {
    this.currentTab = tab;
    if (this.currentUser) {
      localStorage.setItem("saas_active_tab", tab);
      window.location.hash = tab;
    }
    this.render();

    if (tab === "apps") this.loadAppsData();
    if (tab === "campaigns") this.loadCampaignsData();
    if (tab === "analytics") this.loadAnalyticsData();
    if (tab === "apikeys") this.loadAPIKeysData();
  },

  render() {
    const root = document.getElementById("app");
    if (!root) return;

    let content = "";
    if (!this.currentUser) {
      content =
        this.currentTab === "register"
          ? Components.RegisterPage()
          : Components.LoginPage();
    } else {
      switch (this.currentTab) {
        case "apps":
          content = Components.AppsPage();
          break;
        case "campaigns":
          content = Components.CampaignsPage();
          break;
        case "analytics":
          content = Components.AnalyticsPage();
          break;
        case "apikeys":
          content = Components.APIKeysPage();
          break;
        case "docs":
          content = Components.DocsPage();
          break;
        default:
          content = Components.AppsPage();
          break;
      }
    }

    root.innerHTML = Components.Navbar(this.currentUser) + content;
  },

  async handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById("login-email").value;
    const pass = document.getElementById("login-pass").value;

    try {
      const res = await API.login(email, pass);
      API.setToken(res.access_token);
      Toast.success("Успешный вход!");
      this.restoreSession();
    } catch (err) {
      Toast.error(`Ошибка входа: ${err.message}`);
    }
  },

  async handleRegister(e) {
    e.preventDefault();
    const email = document.getElementById("reg-email").value;
    const pass = document.getElementById("reg-pass").value;

    try {
      await API.register(email, pass);
      Toast.success("Компания зарегистрирована! Выполняем вход...");
      const res = await API.login(email, pass);
      API.setToken(res.access_token);
      this.restoreSession();
    } catch (err) {
      Toast.error(`Ошибка регистрации: ${err.message}`);
    }
  },

  handleLogout() {
    API.setToken("");
    localStorage.removeItem("saas_active_tab");
    window.location.hash = "";
    this.currentUser = null;
    Toast.info("Вы вышли из аккаунта");
    this.navigate("login");
  },

  openCreateAppModal() {
    const modal = document.getElementById("create-app-modal");
    if (modal) modal.showModal();
  },

  closeCreateAppModal() {
    const modal = document.getElementById("create-app-modal");
    if (modal) modal.close();
  },

  openCreateCampaignModal() {
    const modal = document.getElementById("create-campaign-modal");
    if (modal) {
      const modalSelect = document.getElementById("modal-camp-app-select");
      if (modalSelect && this.apps) {
        modalSelect.innerHTML = this.apps
          .map(
            (a) => `<option value="${a.id}">${a.name} (${a.app_code})</option>`,
          )
          .join("");
      }
      modal.showModal();
    }
  },

  closeCreateCampaignModal() {
    const modal = document.getElementById("create-campaign-modal");
    if (modal) modal.close();
  },

  async loadAppsData() {
    try {
      this.apps = await API.getApps();
      const list = document.getElementById("apps-list");
      if (!list) return;

      if (!this.apps || this.apps.length === 0) {
        list.innerHTML =
          '<tr><td colspan="3">Нет подключенных сайтов. Нажмите "+ Добавить сайт" выше.</td></tr>';
        return;
      }

      list.innerHTML = this.apps
        .map(
          (a) => `
            <tr>
              <td><strong>${a.name}</strong></td>
              <td><code>${a.app_code}</code></td>
              <td><small>${a.vapid_public_key.substring(0, 30)}...</small></td>
            </tr>
          `,
        )
        .join("");
    } catch (err) {
      Toast.error(`Ошибка загрузки сайтов: ${err.message}`);
    }
  },

  async handleCreateAppForm(e) {
    e.preventDefault();
    const name = document.getElementById("app-name-input").value;
    const code = document.getElementById("app-code-input").value;

    try {
      const app = await API.createApp(name, code);
      Toast.success(`Сайт "${app.name}" подключен!`);
      this.closeCreateAppModal();
      this.loadAppsData();
    } catch (err) {
      Toast.error(`Ошибка создания сайта: ${err.message}`);
    }
  },

  async loadCampaignsData() {
    this.apps = await API.getApps();
    const select = document.getElementById("camp-app-select");
    if (!select) return;

    if (!this.apps || this.apps.length === 0) {
      select.innerHTML =
        '<option value="">Сначала добавьте сайт во вкладке "Сайты"</option>';
      return;
    }

    select.innerHTML = this.apps
      .map((a) => `<option value="${a.id}">${a.name} (${a.app_code})</option>`)
      .join("");
    this.loadCampaignsList(this.apps[0].id);

    select.onchange = (e) => this.loadCampaignsList(e.target.value);
  },

  async loadCampaignsList(appId) {
    const list = document.getElementById("campaigns-list");
    try {
      const campaigns = await API.getCampaigns(appId);
      if (!campaigns || campaigns.length === 0) {
        list.innerHTML =
          '<tr><td colspan="4">Рассылок пока нет. Нажмите "+ Создать рассылку" выше.</td></tr>';
        return;
      }

      list.innerHTML = campaigns
        .map(
          (c) => `
            <tr>
              <td><strong>${c.title}</strong></td>
              <td>${c.body}</td>
              <td><mark>${c.status}</mark></td>
              <td>
                <button onclick="App.handleSendCampaign('${c.id}')">Запустить</button>
              </td>
            </tr>
          `,
        )
        .join("");
    } catch (err) {
      list.innerHTML = `<tr><td colspan="4">Ошибка: ${err.message}</td></tr>`;
    }
  },

  async handleCreateCampaign(e) {
    e.preventDefault();
    const appId = document.getElementById("modal-camp-app-select").value;
    if (!appId) return Toast.error("Выберите сайт!");

    const title = document.getElementById("camp-title").value;
    const body = document.getElementById("camp-body").value;
    const targetUrl = document.getElementById("camp-url").value;

    try {
      await API.createCampaign({
        app_id: appId,
        title,
        body,
        target_url: targetUrl,
      });
      Toast.success(`Рассылка создана!`);
      this.closeCreateCampaignModal();
      this.loadCampaignsList(appId);
    } catch (err) {
      Toast.error(`Ошибка: ${err.message}`);
    }
  },

  async handleSendCampaign(campaignId) {
    try {
      const res = await API.sendCampaign(campaignId);
      Toast.success(res.message);
      const appId = document.getElementById("camp-app-select").value;
      setTimeout(() => this.loadCampaignsList(appId), 1000);
    } catch (err) {
      Toast.error(`Ошибка отправки: ${err.message}`);
    }
  },

  async loadAnalyticsData() {
    this.apps = await API.getApps();
    const select = document.getElementById("analytics-app-select");
    if (!select) return;

    if (!this.apps || this.apps.length === 0) return;

    select.innerHTML = this.apps
      .map((a) => `<option value="${a.id}">${a.name}</option>`)
      .join("");
    this.handleAnalyticsAppChange();
  },

  async handleAnalyticsAppChange() {
    const appId = document.getElementById("analytics-app-select").value;
    const select = document.getElementById("analytics-camp-select");

    const campaigns = await API.getCampaigns(appId);
    if (!campaigns || campaigns.length === 0) {
      select.innerHTML = '<option value="">Нет рассылок</option>';
      document.getElementById("analytics-details").style.display = "none";
      return;
    }

    select.innerHTML = campaigns
      .map((c) => `<option value="${c.id}">${c.title}</option>`)
      .join("");
    this.handleLoadAnalytics();
  },

  async handleLoadAnalytics() {
    const campaignId = document.getElementById("analytics-camp-select").value;
    if (!campaignId) return;

    try {
      const stats = await API.getCampaignStats(campaignId);
      document.getElementById("analytics-details").style.display = "block";

      document.getElementById("m-targeted").innerText = stats.total_targeted;
      document.getElementById("m-sent").innerText = stats.sent_count;
      document.getElementById("m-delivered").innerText = stats.delivered_count;
      document.getElementById("m-clicked").innerText = stats.clicked_count;
      document.getElementById("m-ctr").innerText = stats.ctr + "%";

      renderMetricsChart("analyticsChart", stats);
    } catch (err) {
      Toast.error(`Ошибка аналитики: ${err.message}`);
    }
  },

  async loadAPIKeysData() {
    try {
      const keys = await API.getAPIKeys();
      const list = document.getElementById("apikeys-list");
      if (!list) return;

      if (!keys || keys.length === 0) {
        list.innerHTML =
          '<tr><td colspan="3">Ключи не выпускались. Создайте первый выше!</td></tr>';
        return;
      }

      list.innerHTML = keys
        .map(
          (k) => `
            <tr>
              <td><strong>${k.name}</strong></td>
              <td><code>${k.key_hash.substring(0, 20)}...</code></td>
              <td><small>${new Date(k.created_at).toLocaleDateString()}</small></td>
            </tr>
          `,
        )
        .join("");
    } catch (err) {
      Toast.error(`Ошибка загрузки ключей: ${err.message}`);
    }
  },

  async handleCreateAPIKey() {
    const name = document.getElementById("key-name-input").value;
    if (!name) return Toast.error("Введите название ключа!");

    try {
      const res = await API.createAPIKey(name);
      document.getElementById("raw-key-display").innerHTML = `
          <ins>Скопируйте ключ сейчас: <code>${res.raw_key}</code></ins>
        `;
      Toast.success("API-Ключ создан!");
      this.loadAPIKeysData();
    } catch (err) {
      Toast.error(`Ошибка: ${err.message}`);
    }
  },
};

document.addEventListener("DOMContentLoaded", () => App.init());
