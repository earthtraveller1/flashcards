async function getCards() {
    let response

    let res = await fetch("/api/cardstacks")
    response = await res.json()

    return response
}

getCards().then(cards => {
    let stacksListElement = document.getElementById("stack-list")

    cards.forEach((stack, i, stacks) => {
        stacksListElement.innerHTML += `<h1>${stack.Name}</h1>`
        stacksListElement.innerHTML += `<p>${stack.Description}</h1>`
    })
})
