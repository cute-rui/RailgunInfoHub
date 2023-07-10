package railgun

type needRenewCookieResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl,omitempty"`
	Data    struct {
		Refresh   bool `json:"refresh"`
		Timestamp int  `json:"timestamp"`
	} `json:"data,omitempty"`
}

type refreshCookieResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl,omitempty"`
	Data    struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data,omitempty"`
}

type navResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl,omitempty"`
	Data    struct {
		IsLogin       bool   `json:"isLogin"`
		EmailVerified int    `json:"email_verified"`
		Face          string `json:"face"`
		FaceNft       int    `json:"face_nft"`
		FaceNftType   int    `json:"face_nft_type"`
		LevelInfo     struct {
			CurrentLevel int    `json:"current_level"`
			CurrentMin   int    `json:"current_min"`
			CurrentExp   int    `json:"current_exp"`
			NextExp      string `json:"next_exp"`
		} `json:"level_info"`
		Mid            int     `json:"mid"`
		MobileVerified int     `json:"mobile_verified"`
		Money          float64 `json:"money"`
		Moral          int     `json:"moral"`
		Official       struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"official"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"officialVerify"`
		Pendant struct {
			Pid               int    `json:"pid"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"`
		Scores       int    `json:"scores"`
		Uname        string `json:"uname"`
		VipDueDate   int64  `json:"vipDueDate"`
		VipStatus    int    `json:"vipStatus"`
		VipType      int    `json:"vipType"`
		VipPayType   int    `json:"vip_pay_type"`
		VipThemeType int    `json:"vip_theme_type"`
		VipLabel     struct {
			Path                  string `json:"path"`
			Text                  string `json:"text"`
			LabelTheme            string `json:"label_theme"`
			TextColor             string `json:"text_color"`
			BgStyle               int    `json:"bg_style"`
			BgColor               string `json:"bg_color"`
			BorderColor           string `json:"border_color"`
			UseImgLabel           bool   `json:"use_img_label"`
			ImgLabelUriHans       string `json:"img_label_uri_hans"`
			ImgLabelUriHant       string `json:"img_label_uri_hant"`
			ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
			ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
		} `json:"vip_label"`
		VipAvatarSubscript int    `json:"vip_avatar_subscript"`
		VipNicknameColor   string `json:"vip_nickname_color"`
		Vip                struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path                  string `json:"path"`
				Text                  string `json:"text"`
				LabelTheme            string `json:"label_theme"`
				TextColor             string `json:"text_color"`
				BgStyle               int    `json:"bg_style"`
				BgColor               string `json:"bg_color"`
				BorderColor           string `json:"border_color"`
				UseImgLabel           bool   `json:"use_img_label"`
				ImgLabelUriHans       string `json:"img_label_uri_hans"`
				ImgLabelUriHant       string `json:"img_label_uri_hant"`
				ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
				ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptUrl string `json:"avatar_subscript_url"`
			TvVipStatus        int    `json:"tv_vip_status"`
			TvVipPayType       int    `json:"tv_vip_pay_type"`
			TvDueDate          int    `json:"tv_due_date"`
		} `json:"vip"`
		Wallet struct {
			Mid           int `json:"mid"`
			BcoinBalance  int `json:"bcoin_balance"`
			CouponBalance int `json:"coupon_balance"`
			CouponDueTime int `json:"coupon_due_time"`
		} `json:"wallet"`
		HasShop        bool   `json:"has_shop"`
		ShopUrl        string `json:"shop_url"`
		AllowanceCount int    `json:"allowance_count"`
		AnswerStatus   int    `json:"answer_status"`
		IsSeniorMember int    `json:"is_senior_member"`
		WbiImg         struct {
			ImgUrl string `json:"img_url"`
			SubUrl string `json:"sub_url"`
		} `json:"wbi_img"`
		IsJury bool `json:"is_jury"`
	} `json:"data,omitempty"`
}
type fetchSeasonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Result  struct {
		Areas []struct {
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"areas,omitempty"`
		Evaluate string `json:"evaluate,omitempty"`
		Cover    string `json:"cover,omitempty"`
		NewEp    struct {
			Desc string `json:"desc,omitempty"`
		} `json:"new_ep,omitempty"`
		Title       string `json:"title,omitempty"`
		ShareUrl    string `json:"share_url,omitempty"`
		SeasonId    int    `json:"season_id,omitempty"`
		SeasonTitle string `json:"season_title,omitempty"`
		Episodes    []struct {
			Bvid      string `json:"bvid,omitempty"`
			Cid       int    `json:"cid,omitempty"`
			Title     string `json:"title,omitempty"`
			LongTitle string `json:"long_title,omitempty"`
			Id        int    `json:"id,omitempty"`
		} `json:"episodes,omitempty"`
	} `json:"result,omitempty"`
}

type fetchMediaResp struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Result  struct {
		Media struct {
			Areas []struct {
				Name string `json:"name,omitempty"`
			} `json:"areas,omitempty"`
			Cover string `json:"cover,omitempty"`
			NewEp struct {
				IndexShow string `json:"index_show,omitempty"`
			} `json:"new_ep,omitempty"`
			ShareUrl string `json:"share_url,omitempty"`
			Title    string `json:"title,omitempty"`
			TypeName string `json:"type_name,omitempty"`
			SeasonId int    `json:"season_id,omitempty"`
		} `json:"media,omitempty"`
	} `json:"result,omitempty"`
}

type fetchVideoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Bvid     string `json:"bvid,omitempty"`
		Aid      int64  `json:"aid,omitempty"`
		Pic      string `json:"pic,omitempty"`
		Title    string `json:"title,omitempty"`
		Cid      int    `json:"cid,omitempty"`
		Pubdate  int64  `json:"pubdate,omitempty"`
		Ctime    int64  `json:"ctime,omitempty"`
		Desc     string `json:"desc,omitempty"`
		Duration int64  `json:"duration,omitempty"`
		Owner    struct {
			Name string `json:"name,omitempty"`
		} `json:"owner,omitempty"`
		Dynamic string `json:"dynamic,omitempty"`
		Pages   []struct {
			Cid  int    `json:"cid,omitempty"`
			Part string `json:"part,omitempty"`
		} `json:"pages,omitempty"`
		UgcSeason struct {
			Id int `json:"id"`
		} `json:"ugc_season,omitempty"`
	} `json:"data,omitempty"`
}

type respVideoDetail struct {
	Code    int       `json:"code"`
	Message string    `json:"message,omitempty"`
	Data    mediaData `json:"data,omitempty"`
}

type respSeasonDetail struct {
	Code    int       `json:"code"`
	Message string    `json:"message,omitempty"`
	Data    mediaData `json:"result,omitempty"`
}

type mediaData struct {
	Quality       int    `json:"quality"`
	AcceptFormat  string `json:"accept_format"`
	AcceptQuality []int  `json:"accept_quality"`
	Durl          []struct {
		Order     int      `json:"order"`
		Url       string   `json:"url"`
		BackupUrl []string `json:"backup_url"`
	} `json:"durl,omitempty"`
	Dash struct {
		Video []struct {
			Id         int      `json:"id"`
			BaseUrl    string   `json:"baseUrl"`
			Base_Url   string   `json:"base_url"`
			BackupUrl  []string `json:"backupUrl"`
			Backup_Url []string `json:"backup_url"`
			Codecs     string   `json:"codecs"`
		} `json:"video"`
		Audio []struct {
			Id         int      `json:"id"`
			BaseUrl    string   `json:"baseUrl"`
			Base_Url   string   `json:"base_url"`
			BackupUrl  []string `json:"backupUrl"`
			Backup_Url []string `json:"backup_url"`
			Bandwidth  int      `json:"bandwidth"`
		} `json:"audio"`
		Dolby struct {
			Type  int `json:"type"`
			Audio []struct {
				Id          int      `json:"id"`
				BaseUrl     string   `json:"base_url"`
				BackupUrl   []string `json:"backup_url"`
				Bandwidth   int      `json:"bandwidth"`
				MimeType    string   `json:"mime_type"`
				Codecs      string   `json:"codecs"`
				SegmentBase struct {
					Initialization string `json:"initialization"`
					IndexRange     string `json:"index_range"`
				} `json:"segment_base"`
				Size int `json:"size"`
			} `json:"audio,omitempty"`
		} `json:"dolby,omitempty"`
		Flac struct {
			Display bool `json:"display,omitempty"`
			Audio   struct {
				Id          int      `json:"id"`
				BaseUrl     string   `json:"base_url"`
				BackupUrl   []string `json:"backup_url"`
				Bandwidth   int      `json:"bandwidth"`
				MimeType    string   `json:"mime_type"`
				Codecs      string   `json:"codecs"`
				SegmentBase struct {
					Initialization string `json:"initialization"`
					IndexRange     string `json:"index_range"`
				} `json:"segment_base"`
				Size int `json:"size"`
			} `json:"audio,omitempty"`
		} `json:"flac,omitempty"`
	} `json:"dash,omitempty"`
}

type fetchCollectionResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Aids     []int `json:"aids"`
		Archives []struct {
			Aid              int    `json:"aid"`
			Bvid             string `json:"bvid"`
			Ctime            int    `json:"ctime"`
			Duration         int    `json:"duration"`
			InteractiveVideo bool   `json:"interactive_video"`
			Pic              string `json:"pic"`
			Pubdate          int    `json:"pubdate"`
			Stat             struct {
				View int `json:"view"`
			} `json:"stat"`
			State  int    `json:"state"`
			Title  string `json:"title"`
			UgcPay int    `json:"ugc_pay"`
		} `json:"archives"`
		Meta struct {
			Category    int    `json:"category"`
			Cover       string `json:"cover"`
			Description string `json:"description"`
			Mid         int    `json:"mid"`
			Name        string `json:"name"`
			Ptime       int    `json:"ptime"`
			SeasonId    int    `json:"season_id"`
			Total       int    `json:"total"`
		} `json:"meta"`
		Page struct {
			PageNum  int `json:"page_num"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"page"`
	} `json:"data,omitempty"`
}

type fetchPlayerInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Aid      int    `json:"aid"`
		Bvid     string `json:"bvid"`
		AllowBp  bool   `json:"allow_bp"`
		NoShare  bool   `json:"no_share"`
		Cid      int    `json:"cid"`
		MaxLimit int    `json:"max_limit"`
		PageNo   int    `json:"page_no"`
		HasNext  bool   `json:"has_next"`
		IpInfo   struct {
			Ip       string `json:"ip"`
			ZoneIp   string `json:"zone_ip"`
			ZoneId   int    `json:"zone_id"`
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
		} `json:"ip_info"`
		LoginMid     int    `json:"login_mid"`
		LoginMidHash string `json:"login_mid_hash"`
		IsOwner      bool   `json:"is_owner"`
		Name         string `json:"name"`
		Permission   string `json:"permission"`
		LevelInfo    struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Vip struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path        string `json:"path"`
				Text        string `json:"text"`
				LabelTheme  string `json:"label_theme"`
				TextColor   string `json:"text_color"`
				BgStyle     int    `json:"bg_style"`
				BgColor     string `json:"bg_color"`
				BorderColor string `json:"border_color"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptUrl string `json:"avatar_subscript_url"`
		} `json:"vip"`
		AnswerStatus int    `json:"answer_status"`
		BlockTime    int    `json:"block_time"`
		Role         string `json:"role"`
		LastPlayTime int    `json:"last_play_time"`
		LastPlayCid  int    `json:"last_play_cid"`
		NowTime      int    `json:"now_time"`
		OnlineCount  int    `json:"online_count"`
		Subtitle     struct {
			AllowSubmit bool   `json:"allow_submit"`
			Lan         string `json:"lan"`
			LanDoc      string `json:"lan_doc"`
			Subtitles   []struct {
				Id          int64  `json:"id"`
				Lan         string `json:"lan"`
				LanDoc      string `json:"lan_doc"`
				IsLock      bool   `json:"is_lock"`
				SubtitleUrl string `json:"subtitle_url"`
				Type        int    `json:"type"`
				IdStr       string `json:"id_str"`
				AiType      int    `json:"ai_type"`
			} `json:"subtitles"`
		} `json:"subtitle"`
		ViewPoints      []interface{} `json:"view_points"`
		IsUgcPayPreview bool          `json:"is_ugc_pay_preview"`
		PreviewToast    string        `json:"preview_toast"`
		PcdnLoader      struct {
			Flv struct {
				Labels struct {
					PcdnVideoType string `json:"pcdn_video_type"`
					PcdnStage     string `json:"pcdn_stage"`
					PcdnGroup     string `json:"pcdn_group"`
					PcdnVersion   string `json:"pcdn_version"`
					PcdnVendor    string `json:"pcdn_vendor"`
				} `json:"labels"`
			} `json:"flv"`
			Dash struct {
				Labels struct {
					PcdnVideoType string `json:"pcdn_video_type"`
					PcdnStage     string `json:"pcdn_stage"`
					PcdnGroup     string `json:"pcdn_group"`
					PcdnVersion   string `json:"pcdn_version"`
					PcdnVendor    string `json:"pcdn_vendor"`
				} `json:"labels"`
			} `json:"dash"`
		} `json:"pcdn_loader"`
		Options struct {
			Is360      bool `json:"is_360"`
			WithoutVip bool `json:"without_vip"`
		} `json:"options"`
		GuideAttention []interface{} `json:"guide_attention"`
		JumpCard       []interface{} `json:"jump_card"`
		OperationCard  []interface{} `json:"operation_card"`
		OnlineSwitch   struct {
			EnableGrayDashPlayback string `json:"enable_gray_dash_playback"`
			NewBroadcast           string `json:"new_broadcast"`
			RealtimeDm             string `json:"realtime_dm"`
			SubtitleSubmitSwitch   string `json:"subtitle_submit_switch"`
		} `json:"online_switch"`
		Fawkes struct {
			ConfigVersion int `json:"config_version"`
			FfVersion     int `json:"ff_version"`
		} `json:"fawkes"`
		ShowSwitch struct {
			LongProgress bool `json:"long_progress"`
		} `json:"show_switch"`
	} `json:"data,omitempty"`
}
