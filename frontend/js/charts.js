let metricsChart = null;

function renderMetricsChart(containerId, stats) {
  const canvas = document.getElementById(containerId);
  if (!canvas) return;

  if (metricsChart) {
    metricsChart.destroy();
  }

  const ctx = canvas.getContext("2d");

  metricsChart = new Chart(ctx, {
    type: "bar",
    data: {
      labels: ["Целевые", "Отправлено", "Доставлено", "Клики (CTR)", "Ошибки"],
      datasets: [
        {
          label: "Количество событий",
          data: [
            stats.total_targeted || 0,
            stats.sent_count || 0,
            stats.delivered_count || 0,
            stats.clicked_count || 0,
            stats.failed_count || 0,
          ],
          backgroundColor: [
            "rgba(99, 102, 241, 0.8)",
            "rgba(59, 130, 246, 0.8)",
            "rgba(16, 185, 129, 0.8)",
            "rgba(245, 158, 11, 0.8)",
            "rgba(239, 68, 68, 0.8)",
          ],
          borderRadius: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: "#1f2937",
          titleColor: "#f9fafb",
          bodyColor: "#9ca3af",
          borderColor: "rgba(255,255,255,0.1)",
          borderWidth: 1,
        },
      },
      scales: {
        x: {
          grid: { color: "rgba(255,255,255,0.05)" },
          ticks: { color: "#9ca3af" },
        },
        y: {
          grid: { color: "rgba(255,255,255,0.05)" },
          ticks: { color: "#9ca3af", precision: 0 },
        },
      },
    },
  });
}
