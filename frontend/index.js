async function getCards() {
    let response

    let res = await fetch("/api/cardstacks")
    response = await res.json()

    return response
}

function initMainPage() {
    getCards().then(cards => {
        let stacksListElement = document.getElementById("stack-list")
        stacksListElement.innerHTML = ""

        cards.forEach((stack, _i, _stacks) => {
            let button = document.createElement("div")
            button.classList.add("stack-button")

            let title = document.createElement("h1")
            title.innerText = stack.Name

            let description = document.createElement("p")
            description.innerText = stack.Description

            button.appendChild(title)
            button.appendChild(description)

            stacksListElement.appendChild(button)
        })
    })
}

function initCreateStackPage(mainPage, createStackPage) {
    let button = document.getElementById("create-stack-submit-button")
    button.onclick = () => {
        let name = document.getElementById("new-stack-name").value
        let description = document.getElementById("new-stack-description").value

        const payload = {
            name: name,
            description: description
        }

        fetch("/api/cardstacks", {
            method: "POST", 
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payload)
        }).then(() => {
            createStackPage.parentElement.append(mainPage)
            createStackPage.remove()
            initMainPage()
        }).catch(console.error)
    }
}

function main() {
    let createStackPart = document.getElementById("create-stack-part")
    let mainPage = document.getElementById("main")

    createStackPart.remove()

    let createStackButton = document.getElementById("button-create-stacks")
    createStackButton.onclick = () => {
        mainPage.parentElement.append(createStackPart)
        mainPage.remove()
        initCreateStackPage(mainPage, createStackPart)
    }

    initMainPage()
}

main()
