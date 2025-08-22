const $tooltip = document.querySelector("#tooltip");
let circles = [];
let x = 0;
let y = 0;
let downX, downY;
let since, until;
const dist = (c) => Math.abs(x - c.x);
const distY = (c) => Math.abs(y - c.y);
const updateHovered = () => {
  let nearest;
  for (const c of circles) {
    if (
      !nearest ||
      dist(c) < dist(nearest) ||
      (dist(c) === dist(nearest) && distY(c) < distY(nearest))
    ) {
      nearest = c;
    }
  }
  if (nearest) {
    $tooltip.textContent = nearest.text;
    $tooltip.style.left = `${nearest.x - $tooltip.clientWidth / 2}px`;
    $tooltip.style.top = `${nearest.y - 20}px`;
  }
};
const onMouseDown = (e) => {
  downX = e.clientX;
  downY = e.clientY;
};
const onMouseUp = (e) => {
  onDrag();
};
const onMouseMove = (e) => {
  x = e.clientX;
  y = e.clientY;
  updateHovered();
};
const loadSinceUntil = () => {
  until = parseInt(document.querySelector("#until").textContent);
  since = parseInt(document.querySelector("#since").textContent);
};
const saveSinceUntil = () => {
  if (until > Date.now() / 1000) {
    const currentSize = until - since;
    until = undefined;
    since = Date.now() / 1000 - currentSize;
  }
  const q = new URLSearchParams(window.location.search);
  q.set("since", formatDate(since));
  if (until) {
    q.set("until", formatDate(until));
  } else {
    q.delete("until");
  }
  window.history.replaceState(
    null,
    "",
    `${window.location.origin + window.location.pathname}?${q.toString()}`,
  );
  refresh();
};
const onWheel = (e) => {
  loadSinceUntil();
  const zoomCenter = (e.clientX / window.innerWidth - 0.1) * 1.25;
  const currentSize = until - since;
  const zoomCenterTime = since + zoomCenter * currentSize;
  const newSize = Math.round(currentSize * (1 + 0.1 * Math.sign(e.deltaY)));
  since = Math.round(zoomCenterTime - zoomCenter * newSize);
  until = since + newSize;
  saveSinceUntil();
};
const onDrag = (e) => {
  loadSinceUntil();
  const currentSize = until - since;
  const deltaTime = ((downX - x) / window.innerWidth) * 1.25 * currentSize;
  since = Math.round(since + deltaTime);
  until = since + currentSize;
  saveSinceUntil();
};
document.body.addEventListener("mousedown", onMouseDown);
document.body.addEventListener("mouseup", onMouseUp);
document.body.addEventListener("mousemove", onMouseMove);
document.body.addEventListener("wheel", onWheel);

const xx = (x) => {
  return x.toString().padStart(2, "0");
};

const formatDate = (utc) => {
  const d = new Date(utc * 1000);
  return (
    `${d.getUTCFullYear()}-${xx(d.getUTCMonth() + 1)}-${xx(d.getUTCDate())}` +
    `T${xx(d.getUTCHours())}:${xx(d.getUTCMinutes())}:${xx(d.getUTCSeconds())}`
  );
};

const refresh = async () => {
  const q = new URLSearchParams(window.location.search);
  if (since) {
    q.set("since", formatDate(since));
  }
  if (until) {
    q.set("until", formatDate(until));
  }
  const response = await fetch(
    `${window.location.origin + window.location.pathname}.svg${window.location.search}`,
  );
  const text = await response.text();
  document.querySelector("svg").outerHTML = text;
  circles = Array.from(document.querySelectorAll("circle")).map((c) => ({
    x: (c.getBoundingClientRect().left + c.getBoundingClientRect().right) / 2,
    y: c.getBoundingClientRect().y,
    text: c.textContent,
  }));
  updateHovered();
};
refresh();
setInterval(refresh, 1000);
