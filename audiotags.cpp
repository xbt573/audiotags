/***************************************************************************
    copyright            : (C) 2003 by Scott Wheeler
    email                : wheeler@kde.org
 ***************************************************************************/

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

#include <stdlib.h>
#include <fileref.h>
#include <flacfile.h>
#include <flacpicture.h>
#include <mp4file.h>
#include <tbytevector.h>
#include <tbytevectorstream.h>
#include <tfile.h>
#include <tpropertymap.h>
#include <string.h>
#include <typeinfo>

#include "audiotags.h"

static bool unicodeStrings = true;


class ByteVectorStreamWithName : public TagLib::ByteVectorStream {
    public:
        ByteVectorStreamWithName(const char* name, const TagLib::ByteVector &data) : TagLib::ByteVectorStream(data) {
            this->fileName = TagLib::FileName(name);
        }
        TagLib::FileName name() const {
            return this->fileName;
        }

    private:
        TagLib::FileName fileName;
};


TagLib_FileRefRef *audiotags_file_new(const char *filename)
{
  TagLib::FileRef *fr = new TagLib::FileRef(filename);
  if (fr == NULL || fr->isNull() || !fr->file()->isValid() || fr->tag() == NULL) {
    if (fr) {
      delete fr;
      fr = NULL;
    }
    return NULL;
  }

  TagLib_FileRefRef *holder = new TagLib_FileRefRef();
  holder->fileRef = reinterpret_cast<void *>(fr);
  holder->ioStream = NULL;
  return holder;
}

TagLib_FileRefRef *audiotags_file_memory(const char *data, unsigned int length) {
  TagLib::ByteVectorStream *ioStream = new TagLib::ByteVectorStream(TagLib::ByteVector(data, length));
  TagLib::FileRef *fr = new TagLib::FileRef(ioStream);
  if (fr == NULL || fr->isNull() || !fr->file()->isValid() || fr->tag() == NULL) {
    if (fr) {
      delete fr;
      fr = NULL;
    }
    if (ioStream) {
      delete ioStream;
      ioStream = NULL;
    }
    return NULL;
  }

  TagLib_FileRefRef *holder = new TagLib_FileRefRef();
  holder->fileRef = reinterpret_cast<void *>(fr);
  holder->ioStream = reinterpret_cast<void *>(ioStream);
  return holder;
}

TagLib_FileRefRef *audiotags_file_memory_with_name(const char *fileName, const char *data, unsigned int length) {
  TagLib::ByteVectorStream *ioStream = new ByteVectorStreamWithName(fileName, TagLib::ByteVector(data, length));
  TagLib::FileRef *fr = new TagLib::FileRef(ioStream);
  if (fr == NULL || fr->isNull() || !fr->file()->isValid() || fr->tag() == NULL) {
    if (fr) {
      delete fr;
      fr = NULL;
    }
    if (ioStream) {
      delete ioStream;
      ioStream = NULL;
    }
    return NULL;
  }

  TagLib_FileRefRef *holder = new TagLib_FileRefRef();
  holder->fileRef = reinterpret_cast<void *>(fr);
  holder->ioStream = reinterpret_cast<void *>(ioStream);
  return holder;
}

void audiotags_file_close(TagLib_FileRefRef *fileRefRef)
{
  delete reinterpret_cast<TagLib::FileRef *>(fileRefRef->fileRef);
  if (fileRefRef->ioStream) {
    delete reinterpret_cast<TagLib::IOStream *>(fileRefRef->ioStream);
  }
  delete fileRefRef;
}

void audiotags_file_properties(const TagLib_FileRefRef *fileRefRef, int id)
{
  const TagLib::FileRef *fileRef = reinterpret_cast<const TagLib::FileRef *>(fileRefRef->fileRef);
  TagLib::PropertyMap tags = fileRef->file()->properties();
  for(TagLib::PropertyMap::ConstIterator i = tags.begin(); i != tags.end(); ++i) {
    for(TagLib::StringList::ConstIterator j = i->second.begin(); j != i->second.end(); ++j) {
      char *key = ::strdup(i->first.toCString(unicodeStrings));
      char *val = ::strdup((*j).toCString(unicodeStrings));
      go_map_put(id, key, val);
      free(key);
      free(val);
    }
  }
}

bool audiotags_write_property(TagLib_FileRefRef *fileRefRef, const char *field_c, const char *value_c)
{
  return audiotags_write_properties(fileRefRef, 1, &field_c, &value_c);
}

