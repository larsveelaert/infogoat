const http = require('http')
const process = require('process')
const port = 3000
console.log(process.argv[2])
const requestHandler = (request, response) => {
	  console.log(request.url)
	  if (request.url == "/examples/profile") {
		  response.end('profile');
		 }
	  response.end('Hello Node.js Server!')
}

const server = http.createServer(requestHandler)

server.on('connection', function(sock) {
	  console.log('Client connected from ' + sock.remoteAddress);
	  // Client address at time of connection ----^
});
server.listen(port, '127.0.0.1', (err) => {
	  if (err) {
		      return console.log('something bad happened', err)
		    }
	  console.log(`server is listening on ${port}`)
})
