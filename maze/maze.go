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

var dirs = [4]point {
	{-1,0}, {0,-1}, {1,0}, {0,1}
}

func (p point)add(r point) point{
		return point {p.i+r.i,p.j+r.j}
}
func walk(maze [][]int, start, end point){
	steps := make([][]int ,len(maze))
	for i := range steps {
		steps[i] = make ([]int , len(maze[i]))
	}

	Q := []point{start}

	for len(Q) > 0 {
		cur := Q[0]
		Q = Q[1:]
		for _, dir := range dirs {
			next := cur.add(dir)
		}

	}
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
