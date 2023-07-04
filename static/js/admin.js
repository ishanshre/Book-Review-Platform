// Get the sidebar toggle button and sidebar element
const sidebarToggleBtn = document.querySelector('.sidebar-toggle-btn');
const sidebar = document.querySelector('.nav-sidebar');
const navLinks = document.querySelector(".nav-links")

// Add click event listener to the toggle button
sidebarToggleBtn.addEventListener('click', () => {
  sidebar.classList.toggle('open');
  navLinks.classList.toggle('open')
});


function showDeleteForm(element, id) {
    const form = document.getElementById(id)
    element.classList.add("no-display")
    form.classList.remove("no-display")
}

function removeForm(id) {
    const form = document.getElementById(id)
    form.classList.add("no-display")
    const deleteButton = document.getElementById("delete-"+id)
    deleteButton.classList.remove("no-display")
}