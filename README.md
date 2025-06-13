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

### Media Sources

Thank you to all creators of the included free images and sounds.

* For image sources and attributions see `resources/images/sources/sources.md`
* For sound sources and attributions see `resources/sounds/sources.md`

### Tools Used

* [Go](https://go.dev/) as the general programming language.
* [Ebitengine](https://ebitengine.org/) for building a 2D game.
* [GIMP](https://www.gimp.org) for extracting ingredient images from source pictures.
* [ffmpeg](https://ffmpeg.org/) for adapting audio properties like sample size and volume.

### Improvement Ideas

* Consider an intro screen that explains how to play the game.
* Improve logging, e.g. use `slog` for more structured logging, revisit and introduce logging with proper log levels and configuration.
* Simplify implementation of `Game.move` method. Avoid recursive call for handling lane change not possible.
* Improve audio control for playing sounds, seems like actual volume is not influenced by system volume, at least on Linux.
* Add automated tests where easily possible.
* Consider switching from an all ingredients slice to a map structure to support lookups (ordering is not important).
* Improve randomness of burgers, e.g. exclude mix of meat and veggie patties, avoid too many repetitive etc.
* Check screen resizing and behavior on screens with smaller / larger sizes.
* Improve layout, separate build from order section more clearly.
* Track points, consider adding a timer.