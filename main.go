package main

import (
	"fmt"
        "time"
	"log"
	"net"

        "gobot.io/x/gobot/drivers/gpio"
        "gobot.io/x/gobot/platforms/raspi"
	
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/paypal/gatt/examples/service"
)

func main() {
	Unlock()
	time.Sleep(1 * time.Second)
	Lock()
	time.Sleep(1 * time.Second)
	Unlock()
	time.Sleep(1 * time.Second)
	Lock()
	time.Sleep(1 * time.Second)
	Unlock()
	time.Sleep(1 * time.Second)
	Lock()
	time.Sleep(1 * time.Second)
	Unlock()
	time.Sleep(1 * time.Second)
	Lock()
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			// Setup GAP and GATT services for Linux implementation.
			// OS X doesn't export the access of these services.
			var macs []string
			ifs, _ := net.Interfaces()
			for _, v := range ifs {
			    h := v.HardwareAddr.String()
			    if len(h) == 0 {
			        continue
			    }
			    macs = append(macs, h)
			}
			fmt.Println(macs);
			d.AddService(service.NewGapService("h4shub:" + macs[0])) // no effect on OS X
			d.AddService(service.NewGattService())        // no effect on OS X

			// A simple count service for demo.
			s1 := NewLockService()
			d.AddService(s1)


			// Advertise device name and service's UUIDs.
			d.AdvertiseNameAndServices("h4shub:" + macs[0], []gatt.UUID{s1.UUID()})

			// Advertise as an OpenBeacon iBeacon
			d.AdvertiseIBeacon(gatt.MustParseUUID("AA6062F098CA42118EC4193EB73CCEB6"), 1, 2, -59)

		default:
		}
	}

	d.Init(onStateChanged)
	select {}
}

 
func NewLockService() *gatt.Service {
	s := gatt.NewService(gatt.MustParseUUID("09fc95c0-c111-11e3-9904-0002a5d5c51b"))
	s.AddCharacteristic(gatt.MustParseUUID("11fac9e0-c111-11e3-9246-0002a5d5c51b")).HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
			status := getLockStatus()
			changeLockStatus()
			newStatus := getLockStatus()	
			fmt.Fprintf(rsp, "was %s and is now %s", status, newStatus)
		})

	s.AddCharacteristic(gatt.MustParseUUID("16fe0d80-c111-11e3-b8c8-0002a5d5c51b")).HandleWriteFunc(
		func(r gatt.Request, data []byte) (status byte) {
			oldstatus := getLockStatus()
			changeLockStatus()
			newStatus := getLockStatus()	
			log.Printf("was %s and is now %s \n", oldstatus, newStatus)
			log.Println("Wrote:", string(data))
			return gatt.StatusSuccess
		})

	s.AddCharacteristic(gatt.MustParseUUID("1c927b50-c116-11e3-8a33-0800200c9a66")).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			cnt := 0
			for !n.Done() {
				fmt.Fprintf(n, "Count: %d", cnt)
				cnt++
				time.Sleep(time.Second)
			}
		})

	return s
}

func changeLockStatus(){
	if(getLock().State()) {
		Unlock()
	} else {
		Lock()
	}
}
func getLockStatus() string {
	status := ""
	if(getLock().State()) {
		status = "locked"
	} else {
		status = "unlocked"
	}
	return status
}

func getLock() *gpio.LedDriver {
	r := raspi.NewAdaptor()
        return gpio.NewLedDriver(r, "7")
}
func Unlock() {
	l := getLock()
	l.Off()
}

func Lock() {
	l := getLock()
        l.On()
}
