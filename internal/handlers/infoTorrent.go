package handlers

import (
	"fmt"
	"math"

	"github.com/anacrolix/torrent/metainfo"
)

type TreeNode struct {
	Name     string
	Size     string
	IsDir    bool
	Children map[string]*TreeNode
}

func formatSize(size int64) string {
	if size == 0 {
		return "0B"
	}

	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	i := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))
	if i >= len(sizes) {
		i = len(sizes) - 1
	}

	value := float64(size) / math.Pow(1024, float64(i))
	unit := sizes[i]

	if i == 0 || value == math.Trunc(value) {
		return fmt.Sprintf("%d%s", int64(value), unit)
	} else if value < 10 {
		return fmt.Sprintf("%.2f%s", value, unit)
	} else {
		return fmt.Sprintf("%.1f%s", value, unit)
	}
}

func buildTree(info *metainfo.Info) *TreeNode {
	root := &TreeNode{
		Name:     info.Name,
		Size:     formatSize(info.TotalLength()),
		IsDir:    info.IsDir(),
		Children: make(map[string]*TreeNode),
	}

	if !root.IsDir {
		return root
	}

	for _, file := range info.Files {
		current := root
		for i, part := range file.Path {
			if child, exists := current.Children[part]; exists {
				current = child
				continue
			}

			newNode := &TreeNode{
				Name:     part,
				Size:     formatSize(file.Length),
				IsDir:    i < len(file.Path)-1,
				Children: make(map[string]*TreeNode),
			}
			if newNode.IsDir {
				newNode.Size = ""
			}

			current.Children[part] = newNode
			current = newNode
		}
	}
	return root
}
