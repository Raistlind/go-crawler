package main

import (
	"os"
	"fmt"
)

func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var row, col int
	fmt.Fscanf(file, "%d %d", &row, &col)

	maze := make([][]int,row)
	for i := range maze {
		maze[i] = make ([]int , col)
		for j := range maze[i]{
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}
}

func main(){
	readMaze("maze/maze.in")
}
