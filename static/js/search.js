const host = window.location.host
let lastPage;
let currentPage = 1;
const csrfToken = document.querySelector('meta[name="csrf_token"]').getAttribute('content')
let totalItems;
const paginationNumbers = document.getElementById("pagination-numbers")
const nextButton = document.getElementById("next-button")
const prevButton = document.getElementById("prev-button")
const getData = async (searchType, search, order, limit) => {
    let url = window.location.pathname
    let parts = url.split("/")
    if (parts[1] === "genres") {
        const response = await fetch(`http://${host}/api/genres?search=${search}&sort=${order}&limit=${limit}&page=${currentPage}&genre=${parts[2]}`)
        const content = response.json();
        return content;
    } else if (parts[1] === "languages") {
        const response = await fetch(`http://${host}/api/languages?search=${search}&sort=${order}&limit=${limit}&page=${currentPage}&language=${parts[2]}`)
        const content = response.json();
        return content;
    } else if (parts[1] === "read-list") {
        const response = await fetch(`http://${host}/api/read-list?search=${search}&sort=${order}&limit=${limit}&page=${currentPage}`)
        const content = response.json();
        return content;
    } else if (parts[1] === "buy-list") {
        const response = await fetch(`http://${host}/api/buy-list?search=${search}&sort=${order}&limit=${limit}&page=${currentPage}`)
        const content = response.json();
        return content;
    } else {
        const response = await fetch(`http://${host}/api/${searchType}?search=${search}&sort=${order}&limit=${limit}&page=${currentPage}`)
        const content = response.json();
        return content;
    }
}

let displayDiv = document.getElementById("displayDiv")

