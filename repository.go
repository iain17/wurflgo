package wurflgo

import (
	"errors"
	"github.com/srinathgs/wurflgo/stringSet"
)

var chain = NewChain()

func GetChain() *Chain {
	return chain
}

func GetUtil() *Util {
	return util
}

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
	Parent           *Device
	Children         stringSet.Set
	ActualDeviceRoot bool
	Capabilities     map[string]string
	Properties	 *DeviceProperties
}

type Repository struct {
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

func (r *Repository) register(id, ua string, actualDeviceRoot bool, capabilities map[string]string, parent string) error {
	dev := new(Device)
	dev.Id = id
	dev.UA = ua
	dev.Children = stringSet.New()
	dev.Capabilities = make(map[string]string)
	if parent == "" {
		dev.Capabilities = capabilities
	} else {
		parentDevice, found := r.devices[parent]
		if found == true {
			for k := range parentDevice.Capabilities {
				dev.Capabilities[k] = parentDevice.Capabilities[k]
			}
			for k := range capabilities {
				dev.Capabilities[k] = capabilities[k]
			}
			dev.Parent = parentDevice
			parentDevice.Children.Add(dev.Id)
		} else {
			//fmt.Println(dev.Id)
			return errors.New("Unregistered Parent Device")
		}
	}
	dev.parseDeviceProperties()
	r.devices[dev.Id] = dev
	chain.Filter(dev.UA, dev.Id)
	return nil
}

func (r *Device) parseDeviceProperties() {
	r.Properties = &DeviceProperties{
		BrandName: r.Capabilities["brand_name"],
		ModelName: r.Capabilities["model_name"],
		MarketingName: r.Capabilities["marketing_name"],
		PreferredMarkup: r.Capabilities["preferred_markup"],
		ResolutionWidth: r.Capabilities["resolution_width"],
		ResolutionHeight: r.Capabilities["resolution_height"],
		DeviceOs: r.Capabilities["device_os"],
		DeviceOsVersion: r.Capabilities["device_os_version"],
		BrowserName: r.Capabilities["mobile_browser"],
		BrowserVersion: r.Capabilities["mobile_browser_version"],
	}
}

func (r *Repository) Match(ua string) *Device {
	m := chain.Match(ua)
	return r.find(m)
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
