package xavc

import (
	"fmt"
	"github.com/fukco/media-metadata/internal/manufacturer/sony/nrtmd"
	"github.com/fukco/media-metadata/internal/manufacturer/sony/rtmd"
	"github.com/fukco/media-metadata/internal/meta"
	"strings"
	"time"
)

type NrtmdDisp struct {
	// 厂商
	Manufacturer string
	// 型号
	ModelName string
	// 帧率
	FormatFPS string
	// 捕获帧率
	CaptureFPS string
	// 文件格式以及记录帧速率
	FileFormatAndRecFrameRate string
	// 视频比特率
	VideoBitrate string
	// 色度采样 色深
	Profile string
	// 拍摄模式
	RecordingMode string
	// 是否开启代理
	IsProxyOn bool
	// 创建时间
	CreationTimestamp int64
	// 时间码
	TimecodeSecs  int
	TimecodeFrame int
}

func (nrtmdDisp *NrtmdDisp) parseFromSonyXML(nonRealTimeMeta *nrtmd.NonRealTimeMeta) {
	nrtmdDisp.Manufacturer = nonRealTimeMeta.Device.Manufacturer
	nrtmdDisp.ModelName = nonRealTimeMeta.Device.ModelName
	nrtmdDisp.CaptureFPS = nonRealTimeMeta.VideoFormat.VideoFrame.CaptureFps
	nrtmdDisp.FormatFPS = nonRealTimeMeta.VideoFormat.VideoFrame.FormatFps
	nrtmdDisp.RecordingMode = nonRealTimeMeta.RecordingMode.Type
	var resolution, fileFormat string
	if nonRealTimeMeta.VideoFormat.VideoLayout.Pixel == "3840" {
		resolution = "4K"
	} else if nonRealTimeMeta.VideoFormat.VideoLayout.Pixel == "1080" {
		resolution = "HD"
	}
	codec := nonRealTimeMeta.VideoFormat.VideoFrame.VideoCodec
	codecSplitStrs := strings.Split(codec[:strings.Index(codec, "@")], "_")
	//https://en.wikipedia.org/wiki/Advanced_Video_Coding#Profiles
	//https://en.wikipedia.org/wiki/High_Efficiency_Video_Coding#Profiles
	if strings.HasPrefix(codecSplitStrs[0], "HEVC") {
		fileFormat = "HS"
		if codecSplitStrs[3] == "M10P" {
			// Main 10
			nrtmdDisp.Profile = "4:2:0 10"
		} else if codecSplitStrs[3] == "M42210P" {
			// Main 4:2:2 10
			nrtmdDisp.Profile = "4:2:2 10"
		}
	} else if strings.HasPrefix(codecSplitStrs[0], "AVC") {
		fileFormat = "S"
		if codecSplitStrs[3] == "HP" {
			// High Profile
			nrtmdDisp.Profile = "4:2:0 8"
		} else if codecSplitStrs[3] == "H422P" {
			// High 4:2:2 Profile
			nrtmdDisp.Profile = "4:2:2 10"
		} else if codecSplitStrs[3] == "H422IP" {
			// High 4:2:2 Intra Profile
			fileFormat = "S-I"
			nrtmdDisp.Profile = "4:2:2 10"
		}
	}
	nrtmdDisp.FileFormatAndRecFrameRate = fmt.Sprintf("XAVC %s %s %s", fileFormat, resolution, nrtmdDisp.FormatFPS)
	if len(nonRealTimeMeta.SubStream.Codec) > 0 {
		nrtmdDisp.IsProxyOn = true
	} else {
		nrtmdDisp.IsProxyOn = false
	}
	t, _ := time.Parse(time.RFC3339, nonRealTimeMeta.CreationDate.Value)
	nrtmdDisp.CreationTimestamp = t.Unix()
}

func (nrtmdDisp *NrtmdDisp) parseFromSonyRtmd(rtmd *rtmd.RTMD) {
	nrtmdDisp.TimecodeSecs = rtmd.Timecode.Sec + rtmd.Timecode.Min*60 + rtmd.Timecode.Hour*3600
	nrtmdDisp.TimecodeFrame = rtmd.Timecode.Frame
}

func (nrtmdDisp *NrtmdDisp) parseFromVideoProfile(profile *meta.VideoProfile) {
	nrtmdDisp.VideoBitrate = profile.VideoAvgBitrate
}

func NrtmdDispFromMeta(metadata *meta.Metadata) *NrtmdDisp {
	nrtmdDisp := &NrtmdDisp{}

	if metadata.Mp4Meta.VideoProfile != nil {
		nrtmdDisp.parseFromVideoProfile(metadata.Mp4Meta.VideoProfile)
	}
	if metadata.Sony != nil {
		if metadata.Sony.NonRealTimeMeta != nil {
			nrtmdDisp.parseFromSonyXML(metadata.Sony.NonRealTimeMeta)
		}
		if metadata.Sony.RTMD != nil {
			nrtmdDisp.parseFromSonyRtmd(metadata.Sony.RTMD)
		}
	}
	return nrtmdDisp
}
