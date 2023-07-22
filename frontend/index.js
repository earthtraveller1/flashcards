async function getCards() {
    let response

    let res = await fetch("/api/cardstacks")
    response = await res.json()

    return response
}

function main() {
    getCards().then(cards => {
        let stacksListElement = document.getElementById("stack-list")

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

main()