bool audiotags_write_properties(TagLib_FileRefRef *fileRefRef, unsigned int len, const char *fields_c[], const char *values_c[])
{
  TagLib::FileRef *fileRef = reinterpret_cast<TagLib::FileRef *>(fileRefRef->fileRef);
  TagLib::Tag *t = fileRef->tag();

  bool prop_changed = false;
  for(TagLib::uint i = 0; i < len; i++) {
    TagLib::String field(fields_c[i], TagLib::String::Type::UTF8);
    TagLib::String value(values_c[i], TagLib::String::Type::UTF8);
    if(field == "title") {
      t->setTitle(value);
    } else if(field == "artist") {
      t->setArtist(value);
    } else if(field == "album") {
      t->setAlbum(value);
    } else if(field == "comment") {
      t->setComment(value);
    } else if(field == "genre") {
      t->setGenre(value);
    } else if(field == "year") {
      t->setYear(value.toInt());
    } else if(field == "track") {
      t->setTrack(value.toInt());
    } else {
      TagLib::PropertyMap tags = fileRef->file()->properties();
      if(!tags.contains(field)) {
        tags.insert(field, value);
      } else {
        tags.replace(field, value);
      }
      if((fileRef->file()->setProperties(tags)).size() > 0) {
        return false;
      }
    }
    fileRef->save();
  }

  return true;
}

const TagLib_AudioProperties *audiotags_file_audioproperties(const TagLib_FileRefRef *fileRefRef)
{
  const TagLib::FileRef *fileRef = reinterpret_cast<const TagLib::FileRef *>(fileRefRef->fileRef);
  return reinterpret_cast<const TagLib_AudioProperties *>(fileRef->file()->audioProperties());
}

const TagLib::AudioProperties *props(const TagLib_AudioProperties *audioProperties)
{
  return reinterpret_cast<const TagLib::AudioProperties *>(audioProperties);
}

int audiotags_audioproperties_length(const TagLib_AudioProperties *audioProperties)
{
  return props(audioProperties)->length();
}

int audiotags_audioproperties_length_ms(const TagLib_AudioProperties *audioProperties)
{
  return props(audioProperties)->lengthInMilliseconds();
}

int audiotags_audioproperties_bitrate(const TagLib_AudioProperties *audioProperties)
{
  return props(audioProperties)->bitrate();
}

int audiotags_audioproperties_samplerate(const TagLib_AudioProperties *audioProperties)
{
  return props(audioProperties)->sampleRate();
}

int audiotags_audioproperties_channels(const TagLib_AudioProperties *audioProperties)
{
  return props(audioProperties)->channels();
}

enum img_type {
  JPEG = 0,
  PNG = 1,
  // to be continued...
};

bool audiotags_write_picture(TagLib_FileRefRef *fileRefRef, const char *data, unsigned int length, int w, int h, int type)
{
  const TagLib::FileRef *fileRef = reinterpret_cast<const TagLib::FileRef *>(fileRefRef->fileRef);
  const TagLib::ByteVector byte_vec = TagLib::ByteVector(data, length);

  // check which type the file is (flac, mp4, etc)
  if(TagLib::FLAC::File *flac = dynamic_cast<TagLib::FLAC::File *>(fileRef->file())) {
    TagLib::FLAC::Picture *pic = new TagLib::FLAC::Picture;
    // only front cover type supported for now
    pic->setType(TagLib::FLAC::Picture::Type::FrontCover);
    if(type == PNG) {
      pic->setMimeType("image/png");
    } else if (type == JPEG) {
      pic->setMimeType("image/jpeg");
    } else {
      return false;
    }

    pic->setData(byte_vec);
    pic->setWidth(w);
    pic->setHeight(h);
    pic->setColorDepth(24);
    pic->setNumColors(16777216);
    flac->addPicture(pic);
    flac->save();
  } else if (TagLib::MP4::File *mp4 = dynamic_cast<TagLib::MP4::File *>(fileRef->file())) {
    TagLib::MP4::CoverArt::Format fmt = TagLib::MP4::CoverArt::Format::Unknown;
    if(type == PNG) {
      fmt = TagLib::MP4::CoverArt::Format::PNG;
    } else if (type == JPEG) {
      fmt = TagLib::MP4::CoverArt::Format::JPEG;
    } else {
      return false;
    }

    TagLib::MP4::CoverArtList l = mp4->tag()->item("covr").toCoverArtList();
    l.prepend(TagLib::MP4::CoverArt(fmt, byte_vec));
    mp4->tag()->setItem("covr", l);
    mp4->save();
  } else {
    return false;
  }
  return true;
}

bool audiotags_remove_pictures(TagLib_FileRefRef *fileRefRef)
{
  const TagLib::FileRef *fileRef = reinterpret_cast<const TagLib::FileRef *>(fileRefRef->fileRef);

  // check which type the file is (flac, mp4, etc)
  if(TagLib::FLAC::File *flac = dynamic_cast<TagLib::FLAC::File *>(fileRef->file())) {
    flac->removePictures();
    flac->save();
  } else if (TagLib::MP4::File *mp4 = dynamic_cast<TagLib::MP4::File *>(fileRef->file())) {
    mp4->tag()->removeItem("covr");
    mp4->save();
  } else {
    return false;
  }
  return true;
}
