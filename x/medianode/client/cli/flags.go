package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagURL           = "url"
	FlagLeaseAmount   = "lease-amount"
	FlagLeaseDays     = "lease-days"
	FlagHardwareSpecs = "hardware-specs"
	FlagOwner         = "owner"
	FlagStatus        = "status"
)

var (
	FsRegisterMediaNode = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateMediaNode   = flag.NewFlagSet("", flag.ContinueOnError)
	FsLeaseMediaNode    = flag.NewFlagSet("", flag.ContinueOnError)
	FsCancelLease       = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	// Register flags for registering a media node
	FsRegisterMediaNode.String(FlagURL, "", "URL of the media node")
	FsRegisterMediaNode.String(FlagHardwareSpecs, "", "Hardware specifications of the media node")
	FsRegisterMediaNode.String(FlagLeaseAmount, "", "Lease amount per day")

	// Register flags for updating a media node
	FsUpdateMediaNode.String(FlagURL, "", "New URL of the media node")
	FsUpdateMediaNode.String(FlagHardwareSpecs, "", "Updated hardware specifications")
	FsUpdateMediaNode.String(FlagLeaseAmount, "", "Updated lease amount per day")

	// Register flags for leasing a media node
	FsLeaseMediaNode.Uint64(FlagLeaseDays, 0, "Number of days to lease the media node")

}
