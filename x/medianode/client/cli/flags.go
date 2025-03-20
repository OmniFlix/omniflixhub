package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagURL           = "url"
	FlagPricePerHour  = "price-per-hour"
	FlagDeposit       = "deposit"
	FlagLeaseAmount   = "lease-amount"
	FlagLeaseHours    = "lease-hours"
	FlagHardwareSpecs = "hardware-specs"
	FlagOwner         = "owner"
	FlagStatus        = "status"
	FlagNodeMoniker   = "node-moniker"
	FlagDescription   = "description"
	FlagContact       = "contact"
)

var (
	FsRegisterMediaNode = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateMediaNode   = flag.NewFlagSet("", flag.ContinueOnError)
	FsLeaseMediaNode    = flag.NewFlagSet("", flag.ContinueOnError)
	FsExtendLease       = flag.NewFlagSet("", flag.ContinueOnError)
	FsDepositMediaNode  = flag.NewFlagSet("", flag.ContinueOnError)
	FsCancelLease       = flag.NewFlagSet("", flag.ContinueOnError)
	FsCloseMediaNode    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	// Register flags for registering a media node
	FsRegisterMediaNode.String(FlagURL, "", "URL of the media node")
	FsRegisterMediaNode.String(FlagHardwareSpecs, "", "Hardware specifications of the media node")
	FsRegisterMediaNode.String(FlagPricePerHour, "", "Lease price per hour")
	FsRegisterMediaNode.String(FlagDeposit, "", "initial deposit amount")
	FsRegisterMediaNode.String(FlagNodeMoniker, "", "name of the medianode")
	FsRegisterMediaNode.String(FlagDescription, "", "description of the medianode")
	FsRegisterMediaNode.String(FlagContact, "", "contact details of medianode")

	// Register flags for updating a media node
	FsUpdateMediaNode.String(FlagNodeMoniker, "", "New Moniker of the media node")
	FsUpdateMediaNode.String(FlagDescription, "", "update description of the media node")
	FsUpdateMediaNode.String(FlagContact, "", "update contact of the media node")
	FsUpdateMediaNode.String(FlagHardwareSpecs, "", "Updated hardware specifications")
	FsUpdateMediaNode.String(FlagPricePerHour, "", "Updated lease amount per our")

	// Register flags for leasing a media node
	FsLeaseMediaNode.Uint64(FlagLeaseHours, 0, "Number of hours to lease the media node")
	FsLeaseMediaNode.String(FlagLeaseAmount, "", "lease amount paying")

	// Register flags for leasing a media node
	FsExtendLease.Uint64(FlagLeaseHours, 0, "Number of hours to lease the media node")
	FsExtendLease.String(FlagLeaseAmount, "", "lease amount paying")

	// Register flags for deposit media node
	FsDepositMediaNode.String(FlagDeposit, "", "deposit amount")
}
