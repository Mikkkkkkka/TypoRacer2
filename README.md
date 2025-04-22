# TypoRacer2

A competitive typing game in the early stages of development written in Go.
[Версия ну русском](j./)

## Structure

This project is a mono-repo for both its server-side application and a minimalistic cli-client. The server is intended to be easily started.
The server-side and the client transfer data via HTTP.

## Getting started

If you have docker installed then running the server should be a breeze.

1. Download the repo.
2. Run `docker-compose up --build`
3. The server will be available on port 8080 by default. To change it, modify the corresponding field in the `docker-compose.yml` file.
4. To stop the server run `docker-compose down` and add the `--volumes` flag to remove all the data from the database.

Running the cli-client however, requires having the golang compiler installed.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
