# h4spie2017
bluetooth lock


Install dependencies:
`go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/raspi`
`go get -d -u github.com/paypal/gatt...`


Build: 
`GOARM=7 GOARCH=arm GOOS=linux go build main.go`