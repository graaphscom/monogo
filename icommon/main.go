package main

import (
	"context"
	"fmt"
	"github.com/graaphscom/monogo/icommon/json"
	"github.com/graaphscom/monogo/icommon/unitree"
	"github.com/redis/rueidis"
	"log"
	"strings"
)

func main() {
	manifest, err := json.ReadJson[json.IcoManifest]("testdata/ico_manifest_downloads.json")
	tree, err := unitree.BuildRootTree(manifest)

	if err != nil {
		log.Fatalln(err)
	}

	opt, err := rueidis.ParseURL("redis://localhost:6379")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := rueidis.NewClient(opt)
	if err != nil {
		log.Fatalln(err)
	}

	iconsCount := 0
	tree.Traverse([]string{}, func(_ []string, iconSet unitree.IconSet) {
		for range iconSet.Icons {
			iconsCount++
		}
	})

	commands := make([]rueidis.Completed, iconsCount+1)
	commandsIdx := 0
	tree.Traverse([]string{}, func(segments []string, iconSet unitree.IconSet) {
		for _, icon := range iconSet.Icons {
			commands[commandsIdx] = conn.B().Hset().
				Key(strings.Join(append(segments, icon.Name), ":")).
				FieldValue().
				FieldValue("searchTags", strings.Join(icon.Tags.Search, ",")).
				FieldValue("visualTags", strings.Join(icon.Tags.Visual, ",")).
				Build()
			commandsIdx++
		}
	})

	ctx := context.Background()
	commands[commandsIdx] = conn.B().FtCreate().
		Index("icommon").
		Prefix(1).Prefix("icommon:").
		Schema().
		FieldName("searchTags").Text().
		FieldName("visualTags").Tag().
		Build()

	conn.DoMulti(ctx, commands...)

	fmt.Printf("Total icons count: %d", iconsCount)
}
