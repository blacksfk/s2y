# Spotify to Youtube discord bot
Blast your music all curated in your spotify playlists via the rythm bot in discord!

## Setup
1. Create a copy of `config.json.example` and rename the copy to `config.json`. Alternatively you can name it whatever and pass the absolute path via the config flag.
2. Edit the `config.json` with your API tokens and keys.
3. Change the mode so only the current user/group has read permission (i.e. 400/440).
4. Compile with `go build`.
5. Run with `./s2y` and optionally provide the configuration file (if you have named it differently or it isn't located in the execution directory) with `--config /absolute/path/to/config.json`.

# License
BSD-3-clause
