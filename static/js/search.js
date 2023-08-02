const getData = async (searchType, search, order, limit) => {
    const response = await fetch(`http://localhost:8000/api/${searchType}?search=${search}&sort=${order}&limit=${limit}&page=1`)
    const content = response.json();
    console.log(content)
    return content;
}

let displayDiv = document.getElementById("displayDiv")

const display = async () => {
    let searchType = document.getElementById("search-type").value
    let search = await document.getElementById("search-book").value;
    let order = await document.getElementById("order").value;
    let limit = await document.getElementById("limit").value;
    const payload = await getData(searchType, search, order, limit);
    let data = payload.data;
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
    }
}
display();


