package main

import "fmt"

type TreeNode struct {
	HasToy bool
	Left *TreeNode
	Right *TreeNode
}
func addToy() *TreeNode {
    return &TreeNode{HasToy: true, Left: nil, Right: nil}
}
func addNoToy() *TreeNode {
    return &TreeNode{HasToy: false, Left: nil, Right: nil}
}

// ex00
func areToysBalanced(root *TreeNode) bool {
    if root == nil {
        return true
    }
    leftToys := countToys(root.Left)
    rightToys := countToys(root.Right)
    return leftToys == rightToys
}

func countToys(node *TreeNode) int {
    if node == nil {
        return 0
    }
    count := 0
    if node.HasToy {
        count = 1
    }
    return count + countToys(node.Left) + countToys(node.Right)
}

// ex01
func treeHeight(root *TreeNode) int{ 
    if root == nil {
		return 0
	}
    lHeight := treeHeight(root.Left);
    rHeight := treeHeight(root.Right);
    return max(lHeight, rHeight) + 1;
}

func leftToRightTrav(root *TreeNode, level int, ans *[]bool){
    if root == nil {
        return
    }

    if level == 1 {
        *ans = append(*ans, root.HasToy)
    } else {
        leftToRightTrav(root.Left, level-1, ans)
        leftToRightTrav(root.Right, level-1, ans)
    }
}

func rightToLeftTrav(root *TreeNode, level int, ans *[]bool) {
    if root == nil {
        return
    }

    if level == 1 {
        *ans = append(*ans, root.HasToy)
    } else {
        rightToLeftTrav(root.Right, level-1, ans)
        rightToLeftTrav(root.Left, level-1, ans)
    }
}

func unrollGarland(root *TreeNode) []bool{
	var ans []bool
	leftToRight := false
	height := treeHeight(root)

	for i := 1; i <= height; i++ {
		if leftToRight {
			leftToRightTrav(root, i, &ans)
		} else {
			rightToLeftTrav(root, i, &ans)
		}

		leftToRight = !leftToRight
	}

	return ans
}


func main() {
    root := &TreeNode{
        HasToy: false,
        Left: &TreeNode{
            HasToy: true,				
            Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
            Right:  &TreeNode{HasToy: false, Left: nil, Right: nil},
        },
        Right: &TreeNode{
            HasToy: true,
            Left:   &TreeNode{HasToy: false, Left: nil, Right: nil},
            Right:  &TreeNode{HasToy: false, Left: nil, Right: nil},
        },
    }
	root.Left.Left.Left = addNoToy()
	root.Left.Left.Right = addNoToy()
	root.Right.Right.Left = addNoToy()
	root.Right.Right.Right = addToy()

/*
		  0        
	    /   \
	   1     1     слева направо 
	  / \    / \
	 1   0  0   0   справа налево
    / \        / \
   0   0      0   1   слева направо 
*/

	// ex00
    // fmt.Println(areToysBalanced(root)) 

	// ex01
	ans := unrollGarland(root) 
	fmt.Println(ans)
}