# API Description

## Get all games (/games/all)

**TO DO:** Resolve how to manage data chunks

Request: empty or filter

Chunk structure:

```json
[
    <game_properties>,
    "user_list": int
]
```

## Get profile list (/games/list/<profile_id:int>)

Reuest: empty or filter

Response:

```json
[
    <game_properties>
]
```

## Games filter

A query added to related requests

```json
{
    "filter": int, // 0 - genre, 1 - listed count, 2 - platform
    "ascending" bool
}
```

## Search (/search)

Request:

```json
{
    "query": string // Grand Theft, Red Dead, Final, etc.
}
```

Response:

```json
{
    "games": [
        <game_properties("name", "platforms", "year_released")>
    ]
}
```

## Get profile (/profile/<nickname:str>)

Request: empty

Response:

```json
{
    "profile_info": <profile_info>
}
```

## Update profile info (/profile)

Request: any field of a <profile_info>

Response: updated fields of this profile

## Add game to list (/games/list/add)

Request:

```json
{
    "game_id": int,
    "list_type": int
}
```
