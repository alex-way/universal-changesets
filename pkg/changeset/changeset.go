package changeset

import (
	"bytes"
	"fmt"
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

const CHANGESET_NAME_PARTS = 3
const CHANGESET_DIRECTORY = ".changeset"

type Change struct {
	Type    version.IncrementType
	Message string
}

type Changeset struct {
	Version version.Version
	Changes []Change
}

func getRandomName() string {
	return names[rand.Intn(len(names))]
}

func generateChangesetName() string {
	var parts [CHANGESET_NAME_PARTS]string
	for i := range CHANGESET_NAME_PARTS {
		random_name := getRandomName()
		parts[i] = random_name
	}
	return strings.Join(parts[:], "-")
}

// / Creates a new change set file
func CreateChangeset(incrementType version.IncrementType, message string) string {
	filename := generateChangesetName()

	if _, err := os.Stat(CHANGESET_DIRECTORY); os.IsNotExist(err) {
		err := os.Mkdir(CHANGESET_DIRECTORY, 0755)
		if err != nil {
			panic(err)
		}
	}

	filepath := filepath.Join(CHANGESET_DIRECTORY, filename+".md")

	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file_contents := "---\n" + "changeset/type: " + incrementType.String() + "\n" + "---\n\n# " + message + "\n"

	_, err = file.WriteString(file_contents)
	if err != nil {
		panic(err)
	}
	return filepath
}

func DetermineChangesetType() version.IncrementType {
	// TODO: Implement by reading all the files in the .changeset directory and determining the highest version type.
	// For example if there are multiple patches and one major, then return major
	return version.Patch
}

func ConsumeChangesets() []Changeset {
	// todo: consume all the changesets to create a new version
	files, err := os.ReadDir(CHANGESET_DIRECTORY)
	if err != nil {
		panic(err)
	}

	var changesets []Changeset

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := CHANGESET_DIRECTORY + "/" + file.Name()
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
		title := metaData["type"]
		fmt.Print(title)
	}
	return changesets
}
