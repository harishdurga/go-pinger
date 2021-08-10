# go-pinger

A simple Go app to ping websites.
This is a simple Go application to ping websites for fun 😀😉.

### How it works

There is a CSV file `ping_urls.csv` from which the application will read the list of urls to ping. For each url we will create a goroutine. After that, we will wait to receive channel messages from the goroutines.
Once we receive all the messages for the urls we will bulk insert the results into MongoDB a single document for each of the response of a url.

### How to run

- Clone the repo

```bash
git clone https://github.com/harishdurga/go-pinger
```

- Change directory

```bash
cd go-pinger
```

- Install dependencies

```bash
go install
```

- Set MongoDB URI. Ex: mongodb+srv://username:password@cluster0.bqers.mongodb.net/pinger?retryWrites=true&w=majority

```bash
export PINGER_MONGO_CONNECT_URI=
```

- Run

```bash
go run main.go
```
