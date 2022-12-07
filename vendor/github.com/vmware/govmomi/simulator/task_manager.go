package simulator

import (
	"sync"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/simulator/esx"
	"github.com/vmware/govmomi/simulator/vpx"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var recentTaskMax = 200 // the VC limit

type TaskManager struct {
	mo.TaskManager
	sync.Mutex
}

func NewTaskManager(ref types.ManagedObjectReference) object.Reference {
	s := &TaskManager{}
	s.Self = ref
	if Map.IsESX() {
		s.Description = esx.Description
	} else {
		s.Description = vpx.Description
	}
	Map.AddHandler(s)
	return s
}

func (m *TaskManager) PutObject(obj mo.Reference) {
	ref := obj.Reference()
	if ref.Type != "Task" {
		return
	}

	m.Lock()
	recent := append(m.RecentTask, ref)
	if len(recent) > recentTaskMax {
		recent = recent[1:]
	}

	Map.Update(m, []types.PropertyChange{{Name: "recentTask", Val: recent}})
	m.Unlock()
}

func (*TaskManager) RemoveObject(types.ManagedObjectReference) {}

func (*TaskManager) UpdateObject(mo.Reference, []types.PropertyChange) {}
