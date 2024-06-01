package changeset

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-way/changesets/pkg/version"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

var names = []string{"hello", "world", "dog", "arnold", "cat", "kitten", "puppy", "armadillo", "giraffe", "happy", "sad", "emotional", "earth", "mars", "car", "robot", "whale", "python"}

const CHANGE_NAME_PARTS int8 = 3
const CHANGESET_DIRECTORY string = ".changeset"
const CHANGESET_FILE_KEY string = "changeset/type"

type Change struct {
	BumpType version.BumpType
	Message  string
	FilePath string
}

type Changeset struct {
	// The current version
	CurrentVersion version.Version
	Changes        []Change
}

func getRandomName() string {
	return names[rand.Intn(len(names))]
}

func generateChangeName() string {
	var parts [CHANGE_NAME_PARTS]string
	for i := range CHANGE_NAME_PARTS {
		random_name := getRandomName()
		parts[i] = random_name
	}
	return strings.Join(parts[:], "-")
}

// Creates a new change set file and returns the filepath
func CreateChangeFile(bump_type version.BumpType, message string) (string, error) {
	filename := generateChangeName()

	if _, err := os.Stat(CHANGESET_DIRECTORY); os.IsNotExist(err) {
		err := os.Mkdir(CHANGESET_DIRECTORY, 0755)
		if err != nil {
			return "", err
		}
	}

	filepath := filepath.Join(CHANGESET_DIRECTORY, filename+".md")

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file_contents := "---\n" + CHANGESET_FILE_KEY + ": " + bump_type.String() + "\n" + "---\n\n# " + message + "\n"

	_, err = file.WriteString(file_contents)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

func (cs *Changeset) DetermineFinalBumpType() version.BumpType {
	var highest_version_type version.BumpType = version.Undetermined
	for _, change := range cs.Changes {
		if change.BumpType >= highest_version_type {
			highest_version_type = change.BumpType
		}
	}
	return highest_version_type
}

func GetChanges() ([]Change, error) {
	files, err := filepath.Glob(CHANGESET_DIRECTORY + "/*.md")

	if err != nil {
		panic(err)
	}

	var changes []Change

	for _, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		contents, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		markdown := goldmark.New(
			goldmark.WithExtensions(
				meta.Meta,
			),
		)

		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert([]byte(contents), &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}
		metaData := meta.Get(context)
		ver := metaData[CHANGESET_FILE_KEY].(string)
		if ver == "" {
			return nil, errors.New("changeset file does not have a type")
		}
		parsed_bump_type, err := version.ParseBumpType(ver)
		if err != nil {
			return nil, err
		}
		changes = append(changes, Change{BumpType: parsed_bump_type, Message: "thing", FilePath: filepath})
	}
	return changes, nil
}

func (cs *Changeset) DetermineNextVersion() version.Version {
	next_version := cs.CurrentVersion
	next_version.Bump(cs.DetermineFinalBumpType())
	return next_version
}

// Consumes the associated changes and returns the new version
func (cs *Changeset) ConsumeChanges() (version.Version, error) {
	if len(cs.Changes) == 0 {
		return version.Version{}, errors.New("no changesets found")
	}

	new_version := cs.DetermineNextVersion()
	for _, change := range cs.Changes {
		os.Remove(change.FilePath)
	}

	return new_version, nil
}
