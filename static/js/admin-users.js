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