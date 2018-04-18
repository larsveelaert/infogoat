# InfoGoat
An easy-to-use vulnerable webapp which showcases different issues and how they can be prevented

# Usage

The application is writen in NodeJS and can be started by cloning the repo and running;
```
node goat.js
```
or use the pre-packaged binary, available for all platforms

A webserver will be available on localhost:3000

If you want to demo the application to a group of people in the same network add the demo (only available if running with node)
```
node goat.js demo
```

# Packaging / Building
Use the following software :https://github.com/zeit/pkg
sudo npm install -g pkg
pkg --targets linux,macos,win goat.js --out-path binaries

