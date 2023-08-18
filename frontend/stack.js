class Card {
    /** @type string */
    front;

    /** @type string */
    back;
}

class CardStack {
    /** @type string */
    name;

    /** @type string */
    description;

    /** @type Card[] */
    cards;
}

/** @type CardStack */
let globalStack;

/** @type Node */
let mainPage;
/** @type Node */
let cardsPage;
/** @type Node */
let addCardPage;

/** @returns CardStack */
async function getStack() {
    const response = await fetch(`/api/cardstacks/${serverInfo.stackName}`)
    /** @type CardStack */
    const stack = await response.json()
    return stack
}

function initCreateCardPage() {
    // Clear the input fields
    let newCardFrontInput = document.getElementById("new-card-front")
    let newCardBackInput = document.getElementById("new-card-back")

    newCardFrontInput.value = ""
    newCardBackInput.value = ""

    let cancelButton = document.getElementById("create-card-cancel-button")
    cancelButton.onclick = () => {
        addCardPage.parentElement.appendChild(cardsPage)
        addCardPage.parentElement.removeChild(addCardPage)
        initCardsPage()
    }
}

function initCardsPage() {
    let returnToMainButton = document.getElementById("return-to-main-button")
    returnToMainButton.onclick = () => {
        cardsPage.parentElement.append(mainPage)
        cardsPage.parentElement.removeChild(cardsPage)
        initMainPage()
    }

    let addCardButton = document.getElementById("add-card-button")
    addCardButton.onclick = () => {
        cardsPage.parentElement.appendChild(addCardPage)
        cardsPage.parentElement.removeChild(cardsPage)
        initCreateCardPage()
    }
}

function initMainPage() {
    let manageCardsButton = document.getElementById("manage-cards-button")
    manageCardsButton.onclick = () => {
        mainPage.parentElement.append(cardsPage)
        mainPage.parentElement.removeChild(mainPage)
        initCardsPage()
    }

    getStack().then((stack) => {
        globalStack = stack

        let mainTitle = document.getElementById("main-title")
        mainTitle.innerText = stack.name

        let pageTitle = document.getElementById("page-title")
        pageTitle.innerText = `${stack.name} - Flashcards`

        let mainDescription = document.getElementById("main-description")
        mainDescription.innerText = stack.description

        let deleteButton = document.getElementById("delete-button")
        deleteButton.onclick = () => {
            fetch(`/api/cardstacks/${serverInfo.stackName}`, { method: "DELETE" })
                .then(() => {
                    window.location.href = `${location.origin}`
                })
                .catch(console.error)
        }

    })
}

function main() {
    mainPage = document.getElementById("main")
    cardsPage = document.getElementById("cards-page")
    addCardPage = document.getElementById("add-card-page")

    cardsPage.parentElement.removeChild(cardsPage)
    addCardPage.parentElement.removeChild(addCardPage)

    initMainPage()
}

main()
