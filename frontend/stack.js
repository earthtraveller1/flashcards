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

/** @type Node */
let noCardsMessage;

/** @returns CardStack */
async function getStack() {
    const response = await fetch(`/api/cardstacks/${serverInfo.stackName}`)
    /** @type CardStack */
    const stack = await response.json()
    return stack
}

/** @param {Card} card */
async function submitCard(card) {
    try {
        await fetch(`/api/cardstacks/${serverInfo.stackName}/cards`, { method: "POST", body: JSON.stringify(card) })
    } catch (e) {
        console.error(e)
    }

    // Switch back to the page.
    addCardPage.parentElement.appendChild(cardsPage)
    addCardPage.parentElement.removeChild(addCardPage)
    initCardsPage()
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

    let createButton = document.getElementById("create-card-button")
    let newCardFront = document.getElementById("new-card-front")
    let newCardBack = document.getElementById("new-card-back")

    createButton.onclick = () => {
        let card = new Card()
        card.front = newCardFront.value
        card.back = newCardBack.value

        submitCard(card).catch(console.error)
    }
}

/**
* @param {Node} cards
    * @param {CardStack} stack
*/
function fillCardsList(cards, stack) {
    cards.innerHTML = ""

    if (stack.cards.length <= 0) {
        cards.appendChild(noCardsMessage)
        return
    }

    stack.cards.forEach((card, index, _cards) => {
        let cardFrontElement = document.createElement("div")
        cardFrontElement.innerText = card.front
        cardFrontElement.classList.add("card-front")

        let separator = document.createElement("hr")

        let cardBackElement = document.createElement("div")
        cardBackElement.innerText = card.back
        cardBackElement.classList.add("card-back")

        let cardElement = document.createElement("div")
        cardElement.appendChild(cardFrontElement)
        cardElement.appendChild(separator)
        cardElement.appendChild(cardBackElement)
        cardElement.classList.add("card")

        cardElement.onclick = () => {
            fetch(`/api/cardstacks/${serverInfo.stackName}/cards/${index}`, { method: "DELETE" })
                .then(() => {
                    return fetch(`/api/cardstacks/${serverInfo.stackName}`)
                })
                .then((response) => {
                    return response.json()
                })
                .then((pStack) => {
                    fillCardsList(cards, pStack)
                })
                .catch(console.error)
        }

        cards.appendChild(cardElement)
    })
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

    if (noCardsMessage == undefined) {
        noCardsMessage = document.getElementById("no-cards-message")
    }

    getStack().then((cardStack) => {
        globalStack = cardStack
    }).catch((reason) => {
        console.error(reason)
    }).finally(() => {
        let cards = document.getElementById("cards")
        fillCardsList(cards, globalStack)
    })
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
