package wurflgo

import (
	"errors"
)

var chain = NewChain()

type DeviceProperties struct {
	BrandName string `json:"brand_name"`
	ModelName string `json:"model_name"`
	MarketingName string `json:"marketing_name"`
	PreferredMarkup string `json:"preferred_markup"`
	ResolutionWidth string `json:"resolution_width"`
	ResolutionHeight string `json:"resolution_height"`
	DeviceOs string `json:"device_os"`
	DeviceOsVersion string `json:"device_os_version"`
	BrowserName string `json:"mobile_browser"`
	BrowserVersion string `json:"mobile_browser_version"`
}

type Device struct {
	Id               string
	UA               string
	Parent		 string
	Children         map[string]bool
	ActualDeviceRoot bool
	Capabilities     map[string]string
	Properties	 *DeviceProperties
}

type Repository struct {
	initialized bool
	devices map[string]*Device
}

func NewRepository() *Repository {
	r := new(Repository)
	r.devices = make(map[string]*Device)
	return r
}

func (r *Repository) find(id string) *Device {
	return r.devices[id]
}

func (r *Repository) count() int {
	return len(r.devices)
}

func (r *Repository) Match(ua string) *Device {
	m := chain.Match(ua)
	return r.find(m)
}

func (r *Repository) Initialize() {
	if r.initialized {
		return
	}
	for _, dev := range r.devices {
		chain.Filter(dev.UA, dev.Id)
	}
}

func (r *Repository) register(id, ua string, actualDeviceRoot bool, capabilities map[string]string, parent string) error {
	dev := new(Device)
	dev.Id = id
	dev.UA = ua
	dev.Children = map[string]bool{}
	dev.Capabilities = capabilities
	dev.Parent = parent
	parentDevice, found := r.devices[parent]
	if parent == "" {
		dev.Capabilities = capabilities
	} else {
		if found == true {
			parentDevice.Children[dev.Id] = true
		} else {
			return errors.New("Unregistered Parent Device")
		}
	}
	//Save it
	r.devices[dev.Id] = dev

	return nil
}

func (dev *Device) getCapabilities(r *Repository) map[string]string {
	parentDevice, found := r.devices[dev.Parent]
	capabilities := map[string]string{}
	if found == true {
		for k := range parentDevice.Capabilities {
			capabilities[k] = parentDevice.Capabilities[k]
		}
		for k := range dev.Capabilities {
			capabilities[k] = dev.Capabilities[k]
		}
	}
	return capabilities
}

func (dev *Device) GetProperties(r *Repository) *DeviceProperties {
	capabilities := dev.getCapabilities(r)
	return &DeviceProperties{
		BrandName: capabilities["brand_name"],
		ModelName: capabilities["model_name"],
		MarketingName: capabilities["marketing_name"],
		PreferredMarkup: capabilities["preferred_markup"],
		ResolutionWidth: capabilities["resolution_width"],
		ResolutionHeight: capabilities["resolution_height"],
		DeviceOs: capabilities["device_os"],
		DeviceOsVersion: capabilities["device_os_version"],
		BrowserName: capabilities["mobile_browser"],
		BrowserVersion: capabilities["mobile_browser_version"],
	}
}

