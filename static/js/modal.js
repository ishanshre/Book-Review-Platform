// open modal by id
function openModal(id) {
    if (id === "following") {
        following()
    }
        document.getElementById(id).classList.add('open');
        document.body.classList.add('jw-modal-open');
}

// close currently open modal
function closeModal() {
    document.querySelector('.jw-modal.open').classList.remove('open');
    document.body.classList.remove('jw-modal-open');
}

window.addEventListener('load', function () {
    // close modals on background click
    document.addEventListener('click', event => {
        if (event.target.classList.contains('jw-modal')) {
            closeModal();
        }
    });
});


const getList = async (listType) => {
    const response = await fetch(`http://${window.location.host}/profile/${listType}`)
    const content = await response.json()
    return content
}


const following = async () => {
    const payload = await getList("followings")
    const followingDis = document.getElementById("following-display")
    let authors = payload.data
    let displayItems = authors.map((obj) => {
        const {id, first_name, last_name} = obj;
        return `
            <div class="d-dark text-white m-5 b-radius text-center">
                <p><a href="/authors/${id}">${first_name} ${last_name}</a></p>
            </div>
        `
    }).join("");
    followingDis.innerHTML = displayItems
}