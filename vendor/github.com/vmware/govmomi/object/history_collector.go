package object

import (
	"context"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

type HistoryCollector struct {
	Common
}

func NewHistoryCollector(c *vim25.Client, ref types.ManagedObjectReference) *HistoryCollector {
	return &HistoryCollector{
		Common: NewCommon(c, ref),
	}
}

func (h HistoryCollector) Destroy(ctx context.Context) error {
	req := types.DestroyCollector{
		This: h.Reference(),
	}

	_, err := methods.DestroyCollector(ctx, h.c, &req)
	return err
}

func (h HistoryCollector) Reset(ctx context.Context) error {
	req := types.ResetCollector{
		This: h.Reference(),
	}

	_, err := methods.ResetCollector(ctx, h.c, &req)
	return err
}

func (h HistoryCollector) Rewind(ctx context.Context) error {
	req := types.RewindCollector{
		This: h.Reference(),
	}

	_, err := methods.RewindCollector(ctx, h.c, &req)
	return err
}

func (h HistoryCollector) SetPageSize(ctx context.Context, maxCount int32) error {
	req := types.SetCollectorPageSize{
		This:     h.Reference(),
		MaxCount: maxCount,
	}

	_, err := methods.SetCollectorPageSize(ctx, h.c, &req)
	return err
}
