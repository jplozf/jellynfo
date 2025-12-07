# Jellynfo

Just a command-line tool to see if anyone is connected to your Jellyfin server.

Upon first launch, you will be asked for the server URL and an API key for authentication; this API key must be generated from the server, in the `Dashboard` menu, then `API keys`. This information will subsequently be saved in a `config.json` file located in a hidden `.jellynfo` subfolder in your user's home folder.

```
~/Projets/Go/jplozf/jellynfo>./jellynfo
Successfully connected to Jellyfin server.
Server Name: paulsboutique
Version: 10.11.2
Operating System: (not available)

--- Connected Users ---
Number of active  sessions: 4
Number of playing sessions: 1

User            Device                    Client                    IP Address      Last Activity             Now Playing
----------------------------------------------------------------------------------------------------------------------------------
mca             paulsboutique             Jellyfin Media Player     192.168.1.42    2025-12-06 16:54:46       Sabotage (Movie)
adrock          Galaxy Tab de adrock      Jellyfin Android          142.142.42.242  2025-12-05 22:51:11       Nothing
miked           iPhone de miked           Jellyfin iOS              42.242.142.42   2025-12-05 21:10:33       Nothing
mixmastermike   iPad (2)                  Jellyfin iPadOS           242.142.42.42   2025-12-05 17:10:19       Nothing

~/Projets/Go/jplozf/jellynfo> 
```
