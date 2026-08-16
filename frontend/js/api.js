const API = {
  getToken() {
    return localStorage.getItem("saas_jwt_token") || "";
  },

  setToken(token) {
    localStorage.setItem("saas_jwt_token", token);
  },

  async request(endpoint, options = {}) {
    const headers = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    const token = this.getToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const url = `${window.APP_CONFIG.getApiBaseUrl()}${endpoint}`;

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      throw new Error(data.error || `HTTP Error ${response.status}`);
    }

    return data;
  },

  register(email, password) {
    return this.request("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  login(email, password) {
    return this.request("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  getAppByCode(appCode) {
    return this.request(`/apps/${appCode}`);
  },

  createApp(name, appCode) {
    return this.request("/apps", {
      method: "POST",
      body: JSON.stringify({ name, app_code: appCode }),
    });
  },

  subscribePWA(data) {
    return this.request("/subscribe", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  createCampaign(data) {
    return this.request("/campaigns", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  sendCampaign(campaignId) {
    return this.request(`/campaigns/${campaignId}/send`, {
      method: "POST",
    });
  },

  getCampaignStats(campaignId) {
    return this.request(`/analytics/campaigns/${campaignId}`);
  },

  getApps() {
    return this.request("/apps");
  },

  getCampaigns(appId) {
    return this.request(`/campaigns?app_id=${appId}`);
  },

  getAPIKeys() {
    return this.request("/api-keys");
  },

  createAPIKey(name) {
    return this.request("/api-keys", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
  },
};
