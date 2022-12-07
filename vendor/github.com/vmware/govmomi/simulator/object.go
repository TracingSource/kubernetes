package simulator

import (
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func SetCustomValue(ctx *Context, req *types.SetCustomValue) soap.HasFault {
	ctx.Caller = &req.This
	body := &methods.SetCustomValueBody{}

	cfm := Map.CustomFieldsManager()

	_, field := cfm.findByNameType(req.Key, req.This.Type)
	if field == nil {
		body.Fault_ = Fault("", &types.InvalidArgument{InvalidProperty: "key"})
		return body
	}

	res := cfm.SetField(ctx, &types.SetField{
		This:   cfm.Reference(),
		Entity: req.This,
		Key:    field.Key,
		Value:  req.Value,
	})

	if res.Fault() != nil {
		body.Fault_ = res.Fault()
		return body
	}

	body.Res = &types.SetCustomValueResponse{}
	return body
}
