package main

import (
	"C"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

var idCounter int
var instances map[int]*Plugin
var mu sync.Mutex

func init() {
	// init globals
	instances = make(map[int]*Plugin)
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "plugin", "My Go Plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	// define id
	id := idCounter
	idCounter += 1

	// init plugin
	plugin := &Plugin{}

	// cache instance
	instances[id] = plugin

	// set context
	output.FLBPluginSetContext(ctx, id)

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	fmt.Println("Flush called for unknown instance")
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Return options:
	//
	// output.FLB_OK    = data have been processed.
	// output.FLB_ERROR = unrecoverable error, do not try this again.
	// output.FLB_RETRY = retry to flush later.

	// get plugin from context
	id := output.FLBPluginGetContext(ctx).(int)
	plugin := instances[id]

	// create fluent-bit decoder
	dec := output.NewDecoder(data, int(length))

	// get records
	var records []Record
	for {
		// extract record
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		// parse timestamp
		var timestamp time.Time
		switch t := ts.(type) {
		case output.FLBTime:
			timestamp = ts.(output.FLBTime).Time
		case uint64:
			timestamp = time.Unix(int64(t), 0)
		default:
			timestamp = time.Now()
		}
		record["_flb_ts"] = timestamp

		// append to list
		records = append(records, record)
	}

	plugin.Flush(C.GoString(tag), records)

	// data has been processed
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	mu.Lock()
	idCounter -= 1
	plugin := instances[idCounter]
	mu.Unlock()

	// teardown plugin
	plugin.Teardown()

	return output.FLB_OK
}

type Record = map[interface{}]interface{}

type Plugin struct{}

func (p *Plugin) Flush(tag string, records []Record) {
	fmt.Println("plugin Flush() start")
	fmt.Println(records)
	time.Sleep(5 * time.Second)
	fmt.Println("plugin Flush() ended")
}

func (p *Plugin) Teardown() {
	fmt.Println("plugin Teardown()")
}

func main() {
}
