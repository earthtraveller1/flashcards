async function getStacks() {
    let response

    let res = await fetch("/api/cardstacks")
    response = await res.json()

    return response
}

function initMainPage() {
    getStacks().then(stacks => {
        if (Object.keys(stacks).length == 0) {
            return
        }

        let stacksListElement = document.getElementById("stack-list")
        stacksListElement.innerHTML = ""

        for (let stackId in stacks) {
            let stack = stacks[stackId]

            let button = document.createElement("div")
            button.classList.add("stack-button")
            button.onclick = () => {
                window.location.href = `${location.origin}/stack/${stackId}`
            }

            let title = document.createElement("h1")
            title.innerText = stack.name

            let description = document.createElement("p")
            description.innerText = stack.description

            button.appendChild(title)
            button.appendChild(description)

            stacksListElement.appendChild(button)
        }
    })
}

/**
 * @param {Node} mainPage
 * @param {Node} createStackPage
 */
function initCreateStackPage(mainPage, createStackPage) {
    let button = document.getElementById("create-stack-submit-button")
    let cancelButton = document.getElementById("create-stack-cancel-button")
    let stackName = document.getElementById("new-stack-name")
    let stackDescription = document.getElementById("new-stack-description")

    stackName.value = ""
    stackDescription.value = ""

    button.onclick = () => {
        let name = stackName.value
        let description = stackDescription.value

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

    cancelButton.onclick = () => {
        createStackPage.parentElement.appendChild(mainPage)
        createStackPage.parentElement.removeChild(createStackPage)
        initMainPage()
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
