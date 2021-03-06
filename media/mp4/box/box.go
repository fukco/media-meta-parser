package box

import (
	"github.com/fukco/media-meta-parser/media"
	"io"
	"os"
	"reflect"
)

type IBox interface {
	GetMeta(r io.ReadSeeker, bi *Info, ctx *media.Context, meta *media.Meta) error
}

var Map = make(map[media.BoxType]interface{}, 64)

func AppendBoxMap(boxType media.BoxType, i interface{}) {
	Map[boxType] = i
}

// Box is ISO/IEC 14496-12 Box
type Box struct {
	Size     uint64
	Type     [4]byte
	UserType [16]byte
}

// FullBox is ISO/IEC 14496-12 FullBox
type FullBox struct {
	Box
	Version uint8
	Flags   [3]byte
}

type ContainerBox struct {
}

func (b *ContainerBox) getMeta(io.ReadSeeker, *Info, *media.Context) error {
	return nil
}

func isSupportedBox(bi *Info) bool {
	if _, ok := Map[bi.Type]; ok {
		return true
	}
	return false
}

func IsContainerBox(bi *Info) bool {
	if v, ok := Map[bi.Type]; ok {
		t := reflect.TypeOf(v)
		for i := 0; i < t.NumField(); i++ {
			if reflect.TypeOf(ContainerBox{}) == t.Field(i).Type {
				return true
			}
		}
	}
	return false
}

func GetMetaBoxes(file *os.File) ([]*Info, error) {
	boxInfos := make([]*Info, 0, 8)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	for {
		// read 8 bytes
		if bi, err := ReadBoxInfo(file); err == nil {
			if isSupportedBox(bi) {
				if IsContainerBox(bi) {
					_, err = bi.SeekToPayload(file)
					if err != nil {
						return nil, err
					}
					if isFullBox(bi) { //fullBox has 1 byte version and 3 bytes flags as box body
						_, err := file.Seek(4, io.SeekCurrent)
						if err != nil {
							return nil, err
						}
					}
					continue
				} else {
					boxInfos = append(boxInfos, bi)
				}
			}
			_, err := bi.SeekToEnd(file)
			if err != nil {
				return nil, err
			}
		} else {
			if err == io.EOF {
				break
			}
			return nil, err
		}

	}
	return boxInfos, nil
}
