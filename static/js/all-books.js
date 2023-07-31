const getBooks = async (search, order, limit) => {
    const response = await fetch(`http://localhost:8000/api/books?search=${search}&sort=${order}&limit=${limit}&page=1`)
    const content = response.json();
    return content;
}

let booksDisplay = document.getElementById("books-display")

const displayBooks = async () => {
    let search = await document.getElementById("search-book").value;
    let order = await document.getElementById("order").value;
    let limit = await document.getElementById("limit").value;
    const payload = await getBooks(search, order, limit);
    let data = payload.data;
    let books = data.books;
    let bookDisplay = books.map((obj) => {
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
    booksDisplay.innerHTML = bookDisplay;
}
displayBooks();


