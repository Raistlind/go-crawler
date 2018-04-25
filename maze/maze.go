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
	fmt.Println("row is ",row, " col is " , col)
	maze := make([][]int,row)
	for i := range maze {
		maze[i] = make ([]int , col)
		for j := range maze[i]{
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}
	return maze
}

type point struct{
	i,j int
}

func walk(maze [][]int, start, end point){

}

func main(){
	maze :=	readMaze("maze/maze.in")
	for i := range maze	{
		for j := range maze[i]{
			fmt.Print(maze[i][j]," ")
		}
		fmt.Println()
	}
	walk(maze, point{0,0}, point{len(maze) - 1, len(maze[0]) - 1})
}
