# Ant-farm


## Description

A digital version of an ant farm.

Program `ant-farm` will read from a file (describing the ants and the colony) given in the arguments.

Upon successfully finding the quickest path, `ant-farm` will display the content of the file passed as argument and each move the ants make from room to room.

- At the beginning of the game, all the ants are in the room `##start`. The goal is to bring them to the room `##end` with as few moves as possible.

### Tools and methods
* Ant farm project was approached from `BFS` perspective, as it gives advantaged in case best case scenario comes to life and you don't have to discover all paths to find the best combination.
* For reading incoming data used-> `regex` 
* Data about rooms and paths stored in 2D integer array for efficiency porposes
* For path searching and BFS implementation used `recursive` functions


### Restrictions

- Only the [standard Go](https://golang.org/pkg/) packages are allowed.


## Executing the Program

Use the folowing command, replacing the tag with file name from `examples` folder
```
go run . "<example00.txt>"
```

There is a test file created, check how ants find they way in different maps.

```
bash test.sh
```
