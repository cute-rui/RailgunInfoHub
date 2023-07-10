package railgun

import (
	"reflect"
	"strings"
	"unsafe"
)

func stringBuilder(p ...string) string {
	var (
		b strings.Builder
		c int
	)
	l := len(p)
	for i := 0; i < l; i++ {
		c += len(p[i])
	}
	b.Grow(c)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func stringJoiner(sep string, p ...string) string {
	var notEmpty []string
	for i := range p {
		if p[i] != `` {
			notEmpty = append(notEmpty, p[i])
		}
	}

	return strings.Join(notEmpty, sep)
}

const (
	QualitySP_None = iota
	QualitySP_DolbyVision
	QualitySP_HDR
	QualitySP_DolbyAndHDR
)

func getBestQualityVideo(qlist []int) (int, int) {
	var (
		tmp       int
		qualitySP int
	)

	for i := range qlist {
		switch qlist[i] {
		case 126:
			qualitySP += 1
			continue
		case 125:
			qualitySP += 2
			continue
		default:
			if qlist[i] > tmp {
				tmp = qlist[i]
			}
		}
	}

	return tmp, qualitySP
}

func ByteToString(p []byte) string {
	var b strings.Builder
	l := len(p)
	b.Grow(l)
	for i := 0; i < l; i++ {
		b.WriteByte(p[i])
	}
	return b.String()
}

func StringToByte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func deleteDolbyAndGetBestQualityVideo(qlist []int) int {
	var tmp int

	for i := range qlist {
		if qlist[i] == 126 {
			continue
		}

		if qlist[i] > tmp {
			tmp = qlist[i]
		}
	}

	return tmp
}

func getQualityString(i int) string {
	switch i {
	case 127:
		return `超高清 8K`
	case 126:
		return `杜比视界`
	case 125:
		return `HDR 真彩色`
	case 120:
		return `4K 超清`
	case 116:
		return `1080P60 高帧率`
	case 112:
		return `1080P+ 高码率`
	case 80:
		return `1080P 高清`
	case 74:
		return `720P60 高帧率`
	case 64:
		return `720P 高清`
	case 32:
		return `480P 清晰`
	case 16:
		return `360P 流畅`
	case 6:
		return `240P 极速`
	default:
		return ``
	}
}

func checkRegion(text string) Region {
	switch {
	case strings.Contains(text, `僅限港澳`):
		return REGION_HK
	case strings.Contains(text, `僅限台灣`):
		return REGION_TW
	}

	return Config.getDefaultRegion()
}

type StreamArgs []*DownloadArgs

func (c StreamArgs) Len() int {
	return len(c)
}
func (c StreamArgs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c StreamArgs) Less(i, j int) bool {
	if len(c[i].CID) == len(c[j].CID) {
		return c[i].CID < c[j].CID
	} else {
		return len(c[i].CID) < len(c[j].CID)
	}
}
