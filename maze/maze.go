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
	{-1,0}, {0,-1}, {1,0}, {0,1},
}

func (p point)add(r point) point{
		return point {p.i+r.i,p.j+r.j}
}

func (p point) at (grid [][]int )(int,bool){
	if p.i <0 || p.i >= len(grid) {
		return 0,false
	}

	if p.j < 0 ||p.j >= len(grid[p.i]){
		return 0,false
	}
	return grid[p.i][p.j],true
}

func walk(maze [][]int, start, end point) [][]int{
	steps := make([][]int ,len(maze))
	for i := range steps {
		steps[i] = make ([]int , len(maze[i]))
	}

	Q := []point{start}

	for len(Q) > 0 {
		cur := Q[0]
		Q = Q[1:]

		if cur == end {
			break
		}

		for _, dir := range dirs {
			next := cur.add(dir)

			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}

			val, ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}

			if next == start {
				continue
			}

			curSteps, _ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1

			Q = append(Q, next)

		}

	}
	return steps
}

func getPath(steps [][]int, start,end point) [][]int {
	path := make ([][]int , len(steps))
	for i := range path {
		path[i] = make([]int, len(steps[i]))
	}

	Q := []point{end}

	endNum := steps[end.i][end.j]
	curNum := endNum

	for len(Q) > 0 {
		cur := Q[0]
		Q = Q[1:]

		if cur == start {
			break
		}

		for _,dir := range dirs{
			next := cur.add(dir)
			val,ok := next.at(steps)
			if !ok || val == endNum {
				continue
			}

			if val == curNum -1 {
				curNum = curNum -1

				path[cur.i][cur.j] = steps[cur.i][cur.j]
				Q = append(Q, next)
				break
			}

		}

	}
	return path
}

func main(){
	//读取迷宫定义
	maze :=	readMaze("maze/maze.in")
	for i := range maze	{
		for j := range maze[i]{
			fmt.Print(maze[i][j]," ")
		}
		fmt.Println()
	}

	//计算迷宫路径
	fmt.Println()
	steps := walk(maze, point{0,0}, point{len(maze) - 1, len(maze[0]) - 1})
	for _, row := range steps {
		for _ , val := range row {
			fmt.Printf("%3d",val)
		}
		fmt.Println()
	}

	//去除走不通的路，只显示正确路径
	fmt.Println()
	path := getPath(steps, point{0,0}, point {len(steps)-1,len(steps[0])-1})
	for _, row := range path {
		for _ , val := range row {
			fmt.Printf("%3d",val)
		}
		fmt.Println()

	}

}
