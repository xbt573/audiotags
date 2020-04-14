/***************************************************************************
   copyright            : (C) 2014 by Nick Sellen
   email                : code@nicksellen.co.uk
***************************************************************************/

/***************************************************************************
 *   This library is free software; you can redistribute it and/or modify  *
 *   it  under the terms of the GNU Lesser General Public License version  *
 *   2.1 as published by the Free Software Foundation.                     *
 *                                                                         *
 *   This library is distributed in the hope that it will be useful, but   *
 *   WITHOUT ANY WARRANTY; without even the implied warranty of            *
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU     *
 *   Lesser General Public License for more details.                       *
 *                                                                         *
 *   You should have received a copy of the GNU Lesser General Public      *
 *   License along with this library; if not, write to the Free Software   *
 *   Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  *
 *   USA                                                                   *
 ***************************************************************************/

package audiotags

/*
#cgo pkg-config: taglib
#cgo LDFLAGS: -lstdc++
#include "audiotags.h"
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"unsafe"
)

import "fmt"

const (
	JPEG = iota
	PNG = iota
)

type File C.TagLib_FileRefRef

type AudioProperties struct {
	Length, LengthMs, Bitrate, Samplerate, Channels int
}

func Open(filename string) (*File, error) {
	fp := C.CString(filename)
	defer C.free(unsafe.Pointer(fp))
	f := C.audiotags_file_new(fp)
	if f == nil {
		return nil, fmt.Errorf("cannot process %s", filename)
	}
	return (*File)(f), nil
}

func FromData(data []byte) (*File, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("got empty byte array")
	}

	// actually parse data
	f := C.audiotags_file_memory((*C.char)(unsafe.Pointer(&data[0])), C.uint(len(data)))
	if f == nil {
		return nil, fmt.Errorf("cannot process provided data")
	}
	return (*File)(f), nil
}

func FromDataWithName(filename string, data []byte) (*File, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("got empty byte array")
	}

	// actually parse data
	fp := C.CString(filename)
	defer C.free(unsafe.Pointer(fp))
	f := C.audiotags_file_memory_with_name(fp, (*C.char)(unsafe.Pointer(&data[0])), C.uint(len(data)))
	if f == nil {
		return nil, fmt.Errorf("cannot process provided data")
	}
	return (*File)(f), nil
}

func Read(filename string) (map[string]string, *AudioProperties, error) {
	f, err := Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return f.ReadTags(), f.ReadAudioProperties(), nil
}

func ReadFromData(data []byte) (map[string]string, *AudioProperties, error) {
	f, err := FromData(data)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return f.ReadTags(), f.ReadAudioProperties(), nil
}

func ReadFromDataWithName(filename string, data []byte) (map[string]string, *AudioProperties, error) {
	f, err := FromDataWithName(filename, data)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return f.ReadTags(), f.ReadAudioProperties(), nil
}

func ReadTags(filename string) (map[string]string, error) {
	f, err := Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ReadTags(), nil
}

func ReadAudioProperties(filename string) (*AudioProperties, error) {
	f, err := Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ReadAudioProperties(), nil
}

func (f *File) Close() {
	C.audiotags_file_close((*C.TagLib_FileRefRef)(f))
}

func (f *File) ReadTags() map[string]string {
	id := mapsNextId
	mapsNextId++
	m := make(map[string]string)
	maps[id] = m
	C.audiotags_file_properties((*C.TagLib_FileRefRef)(f), C.int(id))
	delete(maps, id)
	return m
}

func (f *File) WriteTag(tag, value string) bool {
	tagC := C.CString(tag)
	defer C.free(unsafe.Pointer(tagC))
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))
	if C.audiotags_write_property((*C.TagLib_FileRefRef)(f), tagC, valueC) {
		return true
	} else {
		return false
	}
}

func (f *File) WriteTags(tag_map map[string]string) bool {
	tagFields := make([]*C.char, len(tag_map))
	tagValues := make([]*C.char, len(tag_map))
	i := 0
	for field, value := range tag_map {
		fieldC := C.CString(field)
		tagFields[i] = fieldC
		defer C.free(unsafe.Pointer(fieldC))

		valueC := C.CString(value)
		tagValues[i] = valueC
		defer C.free(unsafe.Pointer(valueC))
		i++
	}
	if C.audiotags_write_properties((*C.TagLib_FileRefRef)(f), C.uint(len(tag_map)), &tagFields[0], &tagValues[0]) {
		return true
	} else {
		return false
	}
}

func (f *File) ReadAudioProperties() *AudioProperties {
	ap := C.audiotags_file_audioproperties((*C.TagLib_FileRefRef)(f))
	if ap == nil {
		return nil
	}
	p := AudioProperties{}
	p.Length = int(C.audiotags_audioproperties_length(ap))
	p.LengthMs = int(C.audiotags_audioproperties_length_ms(ap))
	p.Bitrate = int(C.audiotags_audioproperties_bitrate(ap))
	p.Samplerate = int(C.audiotags_audioproperties_samplerate(ap))
	p.Channels = int(C.audiotags_audioproperties_channels(ap))
	return &p
}

func (f *File) WritePicture(data []byte, fmt, w, h int) bool {
	if len(data) == 0 {
		return false
	}

	if C.audiotags_write_picture((*C.TagLib_FileRefRef)(f), (*C.char)(unsafe.Pointer(&data[0])), C.uint(len(data)),
			C.int(w), C.int(h), C.int(fmt)) {
		return true
	} else {
		return false
	}
}

func (f *File) RemovePictures() bool {
	if C.audiotags_remove_pictures((*C.TagLib_FileRefRef)(f)) {
		return true
	} else {
		return false
	}
}

var maps = make(map[int]map[string]string)
var mapsNextId = 0

//export go_map_put
func go_map_put(id C.int, key *C.char, val *C.char) {
	m := maps[int(id)]
	k := strings.ToLower(C.GoString(key))
	v := C.GoString(val)
	m[k] = v
}
