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

const buttons = document.querySelectorAll("[data-carousel-button]")

buttons.forEach(button => {
    button.addEventListener("click", () => {
        const offset = button.dataset.carouselButton === "next" ? 1 : -1
        const slides = button
            .closest("[data-carousel]")
            .querySelector("[data-slides]")
        const activeSlide = slides.querySelector("[data-active]")
        let newIndex = [...slides.children].indexOf(activeSlide) + offset
        if (newIndex < 0) newIndex = slides.children.length - 1
        if (newIndex >= slides.children.length) newIndex = 0

        slides.children[newIndex].dataset.active = true
        delete activeSlide.dataset.active
    })
})