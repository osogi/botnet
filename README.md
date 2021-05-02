# botnet
This little botnet was created for my school project. And I promised to post its sources.
# Preparation (windows)
1) Download and install golang compiler
2) 
```
cd %GOPATH%/src
mkdir github.com
cd github.com
mkdir osogi
cd osogi
git clone https://github.com/osogi/botnet
```
3) Choose interesting files, patch it as you want and enjoy life

# Files
**botnet.go** - a package with some functions for communicating between binaries

**botnet_files/server.go** - a program of management server

**botnet_files/client.go** - a program of bot

**botnet_files/controler.go** - a program of client for controlling the bots
