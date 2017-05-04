# Go Chat!
### A Websocket Experiment

## Welcome
Hi All! I decided I want to try my hand at building a PWA [(Progressive Web App)](https://developers.google.com/web/fundamentals/getting-started/codelabs/your-first-pwapp/ "Your First Progressive Web App")

## Project Dependencies / Technologies
#### Backend
First of all, if you don't have go installed, what are you doing!?! [Go and get go installed!](https://golang.org/doc/install)

I'm trying to keep my third party libraries to a minimum, but I don't want to write the raw duplex connection implementation quite yet.
Therefore, my intention is to use gorilla's implementation of websockets, so you will need to run a:
```
go get github.com/gorilla/websocket
```

#### Frontend
There a number of ways to get the front end dependencies rolling. You can do so with a simple npm and bower install. [npm](https://www.npmjs.com/) [bower](https://bower.io/).

Once You have npm and bower rockin and ready to go, go ahead and run:

```
npm install
bower install
npm run build
```

_(The run build command runs the webpack build that will generate the code you need to get the front end tied together.)_


## Getting Started
Once you have all of the dependencies rolling, you will be able to get the surver running with a simple go run main.go


