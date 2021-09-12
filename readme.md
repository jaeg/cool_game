## Setup
This project used ebiten and installation instructions can be found here: https://ebiten.org/documents/install.html

If you are on Ubuntu chances are all you need is:

`sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`


After installing ebiten dependencies you'll need to run the following command to get all the go deps.

`make vendor`

## Running/Building
`make run`

`make build`

To test with a local copy of the game-engine library add this to the mod file replacing the path/to/game_engine with the location of your copy:

`replace github.com/jaeg/game-engine => /path/to/game_engine`