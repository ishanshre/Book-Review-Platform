const toogleBtn = document.querySelector(".toogle_btn")
const toogleBtnIcon = document.querySelector(".toogle_btn i")
const dropDownMenu = document.querySelector(".dropdown_menu")

toogleBtn.onclick = function () {
    dropDownMenu.classList.toggle('open')
    const isOpen = dropDownMenu.classList.contains("open")
    toogleBtnIcon.classList = isOpen 
        ? 'fa-solid fa-xmark'
        : 'fa-solid fa-bars'
}

function showDeleteForm(element, id) {
    const form = document.getElementById(id)
    element.classList.add("d-none")
    form.classList.remove("d-none")
}

function removeForm(id) {
    const form = document.getElementById(id)
    form.classList.add("d-none")
    const deleteButton = document.getElementById("delete-"+id)
    deleteButton.classList.remove("d-none")
}