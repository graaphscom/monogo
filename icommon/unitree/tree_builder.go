package unitree

import (
	"fmt"
	"github.com/graaphscom/monogo/icommon/json"
	"github.com/graaphscom/monogo/icommon/metadata"
	"github.com/graaphscom/monogo/icommon/tsmakers"
	"os"
	"path"
)

func BuildRootTree(manifest json.IcoManifest) (IconsTree, error) {
	subTrees := make([]IconsTree, 0, len(manifest.VendorsPaths))

	metadataStore := metadata.NewStore(manifest)

	for vendor, vendorPaths := range manifest.VendorsPaths {
		vendorTreeBuilder, ok := treeBuilders[vendor]
		if !ok {
			return IconsTree{}, fmt.Errorf("tree builder doesn't exist for %s", vendor)
		}
		vendorTree, err := vendorTreeBuilder.buildTree(metadataStore, path.Join(manifest.BasePath, vendorPaths.Icons), vendor)
		if err != nil {
			return IconsTree{}, err
		}

		subTrees = append(subTrees, vendorTree)
	}

	return IconsTree{
		Name:     "icommon",
		SubTrees: &subTrees,
	}, nil
}

func (b categoriesTreeBuilder) buildTree(metadata metadata.Store, src, rootName string) (IconsTree, error) {
	srcRootEntries, err := os.ReadDir(src)

	if err != nil {
		return IconsTree{}, err
	}

	subTrees := make([]IconsTree, 0, len(srcRootEntries))

	for _, srcRootEntry := range srcRootEntries {
		if !srcRootEntry.IsDir() {
			continue
		}

		subTree, err := b.iconsTreeBuilder.buildTree(metadata, path.Join(src, srcRootEntry.Name()), srcRootEntry.Name())

		if err != nil {
			return IconsTree{}, err
		}

		subTrees = append(subTrees, subTree)
	}

	return IconsTree{Name: rootName, SubTrees: &subTrees}, nil
}

type categoriesTreeBuilder struct {
	iconsTreeBuilder treeBuilder
}

func (b iconsTreeBuilder) buildTree(metadata metadata.Store, src, rootName string) (IconsTree, error) {
	srcEntries, err := os.ReadDir(path.Join(src, b.srcSuffix))

	if err != nil {
		return IconsTree{}, err
	}

	icons := make([]Icon, 0, len(srcEntries))

	for _, srcEntry := range srcEntries {
		iconName := srcEntry.Name()
		if b.iconNameConverter != nil {
			iconName = b.iconNameConverter(srcEntry.Name())
		}

		tags, err := b.tagsExtractor(metadata, rootName, srcEntry.Name())

		if err != nil {
			return IconsTree{}, err
		}

		icons = append(icons, Icon{
			Name:    iconName,
			SrcFile: path.Join(src, srcEntry.Name()),
			Tags:    tags,
		})
	}

	treeName := rootName
	if b.treeNameConverter != nil {
		treeName = b.treeNameConverter(rootName)
	}

	return IconsTree{Name: treeName, IconSet: &IconSet{Icons: icons, TsMaker: b.tsMaker}}, nil
}

type iconsTreeBuilder struct {
	iconNameConverter func(in string) string
	treeNameConverter func(in string) string
	srcSuffix         string
	tagsExtractor     func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error)
	tsMaker           tsmakers.Maker
}
