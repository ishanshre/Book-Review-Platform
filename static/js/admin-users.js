function showDeleteForm(element, id) {
    const form = document.getElementById(id)
    element.classList.add("deleteForm")
    form.classList.remove("deleteForm")
}

function removeForm(id) {
    const form = document.getElementById(id)
    form.classList.add("deleteForm")
    const deleteButton = document.getElementById("delete-"+id)
    deleteButton.classList.remove("deleteForm")
}