const display = async () => {
    let searchType = document.getElementById("search-type").value
    let search = await document.getElementById("search-book").value;
    let order = await document.getElementById("order").value;
    let limit = await document.getElementById("limit").value;
    const payload = await getData(searchType, search, order, limit);
    let data = payload.data;
    console.log(data)
    lastPage = data.last_page
    currentPage = data.page
    totalItems = data.total
    if (searchType === "books") {
        let books = data.books;
        let displayItems = books.map((obj) => {
            const { title, isbn, cover } = obj;
            return `
        <div class="card-box">
            <a href="/books/${isbn}">
                <div class="card-img">
                    <img src="/${cover}" alt="img-{{$book.Title}}">
                </div>

                <div class="card-text">
                    <span>${title}</span>
                </div>
            </a>
        </div>
        `
        }).join("");
        displayDiv.innerHTML = displayItems;
    } else if (searchType === "authors") {
        let authors = data.authors;
        let displayItems = authors.map((obj) => {
            const { id, first_name, last_name, avatar } = obj
            return `
            <div class="card-author-box d-flex justify-center">
                <a href="/authors/${id}">
                    <div class="card-author-img">
                        <img src="/${avatar}" alt="img-{{$book.Title}}">
                    </div>

                    <div class="card-text text-center">
                        <span>${first_name} ${last_name}</span>
                    </div>
                </a>
            </div>
            `
        }).join("");
        displayDiv.innerHTML = displayItems;
    } else if (searchType === "admin-users") {
        let users = data.users;
        let displayItems = users.map((obj)=> {
            const { id, username, access_level, created_at, is_validated} = obj
            return `
                <tr>
                    <td>${id}</td>
                    <td>${username}</td>
                    <td>${access_level}</td>
                    <td>${created_at}</td>
                    <td>${is_validated}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/users/detail/${id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${id}')" /></button>

                            <div class="jw-modal" id="delete-${id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/users/detail/${id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete @${username}?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-publishers") {
        let publishers = data.publishers;
        let displayItems = publishers.map((obj)=> {
            const { id, name, established_date } = obj
            return `
                <tr>
                    <td>${id}</td>
                    <td>${name}</td>
                    <td>${established_date}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/publishers/detail/${id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${id}')" /></button>

                            <div class="jw-modal" id="delete-${id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/publishers/detail/${id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this publisher ${name}?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-authors") {
        let authors = data.authors;
        let displayItems = authors.map((obj)=> {
            const { id, first_name, last_name } = obj
            return `
                <tr>
                    <td>${id}</td>
                    <td>${first_name} ${last_name}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/authors/detail/${id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${id}')" /></button>

                            <div class="jw-modal" id="delete-${id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/authors/detail/${id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this author ${first_name} ${last_name}?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-books") {
        let books = data.books;
        let displayItems = books.map((obj)=> {
            const { id, title, is_active, added_at } = obj
            return `
                <tr>
                    <td>${id}</td>
                    <td>${title}</td>
                    <td>${is_active}</td>
                    <td>${added_at}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/books/detail/${id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${id}')" /></button>

                            <div class="jw-modal" id="delete-${id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/books/detail/${id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this book ${title}?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-bookauthors") {
        let bookAuthors = data.book_authors;
        let displayItems = bookAuthors.map((obj)=> {
            const { book_id, book_title, author_id, author_first_name, author_last_name } = obj
            return `
                <tr>
                    <td>${book_title}</td>
                    <td>${author_first_name} ${author_last_name}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/bookAuthors/detail/${book_id}/${author_id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${book_id}-${author_id}')" /></button>

                            <div class="jw-modal" id="delete-${book_id}-${author_id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/bookAuthors/detail/${book_id}/${author_id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this relationship?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-readlists") {
        let readLists = data.read_lists;
        let displayItems = readLists.map((obj)=> {
            const { user_id, username, book_id, book_title, created_at } = obj
            return `
                <tr>
                    <td>${book_title}</td>
                    <td>${username}</td>
                    <td>${created_at}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/readLists/detail/${book_id}/${user_id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${book_id}-${user_id}')" /></button>

                            <div class="jw-modal" id="delete-${book_id}-${user_id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/readLists/detail/${book_id}/${user_id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this relationship?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-buylists") {
        let buyLists = data.buy_lists;
        let displayItems = buyLists.map((obj)=> {
            const { user_id, username, book_id, book_title, created_at } = obj
            return `
                <tr>
                    <td>${book_title}</td>
                    <td>${username}</td>
                    <td>${created_at}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/buyLists/detail/${book_id}/${user_id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${book_id}-${user_id}')" /></button>

                            <div class="jw-modal" id="delete-${book_id}-${user_id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/buyLists/detail/${book_id}/${user_id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this relationship?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-followers") {
        let followers = data.followers;
        let displayItems = followers.map((obj)=> {
            const { user_id, username, author_id, author_first_name, author_last_name, followed_at } = obj
            return `
                <tr>
                    <td>${author_first_name} ${author_last_name}</td>
                    <td>${username}</td>
                    <td>${followed_at}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/followers/detail/${author_id}/${user_id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${author_id}-${user_id}')" /></button>

                            <div class="jw-modal" id="delete-${author_id}-${user_id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/followers/detail/${author_id}/${user_id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this relationship?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    } else if (searchType === "admin-reviews") {
        let reviews = data.reviews;
        let displayItems = reviews.map((obj)=> {
            const { id, rating, username, title, created_at, updated_at } = obj
            return `
                <tr>
                    <td>${id}</td>
                    <td>${rating}</td>
                    <td>${username}</td>
                    <td>${title}</td>
                    <td>${created_at}</td>
                    <td>${updated_at}</td>
                    <td>
                        <div class="action-icons">
                            <button><a href="/admin/reviews/detail/${id}"><img src="/static/images/edit-icon.png" alt="update-icon"/></a></button>
                            <button ><img width="19px" height="19px" src="/static/images/del-icon.png" alt="del-icon" onclick="openModal('delete-${id}')" /></button>

                            <div class="jw-modal" id="delete-${id}">
                                <div class="jw-modal-body">
                                    <form action="/admin/reviews/detail/${id}/delete" method="post">
                                        <input type="hidden" name="csrf_token" id="csrf_token" value="${csrfToken}">
                                        <p>Do you want to delete this relationship?</p>
                                        <input type="submit" value="Delete Record">
                                        <button type="button" onclick="closeModal()">No</button>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </td>
                </tr>
            `
        }).join("")
        displayDiv.innerHTML = displayItems
    }
    paginationNumbers.innerHTML = ''
    const getPaginationNumbers = () => {
        for (let i = 1; i <= lastPage; i++) {
            appendPageNumber(i, currentPage)
        }
    }
    getPaginationNumbers()    
}
display();

const appendPageNumber = (index, current) => {
    const pageNumber = document.createElement("button")
    pageNumber.className = "pagination-number";
    pageNumber.innerHTML = index;
    pageNumber.setAttribute("page-index", index);
    pageNumber.setAttribute("aria-label", "Page" + index);
    pageNumber.onclick = () => {
        pageClicked(index);
    }
    if (index === current) {
        pageNumber.classList.add('pg-active')
    }
    paginationNumbers.appendChild(pageNumber)
}

const pageClicked = (index) => {
    console.log("page clicked:", index)
    currentPage = index
    display()
}