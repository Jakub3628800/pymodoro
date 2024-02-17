# pomodoro

Pomodoro is a simple timer that helps you to stay focused on your work. It is a CLI tool to start a timer for 25 minutes and then take a 5 minutes break. After 4 pomodoros, it takes a longer break of 15 minutes.
Usage:
```
pomodoro --help
```



## Install from release

```bash
RELEASE_TAG="v0.0.2"
TARGET_DIR="~/.local/bin"
wget https://github.com/Jakub3628800/pomodoro/releases/download/${RELEASE_TAG}/pomodoro.zip pomodoro.zip
unzip pomodoro.zip
mv pomodoro ${TARGET_DIR}
rm pomodoro.zip
```