func init() {
	genericNormalizers := CreateGenericNormalizers()
	chain.AddHandler(NewJavaMidletHandler(genericNormalizers))
	chain.AddHandler(NewSmartTVHandler(genericNormalizers))
	kindleNormalizer := genericNormalizers.AddNormalizer(NewKindle())
	chain.AddHandler(NewKindleHandler(kindleNormalizer))
	lgPlusNormalizer := genericNormalizers.AddNormalizer(NewLGPLUS())
	chain.AddHandler(NewLGPLUSHandler(lgPlusNormalizer))
	androidNormalizer := genericNormalizers.AddNormalizer(NewAndroid())
	chain.AddHandler(NewAndroidHandler(androidNormalizer))
	chain.AddHandler(NewAppleHandler(genericNormalizers))
	chain.AddHandler(NewWindowsPhoneDesktopHandler(genericNormalizers))
	chain.AddHandler(NewWindowsPhoneHandler(genericNormalizers))
	chain.AddHandler(NewNokiaOviBrowserHandler(genericNormalizers))
	chain.AddHandler(NewNokiaHandler(genericNormalizers))
	chain.AddHandler(NewSamsungHandler(genericNormalizers))
	chain.AddHandler(NewBlackBerryHandler(genericNormalizers))
	chain.AddHandler(NewSonyEricssonHandler(genericNormalizers))
	chain.AddHandler(NewMotorolaHandler(genericNormalizers))
	chain.AddHandler(NewAlcatelHandler(genericNormalizers))
	chain.AddHandler(NewBenQHandler(genericNormalizers))
	chain.AddHandler(NewDoCoMoHandler(genericNormalizers))
	chain.AddHandler(NewGrundigHandler(genericNormalizers))
	htcMacHandler := genericNormalizers.AddNormalizer(NewHTCMac())
	chain.AddHandler(NewHTCMacHandler(htcMacHandler))
	chain.AddHandler(NewHTCHandler(genericNormalizers))
	chain.AddHandler(NewKDDIHandler(genericNormalizers))
	chain.AddHandler(NewKyoceraHandler(genericNormalizers))

	lgNormalizer := genericNormalizers.AddNormalizer(NewLG())
	chain.AddHandler(NewLGHandler(lgNormalizer))

	chain.AddHandler(NewMitsubishiHandler(genericNormalizers))
	chain.AddHandler(NewNecHandler(genericNormalizers))
	chain.AddHandler(NewNintendoHandler(genericNormalizers))
	chain.AddHandler(NewPanasonicHandler(genericNormalizers))
	chain.AddHandler(NewPantechHandler(genericNormalizers))
	chain.AddHandler(NewPhilipsHandler(genericNormalizers))
	chain.AddHandler(NewPortalmmmHandler(genericNormalizers))
	chain.AddHandler(NewQtekHandler(genericNormalizers))
	chain.AddHandler(NewReksioHandler(genericNormalizers))
	chain.AddHandler(NewSagemHandler(genericNormalizers))
	chain.AddHandler(NewSanyoHandler(genericNormalizers))
	chain.AddHandler(NewSharpHandler(genericNormalizers))
	chain.AddHandler(NewSiemensHandler(genericNormalizers))
	chain.AddHandler(NewSPVHandler(genericNormalizers))
	chain.AddHandler(NewToshibaHandler(genericNormalizers))
	chain.AddHandler(NewVodafoneHandler(genericNormalizers))

	webosNormalizer := genericNormalizers.AddNormalizer(NewWebOS())
	chain.AddHandler(NewWebOSHandler(webosNormalizer))

	chain.AddHandler(NewOperaMiniHandler(genericNormalizers))

	// Robots / Crawlers.
	chain.AddHandler(NewBotCrawlerTranscoderHandler(genericNormalizers))

	// Desktop Browsers.
	chromeNormalizer := genericNormalizers.AddNormalizer(NewChrome())
	chain.AddHandler(NewChromeHandler(chromeNormalizer))

	firefoxNormalizer := genericNormalizers.AddNormalizer(NewFirefox())
	chain.AddHandler(NewFirefoxHandler(firefoxNormalizer))

	msieNormalizer := genericNormalizers.AddNormalizer(NewMSIE())
	chain.AddHandler(NewMSIEHandler(msieNormalizer))

	operaNormalizer := genericNormalizers.AddNormalizer(NewOpera())
	chain.AddHandler(NewOperaHandler(operaNormalizer))

	safariNormalizer := genericNormalizers.AddNormalizer(NewSafari())
	chain.AddHandler(NewSafariHandler(safariNormalizer))

	konquerorNormalizer := genericNormalizers.AddNormalizer(NewKonqueror())
	chain.AddHandler(NewKonquerorHandler(konquerorNormalizer))

	// All other requests.
	chain.AddHandler(NewCatchAllHandler(genericNormalizers))

}

func CreateGenericNormalizers() *UserAgentNormalizer {
	return NewUserAgentNormalizer([]Normalizer{
		NewUPLink(),
		NewBlackBerry(),
		NewYesWap(),
		NewBabelFish(),
		NewSerialNumber(),
		NewNovarraGoogleTranslator(),
		NewLocaleRemover(),
		NewUCWEB(),
	})
}
