package simulator

import (
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func RenameTask(e mo.Entity, r *types.Rename_Task) soap.HasFault {
	task := CreateTask(e, "rename", func(t *Task) (types.AnyType, types.BaseMethodFault) {
		obj := Map.Get(r.This).(mo.Entity).Entity()

		if parent, ok := Map.Get(*obj.Parent).(*Folder); ok {
			if Map.FindByName(r.NewName, parent.ChildEntity) != nil {
				return nil, &types.InvalidArgument{InvalidProperty: "name"}
			}
		}

		Map.Update(e, []types.PropertyChange{{Name: "name", Val: r.NewName}})

		return nil, nil
	})

	return &methods.Rename_TaskBody{
		Res: &types.Rename_TaskResponse{
			Returnval: task.Run(),
		},
	}
}
