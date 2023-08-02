async function getStack() {
    const response = await fetch(`/api/cardstacks/${serverInfo.stackName}`)
    const stack = response.json()
    return stack
}

function main() {
    getStack().then((stack) => {
        let mainTitle = document.getElementById("main-title")
        mainTitle.innerText = stack.name

        let mainDescription = document.getElementById("main-description")
        mainDescription.innerText = stack.description
    })
}

main()
