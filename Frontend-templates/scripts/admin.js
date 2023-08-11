// Get the sidebar toggle button and sidebar element
const sidebarToggleBtn = document.querySelector('.sidebar-toggle-btn');
const sidebar = document.querySelector('.nav-sidebar');
const navLinks = document.querySelector(".nav-links")

// Add click event listener to the toggle button
sidebarToggleBtn.addEventListener('click', () => {
  sidebar.classList.toggle('open');
  navLinks.classList.toggle('open')
});
