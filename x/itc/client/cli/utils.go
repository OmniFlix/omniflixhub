package cli

import (
	"encoding/json"
	"fmt"
	"github.com/OmniFlix/omniflixhub/x/itc/types"
	"github.com/spf13/pflag"
	"os"
)

// validate nft details object
func validate(nft *types.NFTDetails) error {
	if nft.DenomId == "" {
		return fmt.Errorf("denom id is required")
	}

	if nft.Name == "" {
		return fmt.Errorf("nft name is required")
	}

	if nft.MediaUri == "" {
		return fmt.Errorf("nft media uri is required")
	}
	return nil
}

func parseNftDetails(fs *pflag.FlagSet) (*types.NFTDetails, error) {
	nftDetails := &types.NFTDetails{}
	nftDetailsFile, err := fs.GetString(FlagNftDetailsFile)
	if err != nil {
		return nil, err
	}
	if nftDetailsFile == "" {
		return nil, fmt.Errorf("file path not provided")
	}
	contents, err := os.ReadFile(nftDetailsFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, nftDetails)
	if err != nil {
		return nil, err
	}

	if err := validate(nftDetails); err != nil {
		return nil, err
	}

	return nftDetails, nil
}
