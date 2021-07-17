# Data structure

## Game properties

```json
{
    "id": int,
    "name": string,
    "platforms": [ // This field currently doesn't work
        <platform>
    ],
    "year_released": int,
    "image_url": string,
    "genres": [ // This field currently doesn't work
        <genre>
    ]
}
```

## Typed game properties

```json
{
    "id": int,
    "name": string,
    "year_released": int,
    "image_url": string,
    "user_list": int // List type of game (0 - Unlisted, 1 - Playing, etc.)
}
```

## Platform

```json
{
    "name": string // PC, Wii, Gamecube, etc.
}
```

## Genre

```json
{
    "name": string // RPG, Action, Rouge-like, Survival, Adventure, etc.
}
```

## Profile info

Currently it's used only for profile creation.

```json
{
    "id": int,
    "nickname": string,
    "description": string,
    "games_listed": int,
    "socials": [
        {
            "type": int, // discord=0, skype=1, twitter=2, twitch=3, youtube=4, etc.
            "data": string // username, email, url, etc.
        }
    ]
}
```
