(function () {
  const root = document.getElementById("active-fast");
  if (!root) return;

  const timerEl = document.getElementById("elapsed-timer");
  const percentEl = document.getElementById("fast-percent");
  const barEl = document.getElementById("progress-bar");

  const startTime = new Date(root.dataset.startTime);
  const goalSeconds = Number(root.dataset.goal);
  const totalSeconds = Math.max(goalSeconds, 1);

  function pad(n) {
    return String(n).padStart(2, "0");
  }

  function formatElapsed(seconds) {
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    return `${pad(hrs)}:${pad(mins)}:${pad(secs)}`;
  }

  function tick() {
    const now = new Date();
    const elapsedSeconds = Math.max(
      0,
      Math.floor((now - startTime) / 1000),
    );
    const ratio = Math.min(
      elapsedSeconds / totalSeconds,
      1,
    );
    const percent = Math.floor(ratio * 100);

    timerEl.textContent = formatElapsed(elapsedSeconds);
    percentEl.textContent = `${percent}%`;
    barEl.style.width = `${percent}%`;
    barEl.textContent = `${percent}%`;
    barEl.parentElement.setAttribute(
      "aria-valuenow",
      String(percent),
    );
  }

  tick();
  setInterval(tick, 1000);
})();
