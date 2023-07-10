package railgun

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WBICache struct {
	imgKey     string
	subKey     string
	lock       sync.RWMutex
	lastUpdate time.Time
}

func newWBICache() *WBICache {
	return &WBICache{}
}

func (r *request) setWBICache() error {
	resp, err := r.getNavInfo()
	if err != nil {
		return err
	}

	r.wbiCache.imgKey = strings.Split(strings.Split(resp.Data.WbiImg.ImgUrl, "/")[len(strings.Split(resp.Data.WbiImg.ImgUrl, "/"))-1], ".")[0]
	r.wbiCache.subKey = strings.Split(strings.Split(resp.Data.WbiImg.SubUrl, "/")[len(strings.Split(resp.Data.WbiImg.SubUrl, "/"))-1], ".")[0]
	return nil
}

func (r *request) appendQueryParamWithWBISignature(params *map[string]string) error {
	err := r.updateWBICache()
	if err != nil {
		return err
	}
	encWbi(params, r.wbiCache.imgKey, r.wbiCache.subKey)

	return nil
}

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
	33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
	61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
	36, 20, 34, 44, 52,
}

func getMixinKey(orig string) string {
	var str strings.Builder
	for _, v := range mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
		}
	}
	return str.String()[:32]
}
func encWbi(params *map[string]string, imgKey string, subKey string) {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	(*params)["wts"] = currTime
	// Sort keys
	keys := make([]string, 0, len(*params))
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Remove unwanted characters
	for k, v := range *params {
		v = strings.ReplaceAll(v, "!", "")
		v = strings.ReplaceAll(v, "'", "")
		v = strings.ReplaceAll(v, "(", "")
		v = strings.ReplaceAll(v, ")", "")
		v = strings.ReplaceAll(v, "*", "")
		(*params)[k] = v
	}
	// Build URL parameters
	var str strings.Builder
	for _, k := range keys {
		str.WriteString(fmt.Sprintf("%s=%s&", k, (*params)[k]))
	}
	query := strings.TrimSuffix(str.String(), "&")
	// Calculate w_rid
	hash := md5.Sum([]byte(query + mixinKey))
	(*params)["w_rid"] = hex.EncodeToString(hash[:])
}

func (r *request) updateWBICache() error {
	if time.Now().Sub(r.wbiCache.lastUpdate).Minutes() < 10 {
		return nil
	}

	if !r.wbiCache.lock.TryLock() {
		return nil
	}

	defer r.wbiCache.lock.Unlock()

	err := r.setWBICache()
	if err != nil {
		return err
	}

	r.wbiCache.lastUpdate = time.Now()

	return nil
}
