package gsignercache

import "github.com/ethereum/go-ethereum/core/types"

func Wrap(signer types.Signer) types.Signer {
	// There used to be a different signer type wrapping the passed in
	// signer to add a global cache for results. However, this feature
	// is no longer needed since go-ethereum's core library is caching
	// results as well. The utilization of both actually caused the
	// accidential disabling of either.
	return signer
}
