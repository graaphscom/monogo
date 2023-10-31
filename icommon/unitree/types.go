package unitree

import (
	"github.com/graaphscom/monogo/icommon/metadata"
	"github.com/graaphscom/monogo/icommon/tsmakers"
)

type treeBuilder interface {
	buildTree(metadata metadata.Store, src, rootName string) (IconsTree, error)
}

type IconsTree struct {
	Name     string
	SubTrees *[]IconsTree
	IconSet  *IconSet
}

func (tree IconsTree) Traverse(segments []string, fn func(segments []string, iconSet IconSet)) {
	if tree.IconSet != nil {
		fn(append(segments, tree.Name), *tree.IconSet)
	}
	if tree.SubTrees != nil {
		for _, subTree := range *tree.SubTrees {
			subTree.Traverse(append(segments, tree.Name), fn)
		}
	}
}

type IconSet struct {
	Icons   []Icon
	TsMaker tsmakers.Maker
}

type Icon struct {
	Name    string
	SrcFile string
	Tags    IconTags
}

type IconTags struct {
	Search []string
	Visual []string
}
