# log-error-detector
Register error related keywords in config file in advance. It can detect and notify when an error message is written in the log in real time.
## Usage
### Install
```
go get -u -v "github.com/fsnotify/fsnotify"
go get -u -v "github.com/go-sql-driver/mysql"
go build
cd web
npm install --save
```
### Configuration
```
vim conf.json
```
```
{
    // Register keywords to be detected
    "patterns": [
        "FATAL",
        "ERROR",
        "DEBUG"
    ],
    // Register log files path
    "logfiles": [
        "sompath/log/log1.log",
        "sompath/log/log2.log",
        "sompath/log/log3.log",
        "sompath/log/log4.log"
    ],
    // The number of lines to read when reading the file
    "maxreadline": 50,
    // Use Multi Core
    "ismulticore": true
}
```
### Run
```
./log-error-detector
npm start
```
## UI 
You can use mysql-live-select in nodejs to push the log history to the screen when an event occurs. At the same time as the history is saved, the client or third party can give notification.

* Ref: [mysql-live-select](https://github.com/numtel/mysql-live-select) - Provide events on `SELECT` statement result set updates

![ScreenShot](/data/log-error-detector-ui.png)