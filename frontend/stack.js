let globalStack;

let mainPage;
let cardsPage;

async function getStack() {
    const response = await fetch(`/api/cardstacks/${serverInfo.stackName}`)
    const stack = response.json()
    return stack
}

function initMainPage() {
    getStack().then((stack) => {
        globalStack = stack

        let mainTitle = document.getElementById("main-title")
        mainTitle.innerText = stack.name

        let mainDescription = document.getElementById("main-description")
        mainDescription.innerText = stack.description
    })
}

function main() {
    mainPage = document.getElementById("main")
    cardsPage = document.getElementById("cards-page")

    cardsPage.parentElement.removeChild(cardsPage)

    initMainPage()
}

main()
