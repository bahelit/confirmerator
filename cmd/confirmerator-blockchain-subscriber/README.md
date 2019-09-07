# Confirmerator

Listen to blockchain networks such as Bitcoin and Ethereum for transactions on given addresses 
and publish messages to NATS to be handled for by other programs such as a notification pusher.

Confirmerator provides users with the ability receive notifications 
when their transactions are confirmed.

Currently functioning as intended for Ethereum but is with a Work-in-Progress.  

## Quickstart

Configuration is done through environment variables.

Configure ethereum connection:

`ETHURL`

`ETHWSURL`

Database connection by string:
 
`MONGOURI`  example  `mongodb0.example.com:27017/my_db`

Another `MONGOURI` example with credentials and replica set
`myDBReader:D1fficultP%40ssw0rd@mongodb0.example.com:27017,mongodb1.example.com:27017,mongodb2.example.com:27017/my_db?replicaSet=myRepl`

Add NATS url:

`NATSURL`

## Docker

Example docker command
`$ docker build -f api/Dockerfile -t confirmerator . && docker run -p 80:80 --name confirmerator confirmerator`

### TODO list

* Connect to bitcoin zmq
* Monitor smart-contracts
* Respond to confirmation count request

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to Confirmerator, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 
 
 ## License
 
 The Confirmerator binaries (i.e. all code inside of this project) is licensed under the
 [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
 in our repository in the `LICENSE` file.