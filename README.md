# Confirmerator

Listen to blockchain networks such as Bitcoin and Ethereum for transactions on given addresses and publish messages to NATS to be handled for by other programs such as a notification pusher.

Currently functioning as intended for Ethereum but is with a Work-in-Progress.  


## Quickstart

Configuration is done through environment variables.

Configure ethereum connection:

`ETHURL`

`ETHWSURL`

Configure the database connection:

`DBHOST`

`DBPORT`

`DBUSER`

`DBPASS`

`DBNAME`

Add NATS url:

`NATSURL`


### TODO list

* Connect to bitcoin zmq
* Monitor smart-contracts
* Add suppoprt for other databases

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to confirmeratorRest, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 
 
 ## License
 
 Copyright (c) 2018 Michael Salmons
 
 The confirmerator binaries (i.e. all code inside of this project) is licensed under the
 [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
 in our repository in the `LICENSE` file.