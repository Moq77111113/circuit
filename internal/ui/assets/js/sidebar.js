function toggleSidebar() {
  const sidebar = document.querySelector(".app__sidebar");
  const overlay = document.querySelector(".mobile-overlay");
  const toggle = document.querySelector(".mobile-menu-toggle");
  sidebar.classList.toggle("is-open");
  overlay.classList.toggle("is-visible");
  toggle.classList.toggle("is-active");
}

function closeSidebar() {
  const sidebar = document.querySelector(".app__sidebar");
  const overlay = document.querySelector(".mobile-overlay");
  const toggle = document.querySelector(".mobile-menu-toggle");
  sidebar.classList.remove("is-open");
  overlay.classList.remove("is-visible");
  toggle.classList.remove("is-active");
}

function toggleNavItem(element) {
  const listItem = element.parentElement;
  listItem.classList.toggle("collapsed");
}
