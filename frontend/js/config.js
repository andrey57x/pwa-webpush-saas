window.APP_CONFIG = {
  // Автоматическое определение протокола и хоста текущего окружения
  getApiBaseUrl() {
    return `${window.location.origin}/api/v1`;
  },
};
