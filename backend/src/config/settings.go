package config

type SettingsEntry interface {
	Key() string
	Default()
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// All global settings as a struct
type GlobalSettings struct {
	RegistrationEnabled         RegistrationEnabled         `json:"registration_enabled"`
	LoggingVerbosity            LoggingVerbosity            `json:"logging_verbosity"`
	UseCdnFonts                 UseCdnFonts                 `json:"use_cdn_fonts"`
	UseIpGeolocation            UseIpGeolocation            `json:"use_ip_geolocation"`
	EnableGravatar              EnableGravatar              `json:"enable_gravatar"`
	CacheProfilePictures        CacheProfilePictures        `json:"cache_profile_pictures"`
	EnableProfilePicturesUpload EnableProfilePicturesUpload `json:"enable_profile_pictures_upload"`
	PublicCalendarEnabled       PublicCalendarEnabled       `json:"public_calendar_enabled"`
}

func (s *GlobalSettings) UpdateSetting(entry SettingsEntry) {
	switch entry.Key() {
	case KeyRegistrationEnabled:
		s.RegistrationEnabled.Enabled = entry.(*RegistrationEnabled).Enabled
	case KeyLoggingVerbosity:
		s.LoggingVerbosity.Verbosity = entry.(*LoggingVerbosity).Verbosity
	case KeyUseCdnFonts:
		s.UseCdnFonts.UseCdn = entry.(*UseCdnFonts).UseCdn
	case KeyUseIpGeolocation:
		s.UseIpGeolocation.UseIpGeolocation = entry.(*UseIpGeolocation).UseIpGeolocation
	case KeyEnableGravatar:
		s.EnableGravatar.Enabled = entry.(*EnableGravatar).Enabled
	case KeyCacheProfilePictures:
		s.CacheProfilePictures.Enabled = entry.(*CacheProfilePictures).Enabled
	case KeyEnableProfilePicturesUpload:
		s.EnableProfilePicturesUpload.Enabled = entry.(*EnableProfilePicturesUpload).Enabled
	case KeyPublicCalendarEnabled:
		s.PublicCalendarEnabled.Enabled = entry.(*PublicCalendarEnabled).Enabled
	default:
		// TODO: warning
	}
}
