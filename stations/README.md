
# Stations pathfinder

## Description 

The program finds all possible paths between the starting station and the ending station in the map. It then takes a given number of trains and sends them along the shortest routes, ensuring that the stations are clear. If the station is busy, the program chooses another path, choosing the next shortest available route. In this way, all the trains get moved from start to end station in a close to optimal manner

This program works if there is contact between stations and there are no other errors described in the task.
The program removes paths that have a conflict between stations (stations have contact in both directions and therefore cannot be free under any conditions).

**Test files** have been created to check **errors**:

`paths_test.go`
`mapping_test.go`
`main_test.go`
`mapfile_test.go`

## Usage program

In the project root folder, run the program with command

```bash

go run . [mapfile name] [start station] [end station] [number of trains]
```

for example:

```Bash

go run . london.map waterloo st_pancras 5
```
## Usage test files

`paths_test.go`
`mapping_test.go`
`main_test.go`
`mapfile_test.go`

1. To check for errors in main.go need to build program file. In the terminal, type:

```go build -o program main.go```

*Don't need to do this for other test files.*

2. Next, click the green arrow
![alt text](<Screenshot 2024-07-23 at 13.36.28.png>)  or ![alt text](<Screenshot 2024-07-23 at 13.36.35.png>)


