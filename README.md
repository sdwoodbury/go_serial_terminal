
### purpose
Really none. If you're using Linux, you're better off to program in C and configure the serial port with termios. If using Windows, just use Putty. Probably the same for Mac. 

### usage 
* go\_serial\_terminal port baud\_rate
* ctrl+c or the like to stop

### dependencies
* go get github.com/albenik/go-serial

### compile and cross-compile 
* go build 
* GOOS=linux GOARCH=arm go build 


