package cdr

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/dhowden/tag"
)

var trackNameReg = regexp.MustCompile("^([0-9]{2}).+")

func NewTrack(file string) (*Track, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error loading file: %v", err)
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}

	return &Track{
		Title:    m.Title(),
		Artist:   m.Artist(),
		Filename: f.Name(),
	}, nil
}

// RenameTrack takes a filename, opens it, reads the metadata, and returns both
// the old and new filename.
func RenameTrack(file string) string {
	t, err := NewTrack(file)
	if err != nil {
		return ""
	}

	// Extract playlist track number from filename
	fMatch := trackNameReg.FindStringSubmatch(t.Filename)
	if len(fMatch) < 2 {
		log.Fatal("Unexpect filename format")
	}
	trackNum := fMatch[1]

	ext := t.Filename[strings.LastIndex(t.Filename, "."):]

	return fmt.Sprintf("%s-%s-%s%s", trackNum, Sanitize(t.Artist), Sanitize(t.Title), ext)
}

// Sanitize takes a string and removes problematic characters from it.
func Sanitize(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '(' || r == ')' || r == '[' || r == ']' || r == '.' {
			return -1
		}
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, s)
	return s
}
