# Data structure

## Game properties

```json
{
    "games": [
        {
            "name": string,
            "platform": string,
            "list_type": int, // 0 - unlisted, 1 - want to play, 2 - playing, 3 - played
            "image_url": string
        }
    ]
}
```

## Profile info

```json
{
    "profile_info": {
        "nickname": string,
        "description": string,
        "games_listed": int,
        "socials": [
            {
                "type": string, // discord, skype, twitter, twitch, youtube, etc.
                "data": string // username, email, url, etc.
            }
        ]
    }
}
```
