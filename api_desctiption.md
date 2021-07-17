# API Description

All api requests have this url structure: `<address>/api/v0/<API>` where `<API>` is a route of api action (in brackets bellow).

## Authorization

Nearly all of requests require authorization. Authorization header structure:

`Authorization: Bearer <token>`

List of routes that don't require authorization:

- POST /profiles
- /aquire-tokens
- /refresh-tokens
- /revoke-token
- not in production mode: all POST and GET requests for genres, platforms, social types and GET methods for profiles and game properties.

## [POST] Get all games (/games/all)

Request:

```json
{
    "last": int, // ID of last game entry on client
    "batch_size": int // amount of games to be sent to client from server
}
```

Response: list of `<typed_game_properties>` with size of `<batch_size>`, but not more than 10.

## [GET] Get games of authorized user (/my-games)

Response: list of `<typed_game_properties>`

## [POST] Add game to list (/list-game)

Request:

```json
{
    "game_id": int,
    "list_type": int
}
```

## [POST] Search (/games/search)

Request:

```json
{
    "name": string // Grand Theft, Red Dead, Final, etc.
}
```

Response:

```json
{
    "games": [
        <game_properties("id", "name")>
    ]
}
```

## [POST] Get game details (/games/details)

Request:

```json
{
    "id": int
}
```

Response:

```json
{
    "game": <typed_game_properties>,
    "platforms": [
        <platform>
    ],
    "genres": [
        <genre>
    ]
}
```

## [POST] Sign Up (/profiles)

Request: `profile_info` data structure

## [POST] Aquire new JWT tokens (/aquire-tokens)

Request:

```json
{
    "nickname": string,
    "email": string,
    "password": string
}
```

It's required that at least email or password where in the request, but not necessarily both of them.

Response:

```json
{
    "token": string,
    "refresh_token": string
}
```

## [POST] Refresh tokens (/refresh-tokens)

Request:

```json
{
    "refresh_token": string
}
```

Response:

```json
{
    "token": string,
    "refresh_token": string
}
```

## [POST] Revoke refresh token (/revoke-token)

Request:

```json
{
    "refresh_token": string
}
```

## [GET] Delete all refresh tokens (/delete-all-refresh-tokens)

## [POST] Add game (/games)

## [GET][POST] List types (/list-types)

## [GET][POST] Genres (/genres)

## [GET][POST] Platforms (/platforms)

## [GET][POST] Social Types (/social-types)

## [NOT IMPLEMENTED] Get profile list (/games/list/<profile_id:int>)

Reuest: empty or filter

Response:

```json
[
    <game_properties>
]
```

## [NOT IMPLEMENTED] Games filter

A query added to related requests

```json
{
    "filter": int, // 0 - genre, 1 - listed count, 2 - platform
    "ascending" bool
}
```

## [NOT IMPLEMENTED] Get profile (/profiles/<nickname:str>)

Request: empty

Response:

```json
{
    "profile_info": <profile_info>
}
```

## [NOT IMPLEMENTED] Update profile info (/profiles)

Request: any field of a <profile_info>

Response: updated fields of this profile
