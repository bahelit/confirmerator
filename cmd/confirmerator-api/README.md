# Confirmerator API

Provides REST API to the Confirmerator project. Add/Remove users and devices


### Quickstart

Configuration is done through environment variables.

Configure the http listener with `LISTEN-ADDRESS` and `LISTEN-PORT`

Database connection by string:
`MONGOURI`  example  `mongodb0.example.com:27017/my_db`

Another `MONGOURI` example with credentials and replica set
`myDBReader:D1fficultP%40ssw0rd@mongodb0.example.com:27017,mongodb1.example.com:27017,mongodb2.example.com:27017/my_db?replicaSet=myRepl`

## Docker

`docker build -f api/Dockerfile -t confirmerator-api . && docker run -p 80:80 --name confirmerator-api confirmerator-api`

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
 
  The Confirmerator binaries (i.e. all code inside of this project) is licensed under the
 [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
 in our repository in the `LICENSE` file.