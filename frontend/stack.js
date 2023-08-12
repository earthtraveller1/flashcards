let globalStack;

/** @type Node */
let mainPage;
/** @type Node */
let cardsPage;

async function getStack() {
    const response = await fetch(`/api/cardstacks/${serverInfo.stackName}`)
    const stack = response.json()
    return stack
}

function initCardsPage() {
    let returnToMainButton = document.getElementById("return-to-main-button")
    returnToMainButton.onclick = () => {
        cardsPage.parentElement.append(mainPage)
        cardsPage.parentElement.removeChild(cardsPage)
        initMainPage()
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

    cardsPage.parentElement.removeChild(cardsPage)

    initMainPage()
}

main()
