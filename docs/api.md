# API Endpoints

Here are the API endpoints for the server that we can use. This is for myself to reference during development.

## GET Endpoints

- `GET /api/cardstacks`: Obtain a full list of card stacks.
- `GET /api/cardstacks/{stackID}`: Obtain information about a specific stack with the given ID.

## POST Endpoints

- `POST /api/cardstacks`: Create a card stack.
- `POST /api/cardstacks/{stackID}/cards`: Create a card in the given card stack.

## DELETE Endpoints

- `DELETE /api/cardstacks/{stackID}`: Delete the stack with the given ID.
- `DELETE /api/cardstacks/{stackID}/cards/{cardIndex}`: Delete the card in the given stack with the given index.
