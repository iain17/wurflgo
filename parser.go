package wurflgo

import (
	"fmt"
	"encoding/xml"
	"os"
	"strings"
	"github.com/iain17/wurflgo/stringSet"
)

type Capability struct{
	Name string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Grp struct{
	//XMLName xml.Name `xml:"group"`
	Id string `xml:"id,attr"`
	Capabilities []Capability `xml:"capability"`
}

type XMLDevice struct{
	//XMLName xml.Name `xml:"device"`
	Id string `xml:"id,attr"`
	Parent string `xml:"fall_back,attr"`
	UserAgent string `xml:"user_agent,attr"`
	ActualDeviceRoot bool`xml:"actual_device_root,attr"`
	Group []Grp `xml:"group"`
}

type WurflProcessor struct{
	Groups stringSet.Set
	DeferredDevices []string
	DeviceList map[string]*XMLDevice
	ProcessedDevices stringSet.Set
	InFile *os.File
	Out *Repository
}

func NewProcessor(groups string, infile string, out *Repository) (*WurflProcessor, error){
	gpSet := stringSet.New()
	gps := strings.Split(groups,",")
	for _,V := range gps {
		gpSet.Add(V)
	}
	var err error
	wurflp := new(WurflProcessor)
	wurflp.Groups = gpSet
	wurflp.Out = out
	wurflp.DeferredDevices = []string{}
	wurflp.ProcessedDevices = stringSet.New()
	wurflp.DeviceList = make(map[string]*XMLDevice)
	wurflp.InFile,err = os.OpenFile(infile,os.O_RDONLY,0666)
	if err != nil{
		return nil,err
	}
	return wurflp,nil
}

func (wp *WurflProcessor) Process(){
	defer wp.InFile.Close()
	var err error
	dec := xml.NewDecoder(wp.InFile)
	for {
		t,_ := dec.Token()
		if t == nil {
			break
		}
		switch se := t.(type){
		case xml.StartElement:
			if se.Name.Local == "device"{
				dev := new(XMLDevice)
				if err = dec.DecodeElement(dev,&se); err == nil{
					wp.DeviceList[dev.Id] = dev
				}
				if dev.Parent == "" || dev.Parent == "root"{
					dev.Parent = ""
					wp.save(dev)
				} else {
					if wp.ProcessedDevices.Get(dev.Parent) {
						wp.save(dev)
					} else {
						wp.DeferredDevices = append(wp.DeferredDevices, dev.Id)
					}
				}

				//fmt.Printf("%#v\n",se)
		}

		//fmt.Printf("%#v\n",t)
		}
	}
	wp.ProcessDeferredDevices()
	wp.Out.Cleanup()
	wp.Out.Initialize()
}


func (wp *WurflProcessor) save(dev *XMLDevice) {
	capabilities := map[string]string{}
	//fmt.Println(dev.Id)
	for _,grp := range dev.Group{
		if wp.Groups.Get(grp.Id){
			for _,Cap := range grp.Capabilities{
				capabilities[Cap.Name] = Cap.Value
			}
		}
	}
	err := wp.Out.register(dev.Id, dev.UserAgent, dev.ActualDeviceRoot, capabilities, dev.Parent)
	if err != nil {
		panic(err)
	}
	wp.ProcessedDevices.Add(dev.Id)
}


func (wp *WurflProcessor)ProcessDeferredDevices(){
	fmt.Println("Processing Deferred Devices...")
	for len(wp.DeferredDevices) > 0{
		devId := wp.DeferredDevices[0]
		dev := wp.DeviceList[devId]
		if wp.ProcessedDevices.Get(dev.Parent){
			wp.save(dev)
			wp.DeferredDevices = wp.DeferredDevices[1:len(wp.DeferredDevices)]
		} else {
			wp.DeferredDevices = append(wp.DeferredDevices[1:len(wp.DeferredDevices)], devId)
		}
	}
}

/**
* Param: groups: list of groups you want to use separated by commas
* Param: database: Path to the wurfl xml file (product_info,xhtml_ui)
*/
func New(database string, groups string) *Repository {
	repository := NewRepository()
	wp, err := NewProcessor(groups, database, repository)
	if err != nil{
		fmt.Println("An Error Occured %s", err.Error())
		return nil
	}
	fmt.Println("Please wait loading wurfl database")
	wp.Process()
	fmt.Println(repository.count(), "devices loaded")
	return repository
}
