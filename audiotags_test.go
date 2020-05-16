package audiotags

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
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

func TestWriteTagLib(t *testing.T) {
	fileName := "test.mp3"
	file, err := Open(fileName)

	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}
	tempDir, err := ioutil.TempDir("", "go-taglib-test")

	if err != nil {
		t.Fatalf("Cannot create temporary file for writing tests: %s", err)
	}

	tempFileName := path.Join(tempDir, "go-taglib-test.mp3")

	defer file.Close()
	defer os.RemoveAll(tempDir)

	err = cp(tempFileName, fileName)

	if err != nil {
		t.Fatalf("Cannot copy file for writing tests: %s", err)
	}

	modifiedFile, err := Open(tempFileName)
	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}

	tags := file.ReadTags()
	writeTags := make(map[string]string)
	writeTags["artist"] = getModifiedString(tags["artist"]) 
	writeTags["album"] = getModifiedString(tags["album"])
	writeTags["title"] = getModifiedString(tags["title"])
	writeTags["comment"] = getModifiedString(tags["comment"])
	writeTags["genre"] = getModifiedString(tags["genre"])
	writeTags["tracknumber"] = getModifiedInt(tags["tracknumber"], t)
	writeTags["date"] = getModifiedInt(tags["date"], t)

	modifiedFile.WriteTags(writeTags)
	if err != nil {
		t.Fatalf("Cannot save file : %s", err)
	}
	modifiedFile.Close()
	
	//Re-open the modified file
	modifiedFile, err = Open(tempFileName)
	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}

	// Test the Tags
	newTags := modifiedFile.ReadTags()
	if newTags["title"] != getModifiedString("The Title") {
		t.Errorf("Got wrong modified title: %s", newTags["title"])
	}

	if newTags["artist"] != getModifiedString("The Artist") {
		t.Errorf("Got wrong modified artist: %s", newTags["artist"])
	}

	if newTags["album"] != getModifiedString("The Album") {
		t.Errorf("Got wrong modified album: %s", newTags["album"])
	}

	if newTags["comment"] != getModifiedString("A Comment") {
		t.Errorf("Got wrong modified comment: %s", newTags["comment"])
	}

	if newTags["genre"] != getModifiedString("Booty Bass") {
		t.Errorf("Got wrong modified genre: %s", newTags["genre"])
	}

	if newTags["date"] != getModifiedInt("1942", t) {
		t.Errorf("Got wrong modified year: %v", newTags["date"])
	}

	if newTags["tracknumber"] != getModifiedInt("42", t) {
		t.Errorf("Got wrong modified track: %v", newTags["tracknumber"])
	}
}

func TestGenericWriteTagLib(t *testing.T) {
	fileName := "test.mp3"
	file, err := Open(fileName)

	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}
	tempDir, err := ioutil.TempDir("", "go-taglib-test")

	if err != nil {
		t.Fatalf("Cannot create temporary file for writing tests: %s", err)
	}

	tempFileName := path.Join(tempDir, "go-taglib-test.mp3")

	defer file.Close()
	defer os.RemoveAll(tempDir)

	err = cp(tempFileName, fileName)

	if err != nil {
		t.Fatalf("Cannot copy file for writing tests: %s", err)
	}

	modifiedFile, err := Open(tempFileName)
	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}
	tags := file.ReadTags()
	writeTags := make(map[string]string)
	writeTags["artist"] = getModifiedString(tags["artist"]) 
	writeTags["album"] = getModifiedString(tags["album"])
	writeTags["title"] = getModifiedString(tags["title"])
	writeTags["comment"] = getModifiedString(tags["comment"])
	writeTags["genre"] = getModifiedString(tags["genre"])
	writeTags["tracknumber"] = getModifiedInt(tags["tracknumber"], t)
	writeTags["date"] = getModifiedInt(tags["date"], t)

	modifiedFile.WriteTags(writeTags)
	if err != nil {
		t.Fatalf("Cannot save file : %s", err)
	}
	modifiedFile.Close()

	//Re-open the modified file
	modifiedFile, err = Open(tempFileName)
	if err != nil {
		t.Fatalf("Read returned error: %s", err)
	}

	// Test the Tags
	newTags := modifiedFile.ReadTags()
	if newTags["title"] != getModifiedString("The Title") {
		t.Errorf("Got wrong modified title: %s", newTags["title"])
	}

	if newTags["artist"] != getModifiedString("The Artist") {
		t.Errorf("Got wrong modified artist: %s", newTags["artist"])
	}

	if newTags["album"] != getModifiedString("The Album") {
		t.Errorf("Got wrong modified album: %s", newTags["album"])
	}

	if newTags["comment"] != getModifiedString("A Comment") {
		t.Errorf("Got wrong modified comment: %s", newTags["comment"])
	}

	if newTags["genre"] != getModifiedString("Booty Bass") {
		t.Errorf("Got wrong modified genre: %s", newTags["genre"])
	}

	if newTags["date"] != getModifiedInt("1942", t) {
		t.Errorf("Got wrong modified year: %v", newTags["date"])
	}

	if newTags["tracknumber"] != getModifiedInt("42", t) {
		t.Errorf("Got wrong modified track: %v", newTags["tracknumber"])
	}
}

func checkModified(original string, modified string) bool {
	return modified == getModifiedString(original)
}

func getModifiedString(s string) string {
	return s + " MODIFIED"
}

func getModifiedInt(s string, t *testing.T) string {
	val, err := strconv.Atoi(s)
	if err != nil {
		t.Fatal(err)
	}

	return strconv.Itoa(val + 1)
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}