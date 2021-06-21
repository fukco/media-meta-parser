package media

type Type string
type SupportMajorBrand string
type Extension string

const (
	MP4 Type = "MP4"
	MOV Type = "MOV(QUICKTIME)"
)

const (
	SONYXAVC SupportMajorBrand = "XAVC"
	QT       SupportMajorBrand = "qt  "
)

const (
	Mp4Extension Extension = ".MP4"
	MovExtension Extension = ".MOV"
)
