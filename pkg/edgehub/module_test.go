package edgehub

import (
	"testing"
	"time"

	"github.com/kubeedge/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/pkg/common/util"
)

// coreContext is beehive context used for communication between modules
var coreContext *context.Context

// edgeHubModule is edgeHub implementation of Module interface
var edgeHubModule core.Module

//TestName is function that registers the module and tests whether the correct name of the module is returned
func TestName(t *testing.T) {
	//Load Configurations as go test runs in /tmp
	err := util.LoadConfig()
	t.Run("AddConfigSource", func(t *testing.T) {
		if err != nil {
			t.Errorf("loading config failed with error: %v", err)
		}
	})
	modules := core.GetModules()
	core.Register(&EdgeHub{controller: NewEdgeHubController()})
	for name, module := range modules {
		if name == ModuleNameEdgeHub {
			edgeHubModule = module
			break
		}
	}
	t.Run("ModuleRegistration", func(t *testing.T) {
		if edgeHubModule == nil {
			t.Errorf("EdgeHub Module not Registered with beehive core")
			return
		}
		if ModuleNameEdgeHub != edgeHubModule.Name() {
			t.Errorf("Name of module is not correct wanted: %v and got: %v", ModuleNameEdgeHub, edgeHubModule.Name())
			return
		}
	})
}

//TestGroup is function that registers the module and tests whether the correct group name is returned
func TestGroup(t *testing.T) {
	//Load Configurations as go test runs in /tmp
	err := util.LoadConfig()
	t.Run("AddConfigSource", func(t *testing.T) {
		if err != nil {
			t.Errorf("loading config failed with error: %v", err)
		}
	})
	modules := core.GetModules()
	core.Register(&EdgeHub{controller: NewEdgeHubController()})
	for name, module := range modules {
		if name == ModuleNameEdgeHub {
			edgeHubModule = module
			break
		}
	}
	t.Run("ModuleRegistration", func(t *testing.T) {
		if edgeHubModule == nil {
			t.Errorf("EdgeHub Module not Registered with beehive core")
			return
		}
		if core.HubGroup != edgeHubModule.Group() {
			t.Errorf("Group of module is not correct wanted: %v and got: %v", core.HubGroup, edgeHubModule.Group())
		}
	})
}

//TestStart is a function to test the start of the edge hub module
func TestStart(t *testing.T) {
	coreContext = context.GetContext(context.MsgCtxTypeChannel)
	modules := core.GetModules()
	for name, module := range modules {
		coreContext.AddModule(name)
		coreContext.AddModuleGroup(name, module.Group())
	}
	go edgeHubModule.Start(coreContext)
	time.Sleep(2 * time.Second)
}

// TestCleanup is function to test cleanup
func TestCleanup(t *testing.T) {
	edgeHubModule.Cleanup()
	var test model.Message

	// Send message to avoid deadlock if channel deletion has failed after cleanup
	go coreContext.Send(ModuleNameEdgeHub, test)

	_, err := coreContext.Receive(ModuleNameEdgeHub)
	t.Run("CheckCleanUp", func(t *testing.T) {
		if err == nil {
			t.Errorf("Edgehub Module still has channel after cleanup")
		}
	})
}
