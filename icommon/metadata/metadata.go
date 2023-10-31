package metadata

import (
	"github.com/graaphscom/monogo/icommon/json"
	"path"
	"strings"
)

func (S Store) GetFluentui(asset string) (Fluentui, error) {
	if S.metadata.fluentui == nil {
		S.metadata.fluentui = make(map[string]Fluentui)
	}

	if v, ok := S.metadata.fluentui[asset]; ok {
		return v, nil
	}

	result, err := json.ReadJson[Fluentui](
		path.Join(S.manifest.BasePath, S.manifest.VendorsPaths["fluentui"].Icons, asset, "metadata.json"),
	)

	if err != nil {
		return Fluentui{}, err
	}

	S.metadata.fluentui[asset] = result

	return result, nil
}

func (S Store) GetFontawesome() (Fontawesome, error) {
	return singleFile[Fontawesome](&S.metadata.fontawesome, S.manifest, "fontawesome")
}

func (S Store) GetOcticons() (Octicons, error) {
	return singleFile[Octicons](&S.metadata.octicons, S.manifest, "octicons")
}

func (S Store) GetRemixicon() (Remixicon, error) {
	return singleFile[Remixicon](&S.metadata.remixicon, S.manifest, "remixicon")
}

func (S Store) GetUnicons(variant string) (UniconsQuickAccess, error) {
	if S.metadata.unicons == nil {
		S.metadata.unicons = make(map[string]UniconsQuickAccess)
	}

	if v, ok := S.metadata.unicons[variant]; ok {
		return v, nil
	}

	jsonContents, err := json.ReadJson[unicons](
		path.Join(S.manifest.BasePath, S.manifest.VendorsPaths["unicons"].Metadata, strings.Join([]string{variant, ".json"}, "")),
	)

	if err != nil {
		return nil, err
	}

	result := make(map[string][]string, len(jsonContents))
	for _, entry := range jsonContents {
		result[entry.Name] = entry.Tags
	}

	S.metadata.unicons[variant] = result

	return result, nil
}

func singleFile[T Fontawesome | Octicons | Remixicon](cacheEntry *T, manifest json.IcoManifest, vendor string) (T, error) {
	if *cacheEntry != nil {
		return *cacheEntry, nil
	}

	result, err := json.ReadJson[T](path.Join(manifest.BasePath, manifest.VendorsPaths[vendor].Metadata))

	if err != nil {
		var empty T
		return empty, err
	}

	*cacheEntry = result

	return result, nil
}

func NewStore(manifest json.IcoManifest) Store {
	return Store{manifest: manifest, metadata: &metadata{}}
}

type Store struct {
	metadata *metadata
	manifest json.IcoManifest
}

type metadata struct {
	fluentui    map[string]Fluentui
	fontawesome Fontawesome
	octicons    Octicons
	remixicon   Remixicon
	unicons     map[string]UniconsQuickAccess
}

type Fluentui struct {
	Metaphor []string
}

type Fontawesome map[string]struct{ Search struct{ Terms []string } }

type Octicons map[string][]string

type Remixicon map[string]map[string]string

type unicons []struct {
	Name string
	Tags []string
}

type UniconsQuickAccess map[string][]string
