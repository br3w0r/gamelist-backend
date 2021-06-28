# Data structure

## Game properties

```json
{
    "name": string,
    "platforms": [
        <platform>
    ],
    "year_released": int,
    "image_url": string,
    "genres": [
        <genre>
    ]
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

```json
{
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
