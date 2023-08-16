const addReadList = document.getElementById("addReadList")
const removeReadList = document.getElementById("removeReadList")
const user_id = document.getElementById("user_id")
const book_id = document.getElementById("book_id")
const csrfToken = document.querySelector('meta[name="csrf_token"]').getAttribute('content')
const host = window.location.host
const protocol = window.location.protocol

const ReadListExists = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/read`, {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
        }
    })
    const data = await response.json()
    let exists = data.exists
    if (exists === true) {
        addReadList.classList.add("d-none")
    } else {
        removeReadList.classList.add("d-none")
    }

}
ReadListExists()

const addToReadList = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/read`,{
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

const removeFromReadList = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/read`,{
        method: 'DELETE',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

addReadList.addEventListener("click", () => {
    addToReadList()
    addReadList.classList.add("d-none")
    removeReadList.classList.remove("d-none")
})

removeReadList.addEventListener("click", () => {
    removeFromReadList()
    addReadList.classList.remove("d-none")
    removeReadList.classList.add("d-none")
})


const addBuyList = document.getElementById("addBuyList")
const removeBuyList = document.getElementById("removeBuyList")

const BuyListExists = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/buy`, {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
        }
    })
    const data = await response.json()
    let exists = data.exists
    if (exists === true) {
        addBuyList.classList.add("d-none")
    } else {
        removeBuyList.classList.add("d-none")
    }

}
BuyListExists()

const addToBuyList = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/buy`,{
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

const removeFromBuyList = async () => {
    const response = await fetch(`${protocol}//${host}/api/books/${book_id.value}/buy`,{
        method: 'DELETE',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

addBuyList.addEventListener("click", () => {
    addToBuyList()
    addBuyList.classList.add("d-none")
    removeBuyList.classList.remove("d-none")
})

removeBuyList.addEventListener("click", () => {
    removeFromBuyList()
    addBuyList.classList.remove("d-none")
    removeBuyList.classList.add("d-none")
})