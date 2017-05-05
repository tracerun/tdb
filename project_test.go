package tdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectMethods(t *testing.T) {
	projectTestFolder := "project_test_folder"
	db, err := Open(projectTestFolder)
	assert.NoError(t, err, "should have no error opening TDB")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(projectTestFolder)

	projectPath := "home"
	targets := []string{filepath.Join(projectPath, "a"), filepath.Join(projectPath, "b")}
	one := filepath.Join(projectPath, "c")

	// can't add target to non-exist project
	err = db.addTargetToProject(projectPath, one)
	assert.Equal(t, ErrProjectUnavailable, err, "project should be unavailable")

	// can't get targets from non-exist project
	_, err = db.GetProjectTargets(projectPath)
	assert.Equal(t, ErrProjectUnavailable, err, "project should be unavailable")

	// add targets[0]
	err = db.addSlot(targets[0], 0, 0)
	assert.NoError(t, err, "should have no error when adding slot")

	// add targets[1]
	err = db.addSlot(targets[1], 0, 0)
	assert.NoError(t, err, "should have no error when adding slot")

	err = db.CreateProject(projectPath)
	assert.NoError(t, err, "should have no error when creating project")

	err = db.addTargetToProject(projectPath, one)
	assert.NoError(t, err, "should have no error when adding a target to project")

	targetResult, err := db.GetProjectTargets(projectPath)
	assert.NoError(t, err, "should have no error when getting targets from a project")
	assert.EqualValues(t, append(targets, one), targetResult, "target result wrong")
}
