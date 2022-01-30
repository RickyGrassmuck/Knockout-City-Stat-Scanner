# Knockout-City-Stat-Scanner

## Credits

This project is mostly a fancy wrapper around the [Extract Table](https://extract-table.com/) [(github)](https://github.com/vegarsti/extract-table) API, they did all the heavy lifting here and deserve all the credit!

## Usage

1. Build with `go build -o stats-scanner scanner.go`
1. Download match stats screen screengrabs from a KoC match to the project folder
1. Run `./stats-scanner`
1. Profit

## Caveats

I have no clue what kind of backend limits thg Extract Table deals with so I have built in a 2 second delay between parsing requests. Please don't change this and hammer the API, in no world does anyone need to parse KoC Match stats any faster than 1 match every 2 seconds.
