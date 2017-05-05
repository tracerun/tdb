package tdb

import (
	"path/filepath"

	"os"
	"strings"

	"github.com/golang/protobuf/proto"
)

const (
	projectIndex = "__project__"
)

// load project index from file
func (db *TDB) loadProjectIndex() error {
	var err error
	db.project, err = createInfo(filepath.Join(db.path, projectIndex))
	return err
}

// CreateProject to create a project in db
// A project is a target container containing group of targets.
func (db *TDB) CreateProject(projPath string) error {
	cleanedProject := filepath.Clean(projPath)

	db.project.contentLock.Lock()
	_, ok := db.project.content[cleanedProject]
	if ok {
		db.project.contentLock.Unlock()
		return nil
	}

	var dbProj Project
	dbProj.Targets = db.getBelongedTargets(cleanedProject)
	bs, err := proto.Marshal(&dbProj)
	if err != nil {
		db.project.contentLock.Unlock()
		return err
	}

	db.project.content[cleanedProject] = bs
	db.project.contentLock.Unlock()

	db.project.writeToDisk()
	return nil
}

// GetProjectTargets to get targets of a project
func (db *TDB) GetProjectTargets(projPath string) ([]string, error) {
	db.project.contentLock.RLock()
	bs, ok := db.project.content[projPath]
	if !ok {
		db.project.contentLock.RUnlock()
		return nil, ErrProjectUnavailable
	}

	var dbProj Project
	err := proto.Unmarshal(bs, &dbProj)
	if err != nil {
		db.project.contentLock.RUnlock()
		return nil, err
	}

	db.project.contentLock.RUnlock()
	return dbProj.Targets, nil
}

// addTargetToProject to add a target to project.
func (db *TDB) addTargetToProject(projPath, target string) error {
	db.project.contentLock.Lock()
	bs, ok := db.project.content[projPath]
	if !ok {
		db.project.contentLock.Unlock()
		return ErrProjectUnavailable
	}

	if !targetBelong(target, projPath) {
		db.project.contentLock.Unlock()
		return ErrTargetNotBelongToProject
	}

	var dbProj Project
	if err := proto.Unmarshal(bs, &dbProj); err != nil {
		db.project.contentLock.Unlock()
		return err
	}

	existed := false
	for i := 0; i < len(dbProj.Targets); i++ {
		if dbProj.Targets[i] == target {
			existed = true
			break
		}
	}
	if existed {
		db.project.contentLock.Unlock()
		return nil
	}

	dbProj.Targets = append(dbProj.Targets, target)
	bs, err := proto.Marshal(&dbProj)
	if err != nil {
		db.project.contentLock.Unlock()
		return err
	}
	db.project.content[projPath] = bs
	db.project.contentLock.Unlock()
	db.project.writeToDisk()
	return nil
}

func (db *TDB) getBelongedTargets(proj string) []string {
	var targets []string
	allTargets := db.slot.getInfoKeys()

	for i := 0; i < len(allTargets); i++ {
		if targetBelong(allTargets[i], proj) {
			targets = append(targets, allTargets[i])
		}
	}
	return targets
}

// target, project should be cleaned first
func targetBelong(target, project string) bool {
	if !strings.HasSuffix(project, string(os.PathSeparator)) {
		project = strings.Join([]string{project, string(os.PathSeparator)}, "")
	}
	return strings.HasPrefix(target, project)
}
