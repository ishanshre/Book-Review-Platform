const follow = document.getElementById("follow")
const unfollow = document.getElementById("unfollow")
const user_id = document.getElementById("user_id")
const author_id = document.getElementById("author_id")
const csrfToken = document.querySelector('meta[name="csrf_token"]').getAttribute('content')
const host = window.location.host

const followExists = async () => {
    const response = await fetch(`http://${host}/api/authors/${author_id.value}/exists`)
    const data = await response.json()
    let exists = data.exists
    if (exists === true) {
        follow.classList.add("d-none")
    } else {
        unfollow.classList.add("d-none")
    }

}
followExists()

const followAuthor = async () => {
    const response = await fetch(`http://${host}/api/authors/${author_id.value}/follow`,{
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

const unfollowAuthor = async () => {
    const response = await fetch(`http://${host}/api/authors/${author_id.value}/unfollow`,{
        method: 'DELETE',
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": csrfToken,
        }
    })
}

follow.addEventListener("click", () => {
    followAuthor()
    follow.classList.add("d-none")
    unfollow.classList.remove("d-none")
})

unfollow.addEventListener("click", () => {
    unfollowAuthor()
    unfollow.classList.add("d-none")
    follow.classList.remove("d-none")
})