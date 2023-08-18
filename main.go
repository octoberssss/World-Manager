package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type WorldManager struct {
	s      *server.Server
	log    *logrus.Logger
	Path   string
	worlds map[string]*world.World
	mutx   sync.Mutex
}

func CreateWorldManager(s *server.Server, path string) *WorldManager {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return &WorldManager{
		s:      s,
		Path:   path,
		worlds: make(map[string]*world.World),
	}
}

func (m *WorldManager) LoadWorldFromPath(worldname string) error {
	m.mutx.Lock()
	defer m.mutx.Unlock()

	if m.worlds[worldname] != nil {
		return nil
	}

	folder := m.Path + "/" + worldname
	worldData, err := mcdb.Config.Open(mcdb.Config{
		BlockSize:   16 * opt.KiB,
		Compression: opt.DefaultCompression,
		ReadOnly:    true,
		Log:         m.log,
	}, folder)

	if err != nil {
		return err
	}

	w := world.Config{
		Dim:      world.Overworld,
		Log:      m.log,
		ReadOnly: true,
		Provider: worldData,
		Entities: entity.DefaultRegistry,
	}.New()
	w.SetTickRange(0)
	w.SetTime(6000)
	w.StopTime()

	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeSurvival)

	m.worlds[worldname] = w
	fmt.Println("Loaded world", worldname)
	return nil
}

func (m *WorldManager) GetWorld(name string) *world.World {
	m.mutx.Lock()
	defer m.mutx.Unlock()

	return m.worlds[name]
}
