package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/textmap"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Generator struct {
	GeneratorConfig
}

type GeneratorConfig struct {
	Characters []*model.AvatarData
	Weapons    []*model.WeaponData
	Artifacts  []*model.ArtifactData
	Enemies    []*model.MonsterData
	Languages  map[string]string // map of languages and their corresponding textmap
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	return &Generator{
		GeneratorConfig: cfg,
	}, nil
}

func (g *Generator) DumpUIJSON(path string) error {
	// delete existing
	err := g.writeTranslationJSON(path + "/names.generated.json")
	if err != nil {
		return err
	}
	return nil
}

type OutData struct {
	CharacterNames map[string]string `json:"character_names"`
	WeaponNames    map[string]string `json:"weapon_names"`
	ArtifactNames  map[string]string `json:"artifact_names"`
	EnemyNames     map[string]string `json:"enemy_names"`
}

func (g *Generator) GetNames(lang string) (OutData, error) {
	data := OutData{
		CharacterNames: make(map[string]string),
		WeaponNames:    make(map[string]string),
		ArtifactNames:  make(map[string]string),
		EnemyNames:     make(map[string]string),
	}
	// load generator for this language
	tp := g.Languages[lang]
	src, err := textmap.NewTextMapSource(tp)
	if err != nil {
		return data, fmt.Errorf("error creating text map src for %v: %w", lang, err)
	}
	// go through all char/weap/art and get names
	for _, v := range g.Characters {
		s, err := src.Get(v.NameTextHashMap)
		if err != nil {
			fmt.Printf("error getting string for char %v id %v\n", v.Key, v.NameTextHashMap)
			continue
		}
		data.CharacterNames[v.Key] = s
	}
	for _, v := range g.Weapons {
		s, err := src.Get(v.NameTextHashMap)
		if err != nil {
			fmt.Printf("error getting string for weapon %v id %v\n", v.Key, v.NameTextHashMap)
			continue
		}
		data.WeaponNames[v.Key] = s
	}
	for _, v := range g.Artifacts {
		s, err := src.Get(v.TextMapId)
		if err != nil {
			fmt.Printf("error getting string for set %v id %v\n", v.Key, v.TextMapId)
			continue
		}
		data.ArtifactNames[v.Key] = s
	}
	for _, v := range g.Enemies {
		s, err := src.Get(v.NameTextHashMap)
		if err != nil {
			fmt.Printf("error getting string for enemy %v id %v\n", v.Key, v.NameTextHashMap)
			continue
		}
		data.EnemyNames[v.Key] = s
	}
	return data, nil
}

func (g *Generator) writeTranslationJSON(path string) error {
	// sort keys first
	var keys []string
	for k := range g.Languages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make(map[string]OutData)

	for _, k := range keys {
		data, err := g.GetNames(k)
		if err != nil {
			return err
		}
		out[k] = data
	}

	data, err := json.MarshalIndent(out, "", "   ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
