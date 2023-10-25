# timout
Simple app used for taking time in meetings. 

## How to use
Copy [beep-01a.mp3](https://github.com/henrikor/timeout/blob/master/beep-01a.mp3) and [cow.mp3](https://github.com/henrikor/timeout/blob/master/cow.mp3) to a folder (maybe same folder as the app itself). Run in terminal/cmd in same folder as the mp3 files.

Example of use: 

timeout.exe 1:30

(this will run for 1 minute and 30 seconds)

For more/ better use: Copy beep-01a.mp3, cow.mp3 amd timeout app to ie: ~/timeout and add the following in your .zshrc or whatever:

```
export TIMEOUT_PATH=${HOME}/timeout
export PATH="${PATH}:${HOME}/timeout"
```
