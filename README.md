# Burger Game

### About the Game

Create tasty burgers by piling up falling ingredients. Pay attention to specific orders as customers don't like unexpected or missing ingredients. Serve them well and you will become a successful burger shop.

A simple 2D game built with Golang and Ebitengine.

### How to Run

* Ensure [Go is installed on your system](https://go.dev/doc/install)
* Download and extract or `git clone` this repositories content to your local machine
* Navigate into your local repository directory (e.g. in the terminal) and run `go run main.go`

### How to Play

* Use the left and right arrow key to move falling ingredients to another lane.
* Use the down arrow key or simply wait to ensure an ingredient is stacked on top of a burger.
* Press digit keys 1, 2 or 3 to sell the burgers 1, 2 or 3 respectively.

### Troubleshooting

#### No Audio

You should hear some audio effects while playing the game, e.g. when selling a burger. If you don't hear any sounds while playing, check your audio output device and volume, make sure it is not muted. If it is still not working, and you are running on Linux, you may need to apply subsequent workaround to disable a problematic audio module. This worked for me, but please use it carefully at your own risk and revert it in case of any issues:
* Open following config file for editing: `/etc/modprobe.d/alsa-base.conf`
* Append an option to disable the possibly problematic module: `options snd-hda-intel model=auto blacklist snd_soc_avs`

### Improvement Ideas

* Improve logging, e.g. use `slog` for more structured logging, revisit and introduce logging with proper log levels and configuration.