# Confirmerator History

Connect to ethereum node parse the transactions out of the box and store them in postgres for quick lookup of historical data.


### Quickstart

Configuration is done through environment variables.

Database connection by parameters:
* `DBHOST="127.0.0.1"`
* `DBPORT="5432"`
* `DBUSER="postgres"`
* `DBPASS="superSecret!"`
* `DBNAME="confirmerator"`

Ethereum node connection parameters
* `ETHURL=/home/treesloth/.ethereum/geth.ipc;`

## Docker

TODO

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to confirmerator, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 
 
 ## License
 
  The Confirmerator binaries (i.e. all code inside of this project) is licensed under the
 [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
 in our repository in the `LICENSE` file.