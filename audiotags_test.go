package audiotags

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadNothing(t *testing.T) {
	file, err := Open("doesnotexist.mp3")

	if file != nil {
		t.Fatal("Returned non nil file struct.")
	}

	if err == nil {
		t.Fatal("Returned nil err.")
	}
}

func TestReadDirectory(t *testing.T) {
	file, err := Open("/")

	if file != nil {
		t.Fatal("Returned non nil file struct.")
	}

	if err == nil {
		t.Fatal("Returned nil err.")
	}
}

func TestTagLib(t *testing.T) {
	file, err := Open("test.mp3")

	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}

	defer file.Close()

	// Test the Tags
	tags := file.ReadTags()
	if tags["title"] != "The Title" {
		t.Errorf("Got wrong title: %s", tags["title"])
	}

	if tags["artist"] != "The Artist" {
		t.Errorf("Got wrong artist: %s", tags["artist"])
	}

	if tags["album"] != "The Album" {
		t.Errorf("Got wrong album: %s", tags["album"])
	}

	if tags["comment"] != "A Comment" {
		t.Errorf("Got wrong comment: %s", tags["comment"])
	}

	if tags["genre"] != "Booty Bass" {
		t.Errorf("Got wrong genre: %s", tags["comment"])
	}

	if tags["date"] != "1942" {
		t.Errorf("Got wrong year: %s", tags["date"])
	}

	if tags["tracknumber"] != "42" {
		t.Errorf("Got wrong track: %s", tags["tracknumber"])
	}

	// Test the properties
	properties := file.ReadAudioProperties()
	if properties.Length != 42 {
		t.Errorf("Got wrong length: %v", properties.Length)
	}

	if properties.Bitrate != 128 {
		t.Errorf("Got wrong bitrate: %d", properties.Bitrate)
	}

	if properties.Samplerate != 44100 {
		t.Errorf("Got wrong samplerate: %d", properties.Samplerate)
	}

	if properties.Channels != 2 {
		t.Errorf("Got wrong channels: %d", properties.Channels)
	}
}

var expectedTags = map[string]string{
	"album":       "Test Album",
	"albumartist": "Test AlbumArtist",
	"artist":      "Test Artist",
	"composer":    "Test Composer",
	"date":        "2000",
	"description": "Test Comment",
	"discnumber":  "02",
	"genre":       "Jazz",
	"title":       "Test Title",
	"tracknumber": "03/06",
}

func TestReads(t *testing.T) {
	var tests = []struct {
		input        string
		expectedTags func(map[string]string) map[string]string
		noImage      bool
	}{
		{
			input:   "sample.ape",
			noImage: true,
		},
		{
			input: "sample.flac",
			expectedTags: func(m map[string]string) map[string]string {
				m["tracknumber"] = "03"
				m["tracktotal"] = "06"
				return m
			},
		},
		{
			input: "sample.id3v11.mp3",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "albumartist")
				delete(m, "composer")
				delete(m, "description")
				delete(m, "discnumber")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "3"
				return m
			},
		},
		{
			input: "sample.id3v22.mp3",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "03/06"
				return m
			},
		},
		{
			input: "sample.id3v23.mp3",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "03/06"
				return m
			},
		},
		{
			input: "sample.id3v24.mp3",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "03/06"
				return m
			},
		},
		{
			input: "sample.m4a",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "3/6"
				m["discnumber"] = "2"
				return m
			},
		},
		{
			input: "sample.mp4",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				delete(m, "tracktotal")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "3/6"
				m["discnumber"] = "2"
				return m
			},
		},
		{
			input: "sample.ogg",
			expectedTags: func(m map[string]string) map[string]string {
				delete(m, "description")
				m["comment"] = "Test Comment"
				m["tracknumber"] = "3"
				m["discnumber"] = "02"
				m["tracktotal"] = "06"
				return m
			},
		},
		{
			input:   "sample.wv",
			noImage: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			file, err := Open(filepath.Join("testdata", test.input))
			assert.NoError(t, err)
			assert.NotNil(t, file)
			assert.True(t, file.HasMedia())
			tags := file.ReadTags()
			expected := expectedTags
			if test.expectedTags != nil {
				tags := map[string]string{}
				for k, v := range expectedTags {
					tags[k] = v
				}
				expected = test.expectedTags(tags)
			}
			assert.Equal(t, expected, tags)

			if !test.noImage {
				img, err := file.ReadImage()
				assert.NoError(t, err)
				assert.NotNil(t, img)
			}
		})
	}
}
